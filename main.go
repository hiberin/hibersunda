package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// inisiasi router make gin
	router := gin.Default()
	// router get pikeun nyokot kabeh kata/kecap nu aya
	router.GET("/kabehkecap", cokotKabehKecap)

	// ngajalankeun router dina port 8080
	router.Run("localhost: 8080")
}

// Fungsi pikeun nyokot kabeh kecap anu aya dina daptarKecap,json
func cokotKabehKecap(c *gin.Context) {
	jsonFile, err := os.Open("daptarKecap.json")

	// lamun aya error
	if err != nil {
		fmt.Println(err)
	}

	//* Lamun euweuh error
	// Teundeun value tina jsonFile dina bentuk byte
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// Inisiasi interface hasil
	var hasil map[string]interface{}

	// parsing json kana variable hasil
	json.Unmarshal([]byte(byteValue), &hasil)

	// Pulangkeun eusi tina file json
	c.IndentedJSON(http.StatusOK, hasil["daptarKecap"])
}
