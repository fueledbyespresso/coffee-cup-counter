package main

import (
	"coffee-cup-counter/commands"
	"coffee-cup-counter/database"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"net/http"
	"os"
)

//Load the environment variables from the projectvars.env file
func initEnv() {
	if _, err := os.Stat(".env"); err == nil {
		err = godotenv.Load(".env")
		if err != nil {
			fmt.Println("Error loading environment.env")
		}
		fmt.Println("Current environment:", os.Getenv("ENV"))
	}
}

//Force SSL in Heroku
func forceSSL() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Header.Get("x-forwarded-proto") != "https" {
			sslUrl := "https://" + c.Request.Host + c.Request.RequestURI
			c.Redirect(http.StatusTemporaryRedirect, sslUrl)
			return
		}
		c.Next()
	}
}

func createServer(dbConnection *database.DB) *gin.Engine {
	r := gin.Default()
	//r.Use(gzip.Gzip(gzip.DefaultCompression))
	if os.Getenv("ENV") != "DEV" {
		r.Use(forceSSL())
	}
	r.POST("/tally", commands.VerifySlackRequest(), commands.Tally(dbConnection))
	return r
}

func main() {
	initEnv()
	database.PerformMigrations("file://database/migrations")
	db := database.InitDBConnection()
	defer db.Close()
	// Run a background goroutine to clean up expired sessions from the database.
	dbConnection := &database.DB{Db: db}
	r := createServer(dbConnection)

	err := r.Run()
	if err != nil {
		return
	}
}
