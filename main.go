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

	openRouterConfig := openrouter.DefaultConfig(openRouterToken)

	openRouterConfig.BaseURL = "https://ai.hackclub.com/proxy/v1"

	openRouterClient := openrouter.NewClientWithConfig(*openRouterConfig)

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
					// There is a 3-second timeout for a reason hopefully
					client.Ack(*evt.Request, nil)

					go func(cmd slack.SlashCommand) {
						resp, err := openRouterClient.CreateChatCompletion(
							context.Background(),
							openrouter.ChatCompletionRequest{
								Model: "~google/gemini-flash-latest", // Use a valid model string
								Messages: []openrouter.ChatCompletionMessage{
									openrouter.UserMessage(cmd.Text),
								},
							})

						if err != nil {
							fmt.Println("LLM Error:", err)
							api.PostEphemeral(cmd.ChannelID, cmd.UserID, slack.MsgOptionText("Sorry, I had trouble thinking: "+err.Error(), false))
							return
						}

						msgText := "Response: \n" + resp.Choices[0].Message.Content.Text
						api.PostMessage(cmd.ChannelID, slack.MsgOptionText(msgText, false))
					}(cmd)
				}
			}
		}
	}()

	client.Run()
}
