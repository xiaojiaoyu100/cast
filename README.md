# cast [![Build Status](https://travis-ci.org/xiaojiaoyu100/cast.svg?branch=master)](https://travis-ci.org/xiaojiaoyu100/cast)

Cast is a http request library written in Golang.

This project is ready for production use and the master branch is always stable. But the API may be broken in the future release.

## Features

+ Use functional options to provide clean constructor.
+ Support url path params. 
+ Encode struct into url query params.
+ Encode struct into http.Request body, including JSON and x-www-form-urlencoded data.
+ Support retry backoff strategies and provide retry hooks to customize retry conditions with respect to http.Response
+ Support timeout http.Request



## Getting started

    go get github.com/xiaojiaoyu100/cast
    
## Usage

### Get

```go
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
```

## License

[MIT License](LICENSE)


