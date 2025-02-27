package client

import (
	"fmt"
	"net/http"

	"github.com/imroc/req/v3"
)

// Client interface defines the HTTP client operations
type Client interface {
	SetProxyUrl(rawurl string) error
	Get(url string, v interface{}) (*http.Response, error)
	GetJSON(url string, v interface{}) error
	Post(url string, v interface{}) (*http.Response, error)
	Download(url, filename string, progress func(current, total int64)) error
}

// New creates a new HTTP client instance
func New() Client {
	return &ReqClient{
		client: req.C().
			SetUserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36"),
	}
}

// DefaultProgress returns a default progress handler for downloads
func DefaultProgress() func(current, total int64) {
	return func(current, total int64) {
		fmt.Printf("%.2f%%\n", float32(current)/float32(total)*100)
	}
}
