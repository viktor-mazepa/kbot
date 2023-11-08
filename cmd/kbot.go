/*
Copyright © 2023 VIKTOR MAZEPA <viktor.mazepa@gmail.com>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
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
			switch payload {
			case "hello":
				err = ctx.Send(fmt.Sprintf("Hello I'm Kbot %s", appVersion))
			}
			return err
		})

		kbot.Start()
	},
}

func init() {
	rootCmd.AddCommand(kbotCmd)
}
