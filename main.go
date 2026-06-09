package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	appToken := os.Getenv("SLACK_APP_TOKEN")
	botToken := os.Getenv("SLACK_BOT_TOKEN")

	// Following the example for the API...
	api := slack.New(botToken, slack.OptionAppLevelToken(appToken))

	client := socketmode.New(api)

	go func() {
		for evt := range client.Events {
			switch evt.Type {
			case socketmode.EventTypeConnecting:
				fmt.Println("Connecting...")
			case socketmode.EventTypeConnectionError:
				fmt.Println("Connection Error...")
			case socketmode.EventTypeConnected:
				fmt.Println("Connected...")
			case socketmode.EventTypeSlashCommand:
				cmd, ok := evt.Data.(slack.SlashCommand)
				cmdName := cmd.Command
				fmt.Println(cmdName)
				if !ok {
					fmt.Println("Ignored ", evt)
					continue
				}

				switch cmdName {
				case "/d4lm-ping":
					payload := &slack.TextBlockObject{
						Type: slack.MarkdownType,
						Text: "Ping!",
					}
					client.Ack(*evt.Request, payload)
				case "/"
				}
			}
		}
	}()

	client.Run()
}
