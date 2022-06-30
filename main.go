package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/slack-go/slack"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

var signingSecret string
var scoreboard map[string]int

func main() {
	r := gin.Default()
	r.POST("/tally", tally())
	r.POST("/scoreboard", scores())

	signingSecret = os.Getenv("SIGNING_SECRET")
	if signingSecret == "" {
		return
	}
	scoreboard = make(map[string]int)

	err := r.Run()
	if err != nil {
		return
	}
}

func tally() gin.HandlerFunc {
	return func(c *gin.Context) {
		w := c.Writer
		r := c.Request
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

		switch s.Command {
		case "/tally":
			fmt.Println("Tally case")
			scoreboard[s.UserID]++
			str := fmt.Sprintf("You had: %d", scoreboard[s.UserID])
			params := &slack.Msg{Text: str}
			c.JSON(200, params)
		default:
			fmt.Println("Incorrect command")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func scores() gin.HandlerFunc {
	return func(c *gin.Context) {
		w := c.Writer
		r := c.Request
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

		switch s.Command {
		case "/scoreboard":
			fmt.Println("Scoreboard case")
			str := fmt.Sprintf("You had: %d", scoreboard[s.UserID])

			params := &slack.Msg{Text: str}
			c.JSON(200, params)
		default:
			fmt.Println("Incorrect command")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
