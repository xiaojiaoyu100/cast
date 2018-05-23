package main

import (
	"log"

	"github.com/xiaojiaoyu100/cast"
)

func main() {
	urlPrefix := "https://status.github.com"
	// Get
	c := cast.New(
		cast.WithUrlPrefix(urlPrefix),
	)
	reply, err := c.WithApi("/api.json").Request()
	if err != nil {
		log.Fatalln(err)
	}
	var ApiUrl struct {
		StatusUrl      string `json:"status_url"`
		MessagesUrl    string `json:"messages_url"`
		LastMessageUrl string `json:"last_message_url"`
		DailySummary   string `json:"daily_summary"`
	}
	log.Println(string(reply.Body()))
	if !reply.StatusOk() {
		return
	}
	if err := reply.DecodeFromJson(&ApiUrl); err != nil {
		log.Fatalln(err)
	}
}
