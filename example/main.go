package main

import (
	"fmt"
	"time"

	"github.com/xiaojiaoyu100/cast"
)

func retry(response *cast.Response, err error) bool {
	if err != nil {
		return true
	}
	if !response.StatusOk() {
		return true
	}
	return false
}

func main() {
	baseUrl := "https://status.github.com"
	c, err := cast.New(
		cast.WithBaseURL(baseUrl),
		cast.WithRetry(3),
		cast.AddRetryHooks(retry),
		cast.WithExponentialBackoffDecorrelatedJitterStrategy(
			time.Millisecond*200,
			time.Millisecond*450,
		),
		cast.AddCircuitConfig("name1"),
		cast.AddCircuitConfig("name2"),
		cast.WithDefaultCircuit("name1"),
	)
	if err != nil {
		return
	}
	for {
		time.Sleep(20 * time.Millisecond)
		request := c.NewRequest().Get().WithPath("/api.json").WithCircuit("name2")
		resp, err := c.Do(request)
		if err != nil {
			fmt.Println(err)
			continue
		}
		var ApiUrl struct {
			StatusUrl      string `json:"status_url"`
			MessagesUrl    string `json:"messages_url"`
			LastMessageUrl string `json:"last_message_url"`
			DailySummary   string `json:"daily_summary"`
		}
		fmt.Println(resp.String())
		if !resp.Success() {
			continue
		}
		if err := resp.DecodeFromJSON(&ApiUrl); err != nil {
			fmt.Println(err)
		}
	}
}
