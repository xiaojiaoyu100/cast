package cast

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/google/go-querystring/query"
)

func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

func TestCast_WithApi(t *testing.T) {
	cast := New()
	api := "/check"
	cast.WithApi(api)
	if api != cast.api {
		t.Fatal("unexpected api")
	}
}

func TestCast_WithMethod(t *testing.T) {
	cast := New()
	method := http.MethodPut
	cast.WithMethod(method)

	assert(t, method == cast.method, "unexpected method")
}

func TestCast_AppendHeader(t *testing.T) {
	originalHeader := http.Header{
		"X": []string{"a"},
	}
	cast := New(WithHeader(originalHeader))
	header := http.Header{
		"X": []string{"b"},
		"Z": []string{"c"},
	}
	cast.AppendHeader(header)
	if cast.header["X"][0] != "a" {
		t.Fatal("unexpected AppendHeader")
	}
	if cast.header["X"][1] != "b" {
		t.Fatal("unexpected AppendHeader")
	}
	if cast.header["Z"][0] != "c" {
		t.Fatal("unexpected AppendHeader")
	}
}

func TestCast_SetHeader(t *testing.T) {
	originalHeader := http.Header{
		"X": []string{"a"},
	}
	cast := New(WithHeader(originalHeader))
	header := http.Header{
		"X": []string{"b", "c"},
	}
	cast.SetHeader(header)
	if cast.header["X"][0] != "c" {
		t.Fatal("unexpected SetHeader")
	}
}

func TestCast_WithQueryParam(t *testing.T) {
	cast := New()
	var query struct {
		Code string `url:"code"`
	}
	cast.WithQueryParam(query)

	if !reflect.DeepEqual(query, cast.queryParam) {
		t.Fatal("unexpected queryParam")
	}
}

func TestCast_WithPathParam(t *testing.T) {
	cast := New()
	pathParam := make(map[string]interface{})
	cast.WithPathParam(pathParam)
	if !reflect.DeepEqual(pathParam, cast.pathParam) {
		t.Fatal("unexpected pathParam")
	}
}

func TestCast_WithJsonBody(t *testing.T) {
	cast := New()
	type payload struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}

	var p payload
	p.Code = 0
	p.Msg = "ok"

	cast.WithJsonBody(p)

	body, err := cast.body.Body()
	ok(t, err)

	bytes, err := ioutil.ReadAll(body)
	ok(t, err)

	var b payload
	err = json.Unmarshal(bytes, &b)
	ok(t, err)

	if p.Code != b.Code || p.Msg != b.Msg {
		t.Fatal("unexpected body")
	}
}

func TestCast_WithUrlEncodedFormBody(t *testing.T) {
	cast := New()

	type payload struct {
		Code int    `url:"code"`
		Msg  string `url:"msg"`
	}

	var p payload
	p.Code = 0
	p.Msg = "ok"

	cast.WithUrlEncodedFormBody(p)

	body, err := cast.body.Body()
	if err != nil {
		t.Fatal(err)
	}

	bytes, err := ioutil.ReadAll(body)
	if err != nil {
		t.Fatal(err)
	}

	values, err := query.Values(&p)
	if err != nil {
		t.Fatal(err)
	}

	if values.Encode() != string(bytes) {
		t.Fatal("unexpected return")
	}
}

func TestCast_WithRetry(t *testing.T) {
	cast := New()
	cast.WithRetry(3)
	if 3 != cast.retry {
		t.Fatal("unexpected retry")
	}
}

func ExampleCast_WithLinearBackoffStrategy() {
	cast := New()
	slope := 1 * time.Second
	cast.WithLinearBackoffStrategy(slope)
	for i := 1; i <= 3; i++ {
		fmt.Println(cast.strat.backoff(i))
	}
	// Output:
	// 1s
	// 2s
	// 3s
}

func ExampleCast_WithConstantBackoffStrategy() {
	cast := New()
	cast.WithConstantBackoffStrategy(2 * time.Second)
	for i := 1; i <= 3; i++ {
		fmt.Println(cast.strat.backoff(i))
	}
	// Output:
	// 2s
	// 2s
	// 2s
}

func ExampleCast_WithExponentialBackoffStrategy() {
	cast := New()
	cast.WithExponentialBackoffStrategy(time.Second, 10*time.Second)
	for i := 1; i <= 5; i++ {
		fmt.Println(cast.strat.backoff(i))
	}
	// Output:
	// 2s
	// 4s
	// 8s
	// 10s
	// 10s
}

func BenchmarkCast_WithExponentialBackoffEqualJitterStrategy(b *testing.B) {
	cast := New()
	cast.WithExponentialBackoffEqualJitterStrategy(time.Second, 10*time.Second)
	for i := 0; i <= b.N; i++ {
		b.Log(cast.strat.backoff(i))
	}
}

func BenchmarkCast_WithExponentialBackoffFullJitterStrategy(b *testing.B) {
	cast := New()
	cast.WithExponentialBackoffFullJitterStrategy(time.Second, 10*time.Second)
	for i := 1; i <= 5; i++ {
		b.Log(cast.strat.backoff(i))
	}
}

func BenchmarkCast_WithExponentialBackoffDecorrelatedJitterStrategy(b *testing.B) {
	cast := New()
	cast.WithExponentialBackoffDecorrelatedJitterStrategy(time.Second, 10*time.Second)
	for i := 1; i <= 5; i++ {
		b.Log(cast.strat.backoff(i))
	}
}

func TestCast_AddRetryHooks(t *testing.T) {
	internalServerErrorHook := func(resp *http.Response) error {
		return nil
	}
	tooManyRequestsHook := func(resp *http.Response) error {
		return nil
	}
	cast := New()
	cast.AddRetryHooks(internalServerErrorHook, tooManyRequestsHook)
	if len(cast.retryHooks) != 2 {
		t.Fatal("fail to add AddRetryHooks.")
	}
}

func TestCast_WithTimeout(t *testing.T) {
	cast := New()
	timeout := 3 * time.Second
	cast.WithTimeout(timeout)
	assert(t, timeout == cast.timeout, "unexpected timeout")
}

func TestCast_Request(t *testing.T) {

}
