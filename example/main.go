package main

import (
	"github.com/xiaojiaoyu100/cast"
	"fmt"
)

func main() {
	baseUrl := "https://status.github.com"
	c, err := cast.New(cast.WithBaseURL(baseUrl))
	if err != nil {
		return
	}
	request := c.NewRequest().Get().WithPath("/api.json")
	resp, err := c.Do(request)
	if err != nil {
		fmt.Println(err)
		return
	}
	var ApiUrl struct {
		StatusUrl      string `json:"status_url"`
		MessagesUrl    string `json:"messages_url"`
		LastMessageUrl string `json:"last_message_url"`
		DailySummary   string `json:"daily_summary"`
	}
	fmt.Println(resp.String())
	if !resp.Success() {
		return
	}
	if err := resp.DecodeFromJSON(&ApiUrl); err != nil {
		fmt.Println(err)
	}
}
