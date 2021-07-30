package scraper

import (
	"fmt"
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

func GetAvailableUrl(orginUrl string) (string, error) {

	resp := &AvailableResp{}
	err := client.GetJSON(fmt.Sprintf(archiveUrl, orginUrl), resp)
	if err != nil {
		return "", err
	}

	return resp.ArchivedSnapshots.Closest.URL, nil
}
