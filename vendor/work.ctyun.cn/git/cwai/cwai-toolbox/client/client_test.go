package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	config := &Config{
		Host:  "http://www.baidu.com",
		Token: "fake-token",
	}
	client := NewClient(config)

	path := "/s"
	params := map[string]string{"wd": "query"}
	// GET http://www.baidu.com/s?wd=query
	content, err := client.Get(path, 0, params)
	assert.Nil(t, err)
	t.Logf("response: %v", string(content)[:15])

	values := url.Values{}
	values.Add("wd", "query")
	content, err = client.Get(path, 0, values)
	assert.Nil(t, err)
}

func TestClientWith(t *testing.T) {
	config := &Config{
		Host:  "http://www.baidu.com",
		Token: "fake-token",
	}
	client := NewClient(config)
	client1 := client.WithHeader("abc", "1")
	client2 := client.WithHeader("def", "2")
	client3 := client2.WithHeader("wwww", "3")
	client3 = client3.WithHeader("wwww", "4")
	fmt.Printf("%#v", client3.Header)
	assert.Len(t, client.Header, 0)
	assert.Len(t, client1.Header, 1)
	assert.Len(t, client2.Header, 1)
	assert.Len(t, client3.Header, 2)
}

func TestClientDoGetWithBody(t *testing.T) {
	config := &Config{
		Host:  "http://www.baidu.com",
		Token: "fake-token",
	}

	client := NewClient(config)
	values := url.Values{}
	values.Add("wd", "query")

	rb, err := client.Do("GET", "/s", 0, values, "body xxxxx")
	assert.Nil(t, err)
	assert.NotEmpty(t, rb)
	t.Logf("Receive body : %s", string(rb))
}

func TestClientDoPostWithNilBody(t *testing.T) {
	config := &Config{
		Host:  "http://www.baidu.com",
		Token: "fake-token",
	}

	client := NewClient(config)
	values := url.Values{}
	values.Add("wd", "query")

	rb, err := client.Do("POST", "/s", 0, values, nil)
	fmt.Println("content: ", string(rb[:50]))
	assert.Nil(t, err)
	assert.NotEmpty(t, rb)
	t.Logf("Receive body : %s\n", string(rb))
}

func TestClientDoPostWithBody(t *testing.T) {
	config := &Config{
		Host:  "http://www.baidu.com",
		Token: "fake-token",
	}

	client := NewClient(config)
	values := url.Values{}
	values.Add("wd", "query")

	rb, err := client.Do("POST", "/s", 0, values, "test body")
	fmt.Println("content: ", string(rb[:50]))
	assert.Nil(t, err)
	assert.NotEmpty(t, rb)
	t.Logf("Receive body : %s\n", string(rb))
}

func TestClientDoGetContext(t *testing.T) {
	config := &Config{
		Host:  "http://www.baidu.com",
		Token: "fake-token",
	}

	client := NewClient(config)
	values := url.Values{}
	values.Add("wd", "query")

	clientContextFunc(t, func(c context.Context) {
		rb, err := client.DoContext(c, "GET", "/s", nil, values, "body xxxxx")
		assert.Errorf(t, err, err.Error())
		t.Logf("Receive error: %T, %+v", err, err)
		assert.Empty(t, rb)
	})
}

func clientContextFunc(t *testing.T, do func(context.Context)) {
	ctx, cancel := context.WithCancel(context.TODO())
	go do(ctx)

	wait := time.Millisecond * 5
	select {
	case <-ctx.Done():
		t.Logf("Request finished with %+v", ctx.Err())
		cancel()
	case <-time.After(wait):
		t.Logf("Request after %+v, cancel", wait)
		cancel()
	}

	time.Sleep(time.Second)
}

func TestClientContext(t *testing.T) {
	config := &Config{
		Host:  "http://www.baidu.com",
		Token: "fake-token",
	}

	client := NewClient(config)
	values := url.Values{}
	values.Add("wd", "query")

	reqGroups := []func(context.Context){
		func(c context.Context) {
			rb, err := client.GetContext(c, "/s", nil, values)
			assert.Errorf(t, err, err.Error())
			t.Logf("Receive error: %T, %+v", err, err)
			assert.Empty(t, rb)
		},
		func(c context.Context) {
			rb, err := client.PostContext(c, "/s", nil, values, "xxxxx body")
			assert.Errorf(t, err, err.Error())
			t.Logf("Receive error: %T, %+v", err, err)
			assert.Empty(t, rb)
		},
		func(c context.Context) {
			rb, err := client.PutContext(c, "/s", nil, values, "xxxxx body")
			assert.Errorf(t, err, err.Error())
			t.Logf("Receive error: %T, %+v", err, err)
			assert.Empty(t, rb)
		},
		func(c context.Context) {
			rb, err := client.PatchContext(c, "/s", nil, values, "xxxxx body")
			assert.Errorf(t, err, err.Error())
			t.Logf("Receive error: %T, %+v", err, err)
			assert.Empty(t, rb)
		},
		func(c context.Context) {
			rb, err := client.DeleteContext(c, "/s", nil, values)
			assert.Errorf(t, err, err.Error())
			t.Logf("Receive error: %T, %+v", err, err)
			assert.Empty(t, rb)
		},
	}

	for _, f := range reqGroups {
		ff := f
		clientContextFunc(t, ff)
	}
}

func TestClientDoGetContextWithHeader(t *testing.T) {
	config := &Config{
		Host:  "http://www.baidu.com",
		Token: "fake-token",
	}

	client := NewClient(config)
	values := url.Values{}
	values.Add("wd", "query")

	header := http.Header{
		"GML-Test-API":  []string{"v1"},
		"GML-Test-Unit": []string{"context", "get"},
	}

	rb, err := client.DoContext(context.TODO(), "GET", "/s", header, values, "body xxxxx")
	assert.Nil(t, err, err)
	t.Logf("Receive error: %T, %+v", err, err)
	assert.NotEmpty(t, rb)
}
