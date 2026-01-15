package handler

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	"no-as-a-service/internal/helper"
	"no-as-a-service/internal/response"

	"github.com/gin-gonic/gin"
)

// reasonCache stores loaded reasons per language to avoid repeated file reads
var (
	reasonCache        = make(map[string][]string)
	cacheMutex         sync.RWMutex
	availableLanguages []string
	langMutex          sync.Once
)

func Reason(ctx *gin.Context) {
	// Get language from query param, default from .env or if not set "en"
	defaultLanguage := helper.GetEnv("DEFAULT_LANGUAGE", "en")
	lang := ctx.DefaultQuery("lang", defaultLanguage)

	// Get available languages from data directory
	languages := getAvailableLanguages()

	// Validate language parameter
	if !isLanguageSupported(lang, languages) {
		ctx.JSON(http.StatusBadRequest, response.DefaultResponse{
			Payload: nil,
			Status: response.StatusResponse{
				Status:  http.StatusBadRequest,
				Message: fmt.Sprintf("Invalid language. Supported: %s", strings.Join(languages, ", ")),
				Code:    "RE_400_INVALID_LANGUAGE",
			},
		})
		return
	}

	// Load reasons for the selected language
	reasons, err := loadReasons(lang)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.DefaultResponse{
			Payload: nil,
			Status: response.StatusResponse{
				Status:  http.StatusInternalServerError,
				Message: "Failed to load reasons",
				Code:    "RE_500_LOAD_ERROR",
			},
		})
		return
	}

	// Select a random reason
	randomReason := reasons[rand.Intn(len(reasons))]

	ctx.JSON(http.StatusOK, response.DefaultResponse{
		Payload: map[string]string{
			"reason": randomReason,
		},
		Status: response.StatusResponse{
			Status:  http.StatusOK,
			Message: "Reason delivered successfully.",
			Code:    "RE_200_REASON_DELIVERED",
		},
	})
}

// loadReasons loads reasons from JSON file for the given language
func loadReasons(lang string) ([]string, error) {
	// Check cache first
	cacheMutex.RLock()
	if reasons, ok := reasonCache[lang]; ok {
		cacheMutex.RUnlock()
		return reasons, nil
	}
	cacheMutex.RUnlock()

	// Load from file
	filePath := fmt.Sprintf("data/reasons.%s.json", lang)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var reasons []string
	if err := json.Unmarshal(data, &reasons); err != nil {
		return nil, err
	}

	// Store in cache
	cacheMutex.Lock()
	reasonCache[lang] = reasons
	cacheMutex.Unlock()

	return reasons, nil
}

// getAvailableLanguages scans the data directory for available language files
func getAvailableLanguages() []string {
	langMutex.Do(func() {
		files, err := filepath.Glob("data/reasons.*.json")
		if err != nil {
			return
		}

		for _, file := range files {
			// Extract language code from filename (e.g., "data/reasons.de.json" -> "de")
			base := filepath.Base(file)
			parts := strings.Split(base, ".")
			if len(parts) == 3 {
				availableLanguages = append(availableLanguages, parts[1])
			}
		}
	})

	return availableLanguages
}

// isLanguageSupported checks if the given language is in the list of available languages
func isLanguageSupported(lang string, languages []string) bool {
	return slices.Contains(languages, lang)
}
