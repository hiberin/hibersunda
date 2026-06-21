package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
	"golang.org/x/text/unicode/norm"
)

//go:embed undakUsukBasa.json
var wordsJSON []byte

type Words struct {
	Words []word `json:"words"`
}

type word struct {
	Sorangan string `json:"sorangan"`
	Batur    string `json:"batur"`
	Loma     string `json:"loma"`
	Bindo    string `json:"bindo"`
	English  string `json:"english"`
}

// CORS Middleware
func CORS(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "GET")
	c.Header("Content-Type", "application/json")
}

func main() {
	// initiate router using gin
	router := gin.Default()
	// We use our custom CORS Middleware
	router.Use(CORS)
	router.GET("/undakusukbasa", getAllWords)
	router.GET("/undakusukbasa/:substring", getWordsBySubstring)

	// run the router on port 8080
	router.Run()
}

/**
* Function to get all the words
 */
func getAllWords(c *gin.Context) {
	words, err := loadWords()
	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "failed to load words"})
		return
	}

	// return json
	c.IndentedJSON(http.StatusOK, words)
}

/**
* Function to get words containing suffix
 */
func getWordsBySubstring(c *gin.Context) {
	substring := normalizeForSearch(c.Param("substring"))
	words, err := loadWords()
	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "failed to load words"})
		return
	}

	foundWords := []interface{}{}
	for i := 0; i < len(words.Words); i++ {
		if wordContainsSubstring(words.Words[i], substring) {
			foundWords = append(foundWords, words.Words[i])
		}
	}
	c.IndentedJSON(http.StatusOK, foundWords)
	return
}

func wordContainsSubstring(w word, substring string) bool {
	return strings.Contains(normalizeForSearch(w.Sorangan), substring) ||
		strings.Contains(normalizeForSearch(w.Batur), substring) ||
		strings.Contains(normalizeForSearch(w.Loma), substring) ||
		strings.Contains(normalizeForSearch(w.Bindo), substring) ||
		strings.Contains(normalizeForSearch(w.English), substring)
}

func normalizeForSearch(value string) string {
	decomposed := norm.NFD.String(strings.ToLower(value))
	var builder strings.Builder
	builder.Grow(len(decomposed))

	for _, r := range decomposed {
		if unicode.Is(unicode.Mn, r) {
			continue
		}
		builder.WriteRune(r)
	}

	return builder.String()
}

func loadWords() (Words, error) {
	var words Words
	err := json.Unmarshal(wordsJSON, &words)
	return words, err
}
