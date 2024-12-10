package client

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type HTTPTestSuite struct {
	suite.Suite
}

func TestHttpClient(t *testing.T) {
	suite.Run(t, new(HTTPTestSuite))
}

func (hs *HTTPTestSuite) TestClientRequestTimeout() {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hs.T().Logf("Before serve")
		time.Sleep(10 * time.Second)
		hs.T().Logf("After serve")
		fmt.Fprintf(w, "after sleep 360s")
	}))

	defer svr.Close()

	hs.T().Logf("svr: start")
	c := NewHTTPClient(&HTTPConfig{
		Timeout: 1000, // 1s
	})

	hs.T().Logf("ready to get %s", svr.URL)
	_, err := c.Get(svr.URL)
	hs.T().Logf("get %s err: %+v", svr.URL, err)
	hs.Error(err)
	nErr, ok := err.(net.Error)
	hs.True(ok, err)
	hs.True(os.IsTimeout(nErr), err)
	hs.Contains(nErr.Error(), "(Client.Timeout exceeded while awaiting headers)", nErr)
}

func (hs *HTTPTestSuite) TestClientContext() {
	c := NewHTTPClient(&HTTPConfig{
		Timeout:     10000, // 10s
		DialTimeout: 10000, // 10s
	})

	nonRoutable := "10.255.255.255"
	url := "http://" + nonRoutable
	rctx, cancel := context.WithCancel(context.TODO())
	go func(rtx context.Context) {
		req, err := http.NewRequestWithContext(rtx, http.MethodGet, url, nil)
		hs.Nil(err)
		_, err = c.Do(req)
		hs.T().Logf("get %s err: %T, %+v", url, err, err)
		hs.Error(err)
		nErr, ok := err.(net.Error)
		hs.True(ok, err)
		hs.Contains(nErr.Error(), "context canceled", nErr)
	}(rctx)

	wait := time.Second * 1
	select {
	case <-rctx.Done():
		hs.T().Log("Request finished")
		cancel()
	case <-time.After(wait):
		hs.T().Log("Ready to cancel")
		cancel()
	}

	time.Sleep(time.Second)
}
