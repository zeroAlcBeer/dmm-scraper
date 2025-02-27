package scraper

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	avwikiSearchUrl = "https://av-wiki.net/?s=%s"
)

type AVWikiScraper struct {
	DefaultScraper
	formatQuery string
}

func (s *AVWikiScraper) GetType() string {
	return "AVWikiScraper"
}

func (s *AVWikiScraper) FetchDoc(query string) (err error) {
	l, i := GetLabelNumber(query)
	if l == "" {
		return fmt.Errorf("unable to GetLabelNumber")
	}
	s.formatQuery = strings.ToUpper(fmt.Sprintf("%s-%03d", l, i))
	u := fmt.Sprintf(avwikiSearchUrl, s.formatQuery)

	err = s.GetDocFromURL(u)
	if err != nil {
		return err
	}

	var hrefs []string
	s.doc.Find(".archive-list").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Find(".read-more a").Attr("href")
		hrefs = append(hrefs, href)
	})

	if len(hrefs) == 0 {
		return errors.New("record not found")
	}

	return s.GetDocFromURL(hrefs[0])
}

func (s *AVWikiScraper) GetPlot() string {
	if s.doc == nil {
		return ""
	}
	title, _ := s.doc.Find(".article-thumbnail a img").Attr("alt")
	return strings.TrimSpace(title)
}

func (s *AVWikiScraper) GetTitle() string {
	if s.doc == nil {
		return ""
	}
	title, _ := s.doc.Find(".article-thumbnail a img").Attr("alt")
	return strings.TrimSpace(title)
}

func (s *AVWikiScraper) GetDirector() string {
	if s.doc == nil {
		return ""
	}
	return strings.TrimSpace(s.doc.Find("span[itemprop=director]").First().Text())
}

func (s *AVWikiScraper) GetRuntime() string {
	return ""
}

func (s *AVWikiScraper) GetTags() (tags []string) {
	if s.doc == nil {
		return
	}

	tags = make([]string, 0)

	s.doc.Find("div.cat-link a").Each(func(i int, selection *goquery.Selection) {
		tag := strings.TrimSpace(selection.Text())
		if tag != "" {
			tags = append(tags, tag)
		}
	})

	return tags
}

func (s *AVWikiScraper) GetMaker() string {
	if s.doc == nil {
		return ""
	}
	return strings.TrimSpace(s.doc.Find("dl.dltable dt:contains('メーカー')").First().Next().Text())
}

func (s *AVWikiScraper) GetActors() (actors []string) {
	if s.doc == nil {
		return nil
	}
	var actresses []string
	s.doc.Find("dl.dltable dt:contains('AV女優名')").First().Next().Find("a").Each(func(i int, s *goquery.Selection) {
		actresses = append(actresses, strings.TrimSpace(s.Text()))
	})
	return actresses
}

func (s *AVWikiScraper) GetLabel() string {
	if s.doc == nil {
		return ""
	}
	return strings.TrimSpace(s.doc.Find("dl.dltable dt:contains('レーベル')").First().Next().Text())
}

func (s *AVWikiScraper) GetNumber() string {
	return s.formatQuery
}

func (s *AVWikiScraper) GetCover() string {
	if s.doc == nil {
		return ""
	}
	u, _ := s.doc.Find(".article-thumbnail a img").Attr("src")
	return u
}

func (s *AVWikiScraper) GetPremiered() (rel string) {
	if s.doc == nil {
		return ""
	}
	return strings.TrimSpace(s.doc.Find("dl.dltable dt:contains('配信開始日')").First().Next().Text())
}

func (s *AVWikiScraper) GetYear() (rel string) {
	if s.doc == nil {
		return ""
	}
	return regexp.MustCompile(`\d{4}`).FindString(s.GetPremiered())
}

func (s *AVWikiScraper) GetSeries() string {
	return ""
}

func (s *AVWikiScraper) GetFormatNumber() string {
	return s.GetNumber()
}
