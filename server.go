package main

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/johnforster/racetrack-go/circuit"
	"github.com/johnforster/racetrack-go/track_analyser"
	"github.com/sergeymakinen/go-bmp"
)

func main() {
	port := os.Getenv("PORT")
	frontend_url := os.Getenv("FRONTEND_URL")
	r := gin.Default()
	r.LoadHTMLGlob("static/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
	r.POST("/tracks", postTracks)

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{frontend_url}

	r.Use(cors.New(config))

	r.Run("localhost:" + port)
}

type FormSubmissionJSON struct {
	Image *multipart.FileHeader `form:"trackimage" binding:"required"`
	// TrackInfo   struct {
	// 	TrackTitle    string `json:"title" binding:"required,min=4,max=32"`
	// 	Author string `json:"author" binding:"required,max=20"`
	// } `form:"user" binding:"required"`
}

func postTracks(c *gin.Context) {
	file, header, err := c.Request.FormFile("trackimage")

	if err != nil {
		fmt.Println("Error opening file")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	fmt.Println("Received file with filename:" + header.Filename)

	image, err := bmp.Decode(file)
	if err != nil {
		fmt.Println("Error decoding image")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tracks, err := track_analyser.GetTracksFromImage(image)

	if err != nil {
		fmt.Println("Error getting tracks from image")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(tracks) != 2 {
		fmt.Println("Error getting tracks from image")
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%v tracks found", len(tracks))})
		return
	}

	circuit := circuit.NewCircuit(tracks[1], tracks[0])

	json, _ := circuit.ToJSON()

	c.Data(http.StatusOK, "application/json", json)
}
