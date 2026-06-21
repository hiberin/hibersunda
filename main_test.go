package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetWordsBySubstringMatchesAccentInsensitiveSlug(t *testing.T) {
	tests := []struct {
		name          string
		slug          string
		expectedField string
		expectedValue string
	}{
		{
			name:          "plain e matches e acute",
			slug:          "sae",
			expectedField: "sorangan",
			expectedValue: "saé",
		},
		{
			name:          "plain slug matches accented word",
			slug:          "kangge",
			expectedField: "sorangan",
			expectedValue: "kanggo, kanggé",
		},
		{
			name:          "existing plain substring still matches",
			slug:          "ajar",
			expectedField: "sorangan",
			expectedValue: "ajar",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			words := getWordsBySubstringResponse(t, tt.slug)

			for _, w := range words {
				if wordField(w, tt.expectedField) == tt.expectedValue {
					return
				}
			}

			t.Fatalf("expected response for slug %q to include %s=%q, got %#v", tt.slug, tt.expectedField, tt.expectedValue, words)
		})
	}
}

func getWordsBySubstringResponse(t *testing.T, slug string) []word {
	t.Helper()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/undakusukbasa/:substring", getWordsBySubstring)

	req := httptest.NewRequest(http.MethodGet, "/undakusukbasa/"+slug, nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d with body %s", recorder.Code, recorder.Body.String())
	}

	var words []word
	if err := json.Unmarshal(recorder.Body.Bytes(), &words); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	return words
}

func wordField(w word, field string) string {
	switch field {
	case "sorangan":
		return w.Sorangan
	case "batur":
		return w.Batur
	case "loma":
		return w.Loma
	case "bindo":
		return w.Bindo
	case "english":
		return w.English
	default:
		return ""
	}
}
