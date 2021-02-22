package client

import (
	"fmt"
	"net/http"

	"github.com/imroc/req"
)

// Client ...
type Client interface {
	SetProxyUrl(rawurl string) error
	Get(url string, v ...interface{}) (*http.Response, error)
	GetJSON(url string, v interface{}) error
	Post(url string, v ...interface{}) (*http.Response, error)
	Download(url, filename string, progress func(current, total int64)) error
}

// New ...
func New() Client {
	return &ReqClient{
		req.New(),
	}
}

// DefaultProgress ...
func DefaultProgress() func(current, total int64) {
	return func(current, total int64) {
		fmt.Println(float32(current)/float32(total)*100, "%")
	}
}
