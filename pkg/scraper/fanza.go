package scraper

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	dmmDigitalSearchUrl = "https://www.dmm.co.jp/digital/-/list/search/=/?searchstr=%s"
)

type FanzaScraper struct {
	DMMScraper
}

func (s *FanzaScraper) GetType() string {
	return "FanzaScraper"
}

// FetchDoc search once or twice to get a detail page
func (s *FanzaScraper) FetchDoc(query string) (err error) {
	s.cookie = &http.Cookie{
		Name:    "age_check_done",
		Value:   "1",
		Path:    "/",
		Domain:  "dmm.co.jp",
		Expires: time.Now().Add(1 * time.Hour),
	}

	// dmm 搜索页
	if strings.Contains(query, "-") {
		query = strings.Replace(query, "-", "00", 1)
	}
	err = s.GetDocFromURL(fmt.Sprintf(dmmDigitalSearchUrl, query))
	if err != nil {
		return err
	}

	var hrefs []string
	s.doc.Find("#list li").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Find(".tmb a").Attr("href")
		hrefs = append(hrefs, href)
	})

	if len(hrefs) == 0 {
		return errors.New("record not found")
	}
	// 多个结果时，取最短长度
	var detail string
	for _, href := range hrefs {
		if isURLMatchQuery(href, query) {
			detail = href
		}
	}
	if detail == "" {
		return fmt.Errorf("unable to match in hrefs %v", hrefs)
	}

	return s.GetDocFromURL(detail)
}

func (s *FanzaScraper) NeedCut() bool {
	return needCut
}
