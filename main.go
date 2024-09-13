package main

import (
	analyzer "MIA_2S_P1_202200041/analyzer"
	"fmt"
	"net/http"
	"strings"

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
		//analyzer.Analyzer(data.Text)
		c.JSON(http.StatusOK, gin.H{"text": processedText})

	})

	r.Run(":8080")
}

func processText(text string) string {
	lineas := strings.Split(text, "\n")
	var contenido []string

	for _, line := range lineas {
		if line == "" {
			continue
		}
		result, err := analyzer.Analyzer(line)

		if err != nil {
			contenido = append(contenido, fmt.Sprintf("Error: %v", err))
		} else {
			contenido = append(contenido, fmt.Sprintf("Result: %v", result))
		}

	}

	//analyzer.Analyzer(text)
	return "Texto Obtenido: " + text
}
