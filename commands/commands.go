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
	"strconv"
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

func JoinContest(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		sStr, exists := c.Get("SlackCommand")
		if !exists {
			c.AbortWithStatusJSON(500, "Could not verify Slack Request")
			return
		}
		s := sStr.(slack.SlashCommand)

		row := db.Db.QueryRow(`INSERT INTO contestant (username, user_id) 
									 	VALUES ($1, $2) on conflict do nothing`, s.UserName, s.UserID)
		if row.Err() != nil {
			c.AbortWithStatusJSON(400, "Invalid parameters.")
			return
		}

		params := &slack.Msg{Text: "You have joined the contest! If you were already in the competition, nothing changed."}
		c.JSON(200, params)
	}
}

func ListMembers(dbConnection *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		sStr, exists := c.Get("SlackCommand")
		if !exists {
			c.AbortWithStatusJSON(500, "Could not verify Slack Request")
			return
		}
		s := sStr.(slack.SlashCommand)

		temp := slack.GetUsersInConversationParameters{
			ChannelID: s.ChannelID,
			Cursor:    "",
			Limit:     0,
		}
		conversation, _, err := slack.New(os.Getenv("BOT_TOKEN")).GetUsersInConversation(&temp)
		if err != nil {
			return
		}

		for i, s := range conversation {
			fmt.Println(i, s)
		}

		params := &slack.Msg{Text: s.Text}
		c.JSON(200, params)
	}
}

func Tally(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		sStr, exists := c.Get("SlackCommand")
		if !exists {
			c.AbortWithStatusJSON(500, "Could not verify Slack Request")
			return
		}
		s := sStr.(slack.SlashCommand)
		_, err := db.Db.Query(`INSERT INTO tally (contestant) VALUES ($1)`, s.UserID)
		if err != nil {
			c.AbortWithStatusJSON(500, &slack.Msg{Text: "Cannot tally!"})
			return
		}

		row := db.Db.QueryRow(`SELECT count() FROM tally WHERE contestant=$1`, s.UserID)
		tally := 0
		if row.Err() != nil {
			c.AbortWithStatusJSON(500, &slack.Msg{Text: "Cannot display tally!"})
			return
		}

		err = row.Scan(&tally)

		params := &slack.Msg{Text: strconv.Itoa(tally)}
		c.JSON(200, params)
	}
}

func Scoreboard(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Db.Query(`SELECT username, count(username) FROM tally 
    									JOIN contestant c on c.user_id = tally.contestant 
                                 		group by username`)
		if err != nil {
			c.AbortWithStatusJSON(500, &slack.Msg{Text: "Cannot display scoreboard!"})
			return
		}
		scoreboardStr := ""
		for rows.Next() {
			var contestant string
			var count int
			err = rows.Scan(&contestant, &count)
			if err != nil {
				c.AbortWithStatusJSON(500, &slack.Msg{Text: "Error retrieving drink tallies."})
				return
			}

			scoreboardStr += fmt.Sprintf("%s drank %d drinks\n", contestant, count)
		}
		params := &slack.Msg{Text: scoreboardStr}
		c.JSON(200, params)
	}
}
