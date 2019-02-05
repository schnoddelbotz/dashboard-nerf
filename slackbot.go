package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/nlopes/slack"
)

var botUserID string

// based on https://github.com/nlopes/slack/blob/master/examples/websocket/websocket.go
func doSlack() {
	api := slack.New(
		slackToken,
		slack.OptionDebug(false),
		slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)),
	)

	authTest, err := api.AuthTest()
	if err != nil {
		log.Fatalf("Error getting channels: %s\n", err)
	}
	botUserID = authTest.UserID

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {

		case *slack.ConnectingEvent:
			fmt.Println("Slack: Connecting...")

		case *slack.ConnectedEvent:
			fmt.Println("Slack: Successfully connected")

		case *slack.MessageEvent:
			response := handleSlackMessage(msg.Data.(*slack.MessageEvent).Msg.Channel[0], msg.Data.(*slack.MessageEvent).Msg.Text)
			if response != "" {
				rtm.SendMessage(rtm.NewOutgoingMessage(response, msg.Data.(*slack.MessageEvent).Msg.Channel))
			}

		case *slack.RTMError:
			fmt.Printf("Slack Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			fmt.Println("Slack Error: Invalid credentials")
			return

		default:
			// Ignore other events..
		}
	}
}

func handleSlackMessage(mtype byte, msg string) string {
	// https://stackoverflow.com/questions/41111227/how-can-a-slack-bot-detect-a-direct-message-vs-a-message-in-a-channel
	if string(mtype) == "D" || strings.HasPrefix(msg, "<@"+botUserID+">") {
		msg = strings.TrimPrefix(msg, "<@"+botUserID+"> ")
		if msg == "stop" {
			if cmd != nil {
				cmd.Process.Kill()
				return "Okay, okay... stopped it."
			}
			return "Nothing to stop. Try to `play` something first!"
		} else if strings.HasPrefix(msg, "play") {
			requestedGlob := strings.TrimPrefix(msg, "play ")
			glob := mediaRoot + "/*" + requestedGlob + "*"
			candidates, _ := filepath.Glob(glob)
			if len(candidates) == 0 {
				return "No matches. Try `play *` to list all media I can play."
			} else if len(candidates) == 1 {
				mType := getMediaType(candidates[0])
				playQueue <- playRequest{Filename: candidates[0], MediaType: mType}
				return fmt.Sprintf("Playing %s for you!", candidates[0])
			} else {
				var baseNames []string
				for _, f := range candidates {
					baseNames = append(baseNames, filepath.Base(f))
				}
				return fmt.Sprintf("Multiple candidates: %s", baseNames)
			}
		}
		return "huh? I only understand: `play GLOB-PATTERN` or `stop`."
	}
	return ""
}
