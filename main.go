package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/revrost/go-openrouter"
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
	openRouterToken := os.Getenv("HACK_CLUB_API_KEY")

	openRouterClient := openrouter.NewClient(openRouterToken)


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
				cmdInput := cmd.Text
				println(cmdInput)
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
				case "/d4lm-ask":
					resp, _ := openRouterClient.CreateChatCompletion(
						context.Background(),
						openrouter.ChatCompletionRequest{
							Model: "~google/gemini-flash-latest",
							Messages: []openrouter.ChatCompletionMessage{
								openrouter.UserMessage(cmdInput),
							},
						})

					payload := &slack.TextBlockObject{
						Type: slack.MarkdownType,
						Text: resp.Choices[0].Message.Content.Text,
					}
					client.Ack(*evt.Request, payload)
				}
			}
		}
	}()

	client.Run()
}
