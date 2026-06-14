package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
	"google.golang.org/genai"
)

const model = "gemini-3.1-flash-lite"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	appToken := os.Getenv("SLACK_APP_TOKEN")
	botToken := os.Getenv("SLACK_BOT_TOKEN")

	geminiClient, err := genai.NewClient(context.Background(), &genai.ClientConfig{})

	if err != nil {
		log.Fatal(err)
	}

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
					response, err := geminiClient.Models.GenerateContent(
						context.Background(),
						model,
						genai.Text("Very quickly (less than 3 seconds), answer this question: "+cmd.Text),
						nil,
					)

					if err != nil {
						log.Fatal(err)
					}

					payload := &slack.TextBlockObject{
						Type: slack.MarkdownType,
						Text: "Please note that this command is quickly processed: \n" + response.Text(),
					}
					client.Ack(*evt.Request, payload)

				case "/d4lm-ascii-art":
					response, err := geminiClient.Models.GenerateContent(
						context.Background(),
						model,
						genai.Text("very quickly (less than 3 seconds), generate ascii art for this user following the prompt: "+cmd.Text),
						nil)

					if err != nil {
						log.Fatal(err)
					}

					payload := &slack.TextBlockObject{
						Type: slack.MarkdownType,
						Text: "Please note that this command is quickly processed, here is your ASCII Art: \n" + response.Text(),
					}

					client.Ack(*evt.Request, payload)

				case "/d4lm-catfact":
					response, err := geminiClient.Models.GenerateContent(
						context.Background(),
						model,
						genai.Text("In less than 3 seconds, send me a catfact!"),
						nil)

					if err != nil {
						log.Fatal(err)
					}

					payload := &slack.TextBlockObject{
						Type: slack.MarkdownType,
						Text: "Please note that this command is quickly processed: \n" + response.Text(),
					}

					client.Ack(*evt.Request, payload)

				case "/d4lm-joke":
					response, err := geminiClient.Models.GenerateContent(
						context.Background(),
						model,
						genai.Text("In less than 3 seconds, reply with a joke!"),
						nil)

					if err != nil {
						log.Fatal(err)
					}

					payload := &slack.TextBlockObject{
						Type: slack.MarkdownType,
						Text: "Please note that this command is quickly processed: \n" + response.Text(),
					}

					client.Ack(*evt.Request, payload)

				case "/d4lm-help":
					payload := &slack.TextBlockObject{
						Type: slack.MarkdownType,
						Text: "/d4lm-ping			Pings the bot!\n" +
							"/d4lm-help			Shows you all the commands!\n" +
							"/d4lm-catfact		Replies with a cat fact!\n" +
							"/d4lm-joke			Replies with a joke!\n" +
							"/d4lm-ask			Replies with the answer you seek!\n" +
							"/d4lm-ascii-art		Gnerates ascii art!",
					}

					client.Ack(*evt.Request, payload)
				}
			}
		}
	}()
	client.Run()
}
