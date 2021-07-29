package scraper

import (
	myclient "better-av-tool/pkg/client"
	"better-av-tool/pkg/logger"
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

func Init(c myclient.Client, l logger.Logger) {
	client = c
	log = l
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
