package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// photo represents data about a photh.
type photo struct {
	ImageData        string `json:"imageData"`
	ImageType        string `json:"imageType"`
	Greeting         string `json:"greeting"`
	DateFormatString string `json:"dateFormatString"`
	Language         string `json:"language"`
	Location         string `json:"location"`
}

type picturetext struct {
	Id          int    `json:"id"`
	PictureText string `json:"picturetext"`
}

func main() {
	router := gin.Default()

	router.Use(CORSMiddleware())
	router.POST("/overlayImage", overlayImage)
	router.Run("0.0.0.0:8080")

}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func postToOverlayImage(photoToBeModified photo) photo {

	pt := getPictureText()
	photoToBeModified.Greeting = pt.PictureText

	json_data, err := json.Marshal(photoToBeModified)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Calling API...")
	// Get URL from environment variable
	overlayURL := os.Getenv("OVERLAY_IMAGE_URL")
	resp, err := http.Post(overlayURL, "application/json",
		bytes.NewBuffer(json_data))

	if err != nil {
		fmt.Println("Error in function 'postToOverlayImage'")
		log.Fatal(err)
	}

	var res map[string]interface{}

	var mp photo
	err = json.NewDecoder(resp.Body).Decode(&mp)
	if err != nil {
		panic(err)
	}

	json.NewDecoder(resp.Body).Decode(&res)

	var moddedPhoto photo
	err = json.Unmarshal([]byte(json_data), &moddedPhoto)
	if err != nil {
		log.Fatal(err)
	}

	return mp
}

func overlayImage(c *gin.Context) {

	var p photo
	if err := c.BindJSON(&p); err != nil {
		return
	}

	var newlyModifiedPhoto photo
	newlyModifiedPhoto = postToOverlayImage(p)
	c.IndentedJSON(http.StatusCreated, newlyModifiedPhoto)
}

func getPictureText() picturetext {
	fmt.Println("Calling API...")
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://getpicturetext-rhn-engineering-dsch-dev.apps.sandbox-m3.1530.p1.openshiftapps.com/text", nil)
	if err != nil {
		fmt.Print(err.Error())
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Print(err.Error())
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
	}
	var pt picturetext
	json.Unmarshal(bodyBytes, &pt)
	return pt
}
