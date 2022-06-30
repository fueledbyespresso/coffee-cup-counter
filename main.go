package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()
	r.GET("/home", home())
	err := r.Run()
	if err != nil {
		return
	}
}

func home() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"test": "test",
		})
	}
}
