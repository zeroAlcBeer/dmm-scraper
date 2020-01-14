package archive

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type AvailableResp struct {
	ArchivedSnapshots struct {
		Closest struct {
			Status    string `json:"status"`
			Available bool   `json:"available"`
			URL       string `json:"url"`
			Timestamp string `json:"timestamp"`
		} `json:"closest"`
	} `json:"archived_snapshots"`
	URL string `json:"url"`
}

const (
	archiveUrl = "https://archive.org/wayback/available?url=%s"
)

func GetAvailableUrl(orginUrl string, c *http.Client) (string, error) {
	res, err := c.Get(fmt.Sprintf(archiveUrl, orginUrl))
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return "", errors.New(fmt.Sprintf("status code error: %d %s", res.StatusCode, res.Status))
	}
	b, _ := ioutil.ReadAll(res.Body)
	resp := &AvailableResp{}
	err = json.Unmarshal(b, resp)
	if err != nil {
		return "", err
	}
	return resp.ArchivedSnapshots.Closest.URL, nil
}