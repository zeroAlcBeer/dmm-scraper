package archive

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
	GetAvailableUrl = "https://archive.org/wayback/available?url=%s"
)


