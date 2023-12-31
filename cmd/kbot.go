/*
Copyright © 2023 VIKTOR MAZEPA <viktor.mazepa@gmail.com>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	telebot "gopkg.in/telebot.v3"
)

var (
	TeleToken = os.Getenv("TELE_TOKEN")
)

var kbotCmd = &cobra.Command{
	Use:     "kbot",
	Aliases: []string{"start"},
	Short:   "Start telegram kBot application",
	Long: `A simple Telegram bot that can handle text messages.
	You can write something to https://t.me/viktormazepa_bot and sometimes it answer:)`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("kbot %s started\n", appVersion)

		kbot, err := telebot.NewBot(telebot.Settings{
			URL:    "",
			Token:  TeleToken,
			Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
		})

		if err != nil {
			log.Fatalf("please check TELE_TOKEN env variable, %s", err)
		}

		kbot.Handle(telebot.OnText, func(ctx telebot.Context) error {
			log.Println(ctx.Message().Payload, ctx.Text())
			payload := ctx.Message().Payload
			answerStr := handlePayload(payload)
			err = ctx.Send(answerStr)
			return err
		})

		kbot.Handle(telebot.OnVoice, func(ctx telebot.Context) error {
			answerStr := "I don' have ears, so I can't hear you :)"
			err = ctx.Send(answerStr)
			return err
		})

		kbot.Handle(telebot.OnPhoto, func(ctx telebot.Context) error {
			answerStr := "Nice picture... or not"
			err = ctx.Send(answerStr)
			return err
		})

		kbot.Handle(telebot.OnSticker, func(ctx telebot.Context) error {
			answerStr := "It's ver funny... actually no 😠"
			err = ctx.Send(answerStr)
			return err
		})

		kbot.Start()
	},
}

func handlePayload(payload string) string {
	switch strings.ToLower(payload) {
	case "hello":
		return fmt.Sprintf("Hello I'm Kbot %s", appVersion)
	case "version":
		return fmt.Sprintf("Current kBot version %s", appVersion)
	case "how are you?":
		return "Thank you, I'm fine."
	case "time":
		dt := time.Now()
		return dt.Format("02-01-2006 15:04:05")
	default:
		return "Unknown request"
	}
}

func init() {
	rootCmd.AddCommand(kbotCmd)
}
