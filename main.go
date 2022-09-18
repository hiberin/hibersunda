package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type Words struct {
	Words []word `json:"words"`
}

type word struct {
	Sorangan string `json:"sorangan"`
	Batur    string `json:"batur"`
	Loma     string `json:"loma"`
	Bindo    string `json:"bindo"`
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
	jsonFile, err := os.Open("./undakUsukBasa.json")

	// if error happens
	if err != nil {
		fmt.Println(err)
	}

	// If no error
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	// Save json content into byte
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// Initiate interface
	var words Words
	json.Unmarshal(byteValue, &words)

	// return json
	c.IndentedJSON(http.StatusOK, words)
}

/**
* Function to get words containing suffix
 */
func getWordsBySubstring(c *gin.Context) {
	substring := c.Param("substring")
	jsonFile, err := os.Open("./undakUsukBasa.json")

	// if error happens
	if err != nil {
		fmt.Println(err)
	}

	// If no error
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	// Save json content into byte
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// Initiate interface
	var words Words
	json.Unmarshal(byteValue, &words)

	foundWords := []interface{}{}
	for i := 0; i < len(words.Words); i++ {
		soranganWord := strings.ToLower(words.Words[i].Sorangan)
		baturWord := strings.ToLower(words.Words[i].Batur)
		lomaWord := strings.ToLower(words.Words[i].Loma)
		bindoWord := strings.ToLower(words.Words[i].Bindo)
		substring = strings.ToLower(substring)
		if strings.Contains(soranganWord, substring) || strings.Contains(baturWord, substring) || strings.Contains(lomaWord, substring) || strings.Contains(bindoWord, substring) {
			foundWords = append(foundWords, words.Words[i])
		}
	}
	c.IndentedJSON(http.StatusOK, foundWords)
	return
}
