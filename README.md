# cast

[![Build Status](https://travis-ci.org/xiaojiaoyu100/cast.svg?branch=master)](https://travis-ci.org/xiaojiaoyu100/cast)
[![Go Report Card](https://goreportcard.com/badge/github.com/xiaojiaoyu100/cast)](https://goreportcard.com/report/github.com/xiaojiaoyu100/cast)
[![GoDoc](https://godoc.org/github.com/xiaojiaoyu100/cast?status.svg)](https://godoc.org/github.com/xiaojiaoyu100/cast)

Cast is a http request library written in Golang.

This project is ready for production use and the master branch is always stable. But the API may be broken in the future release.

## Getting started

    dep ensure -add github.com/xiaojiaoyu100/cast
    
## Usage

### Generate a Cast

```go
c, err := cast.New(cast.WithBaseURL("https://status.github.com"))
```

### Generate a request

```go
request := c.NewRequest()
```

### Get


```go
request := c.NewRequest().Get().WithPath("/api.json")
response, err := c.Do(request)
```

### POST X-WWW-FORM-URLENCODED

```go
request := c.NewRequest().Get().WithPath("/api.json").WithFormURLEncodedBody(body)
resp, err := c.Do(request)
```

### POST JSON 

```go
request := c.NewRequest().Post().WithPath("/api.json").WithJSONBody(body)
response, err := c.Do(request)
```

### POST XML

```go
request := c.NewRequest().Post().WithPath("/api.json").WithXMLBody(body)
response, err := c.Do(request)
```

### POST MULTIPART FORM DATA

```go
request := c.NewRequest().Post().WithPath("/api.json").WithMultipartFormDataBody(formData)
resp, err := c.Do(request)
```

### Timeout

```go
c.NewRequest().WithTimeout(3 * time.Second)
```

### Retry

```go
cast.WithRetry(3)
```

### Backoff

```go
cast.WithXXXBackoffStrategy()
```

## License

[MIT License](LICENSE)



