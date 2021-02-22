package client

import (
	"net/http"

	"github.com/imroc/req"
)

// ReqClient ...
type ReqClient struct {
	*req.Req
}

// Get ...
func (rc *ReqClient) Get(url string, v ...interface{}) (*http.Response, error) {
	r, err := rc.Req.Get(url, v...)
	return r.Response(), err
}

// GetJSON ...
func (rc *ReqClient) GetJSON(url string, v interface{}) error {
	r, err := rc.Req.Get(url)
	if err != nil {
		return err
	}

	return r.ToJSON(v)
}

// Post ...
func (rc *ReqClient) Post(url string, v ...interface{}) (*http.Response, error) {
	r, err := rc.Req.Post(url, v...)
	return r.Response(), err
}

// Download ...
func (rc *ReqClient) Download(url, filename string, progress func(current, total int64)) error {
	r, err := rc.Req.Get(url, req.DownloadProgress(progress))
	if err != nil {
		return err
	}
	return r.ToFile(filename)
}
