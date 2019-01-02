package main

import (
	"log"

	"github.com/xiaojiaoyu100/cast"
)

func main() {
	baseUrl := "https://status.github.com"
	c := cast.New(cast.WithBaseURL(baseUrl))
	request := c.NewRequest().Get().WithPath("/api.json")
	resp, err := c.Do(request)
	if err != nil {
		log.Fatalln(err)
	}
	var ApiUrl struct {
		StatusUrl      string `json:"status_url"`
		MessagesUrl    string `json:"messages_url"`
		LastMessageUrl string `json:"last_message_url"`
		DailySummary   string `json:"daily_summary"`
	}
	log.Println(string(resp.Body()))
	if !resp.StatusOk() {
		return
	}
	if err := resp.DecodeFromJSON(&ApiUrl); err != nil {
		log.Fatalln(err)
	}
}
