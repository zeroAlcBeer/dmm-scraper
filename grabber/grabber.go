package grabber

import (
	"better-av-tool/log"
	"github.com/cavaliercoder/grab"
	"net/http"
)

var grabClient *grab.Client

func init() {
	grabClient = grab.NewClient()
}

func SetHTTPClient(client *http.Client) {
	grabClient.HTTPClient = client
}

func Download(u string) (string, error) {
	log.Infof("Downloading from %s", u)

	req, _ := grab.NewRequest("", u)
	resp := grabClient.Do(req)

	if err := resp.Err(); err != nil {
		return "", err
	}

	if resp.IsComplete() {
		log.Infof("Finished %s %d / %d bytes (%d%%)", resp.Filename, resp.BytesComplete(), resp.Size, int(100*resp.Progress()))
	}
	return resp.Filename, nil
}
