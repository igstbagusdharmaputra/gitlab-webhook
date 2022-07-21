package main

import (
	"fmt"
	"os"
)

const (
	telegramSendMessageURL = "https://api.telegram.org/bot%s/sendMessage?chat_id=%s"
)

type TelegramOptions struct {
	IDToken string
	ChatID  string
}
type sendMessageReqBody struct {
	Text string `json:"text"`
}

func getSendMessageURL() string {
	var telegramOptions = TelegramOptions{
		IDToken: os.Getenv("TOKEN_TELEGRAM"),
		ChatID:  os.Getenv("CHAT_ID_TELEGRAM"),
	}
	return fmt.Sprintf(telegramSendMessageURL, telegramOptions.IDToken, telegramOptions.ChatID)
}
