package scraper

import (
	myclient "better-av-tool/internal/client"
	"better-av-tool/internal/logger"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

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
