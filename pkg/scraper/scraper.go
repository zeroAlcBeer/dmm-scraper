package scraper

import (
	"better-av-tool/pkg/archive"
	myclient "better-av-tool/pkg/client"
	"better-av-tool/pkg/config"
	"better-av-tool/pkg/logger"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
)

// Scraper is interface
type Scraper interface {
	// Remote
	FetchDoc(query string) (err error)

	// Local
	GetPlot() string
	GetTitle() string
	GetDirector() string
	GetRuntime() string
	GetTags() []string
	GetMaker() string
	GetActors() []string
	GetLabel() string
	GetNumber() string
	GetCover() string
	GetWebsite() string
	GetPremiered() string
	GetYear() string
	GetSeries() string

	// Operation
	NeedCut() bool
}

var (
	client myclient.Client
	log    logger.Logger
	cookie *http.Cookie
)

// Setup ...
func Setup(p config.Proxy) {
	log = logger.New()
	client = myclient.New()
	if p.Enable {
		err := client.SetProxyUrl(p.Socket)
		if err != nil {
			log.Errorf("Error parse proxy url, %s, proxy disabled", err)
		}
	}
	cookie = &http.Cookie{}
}

// GetDocFromURL ...
func GetDocFromURL(u string) (*goquery.Document, error) {
	log.Infof("fetching %s", u)
	res, err := client.Get(u, cookie)
	if err != nil {
		return nil, err
	}
	utfBody, _ := myclient.DecodeHTMLBody(res.Body, "")
	doc, err := goquery.NewDocumentFromReader(utfBody)
	if err != nil {
		return nil, err
	}
	doc.Url = res.Request.URL
	return doc, nil
}

// Download ...
func Download(url, filename string, progress func(current, total int64)) error {
	return client.Download(url, filename, progress)
}

func GetAvailableUrl(orginUrl string) (string, error) {

	resp := &archive.AvailableResp{}
	err := client.GetJSON(fmt.Sprintf(archive.GetAvailableUrl, orginUrl), resp)
	if err != nil {
		return "", err
	}

	return resp.ArchivedSnapshots.Closest.URL, nil
}
