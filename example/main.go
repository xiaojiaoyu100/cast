package main

import (
	"log"

	"github.com/xiaojiaoyu100/cast"
	"time"
)

func main() {
	baseUrl := "https://status.github.com"
	// Get
	c := cast.New(
		cast.WithBaseUrl(baseUrl),
	)
	reply, err := c.Get().WithPath("/api.json").WithTimeout(10 * time.Millisecond).Request()
	if err != nil {
		log.Fatalln(err)
	}
	var ApiUrl struct {
		StatusUrl      string `json:"status_url"`
		MessagesUrl    string `json:"messages_url"`
		LastMessageUrl string `json:"last_message_url"`
		DailySummary   string `json:"daily_summary"`
	}
	log.Println(string(reply.Body()), reply.Url())
	if !reply.StatusOk() {
		return
	}
	if err := reply.DecodeFromJson(&ApiUrl); err != nil {
		log.Fatalln(err)
	}
}
