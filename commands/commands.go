package commands

import (
	"coffee-cup-counter/database"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/slack-go/slack"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func VerifySlackRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		w := c.Writer
		r := c.Request
		signingSecret := os.Getenv("SIGNING_SECRET")
		verifier, err := slack.NewSecretsVerifier(r.Header, signingSecret)
		if err != nil {
			fmt.Println("Step 1", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		r.Body = ioutil.NopCloser(io.TeeReader(r.Body, &verifier))
		s, err := slack.SlashCommandParse(r)
		if err != nil {
			fmt.Println("Step 2", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err = verifier.Ensure(); err != nil {
			fmt.Println("Step 3", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		c.Set("SlackCommand", s)
	}
}

func Tally(dbConnection *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		sStr, exists := c.Get("SlackCommand")
		if !exists {
			c.AbortWithStatusJSON(500, "Could not verify Slack Request")
			return
		}
		s := sStr.(slack.SlashCommand)
		switch s.Command {
		case "/tally":
			fmt.Println("Tally case")

			params := &slack.Msg{Text: "TALLIED!"}
			c.JSON(200, params)
		default:
			fmt.Println("Incorrect command")
			c.AbortWithStatus(500)
			return
		}
	}
}
