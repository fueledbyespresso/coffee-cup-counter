package main

import (
	"encoding/json"
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/slack-go/slack"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

var signingSecret string

func main() {
	r := gin.Default()
	r.POST("/tally", home())

	signingSecret = os.Getenv("SIGNING_SECRET")
	if signingSecret == "" {
		return
	}
	flag.Parse()

	err := r.Run()
	if err != nil {
		return
	}
}

func home() gin.HandlerFunc {
	return func(c *gin.Context) {
		w := c.Writer
		r := c.Request
		verifier, err := slack.NewSecretsVerifier(r.Header, signingSecret)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		r.Body = ioutil.NopCloser(io.TeeReader(r.Body, &verifier))
		s, err := slack.SlashCommandParse(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err = verifier.Ensure(); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		switch s.Command {
		case "/echo":
			params := &slack.Msg{Text: s.Text}
			b, err := json.Marshal(params)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_, err = w.Write(b)
			if err != nil {
				return
			}
		default:
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
