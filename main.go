package main

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type TextData struct {
	Text string `json:"text"`
}

func main() {
	r := gin.Default()
	r.Use(cors.Default())

	/*r.GET("/api/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})*/

	r.POST("/api/process", func(c *gin.Context) {

		var data TextData
		if err := c.BindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		processedText := processText(data.Text)
		c.JSON(http.StatusOK, gin.H{"text": processedText})

	})

	r.Run(":8080")
}

func processText(text string) string {
	fmt.Println("Texto Obtenido: " + text)
	return "Texto Obtenido: " + text
}
