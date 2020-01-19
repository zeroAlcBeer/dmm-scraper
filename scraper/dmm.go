package scraper

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"better-av-tool/log"
)

const (
	dmmSearchUrl = "https://www.dmm.co.jp/mono/dvd/-/search/=/searchstr=%s/"
)

type DMMScraper struct {
	doc        *goquery.Document
	docUrl     string
	HTTPClient *http.Client
}

func (s *DMMScraper) SetHTTPClient(client *http.Client) {
	s.HTTPClient = client
}

func (s *DMMScraper) FetchDoc(num string) error {
	if s.HTTPClient == nil {
		s.HTTPClient = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
			},
		}
	}
	log.Infof("fetching %s", fmt.Sprintf(dmmSearchUrl, num))
	res, err := s.HTTPClient.Get(fmt.Sprintf(dmmSearchUrl, num))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return errors.New(fmt.Sprintf("status code error: %d %s", res.StatusCode, res.Status))
	}

	listDoc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		return err
	}

	// Find the ul id=list items
	itemCount := listDoc.Find("#list li").Length()
	if itemCount == 0 {
		return errors.New("record not found")
	}
	var hrefs []string
	listDoc.Find("#list li").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Find(".tmb a").Attr("href")
		if strings.Contains(href, strings.Replace(strings.ToLower(num), "-", "", 1)) {
			hrefs = append(hrefs, href)
		}
	})

	if len(hrefs) == 0 {
		return errors.New("fail to make number specific")
	}

	s.docUrl = hrefs[0]
	if len(hrefs) > 1 {
		for _, href := range hrefs[1:] {
			if len(href) < len(s.docUrl) {
				s.docUrl = href
			}
		}
	}
	log.Infof("fetching %s", s.docUrl)
	resDetail, err := s.HTTPClient.Get(s.docUrl)
	if err != nil {
		return err
	}
	defer resDetail.Body.Close()
	if resDetail.StatusCode != 200 {
		return errors.New(fmt.Sprintf("status code error: %d %s", resDetail.StatusCode, resDetail.Status))
	}

	s.doc, err = goquery.NewDocumentFromReader(resDetail.Body)
	return err
}

func (s *DMMScraper) GetPlot() string {
	if s.doc == nil {
		return ""
	}
	tempDoc := s.doc.Find("p[class=mg-b20]")
	return strings.TrimSpace(tempDoc.Children().Remove().End().Text())
}

func (s *DMMScraper) GetTitle() string {
	if s.doc == nil {
		return ""
	}
	return s.doc.Find("#title").Text()
}

func (s *DMMScraper) GetDirector() string {
	if s.doc == nil {
		return ""
	}
	return getDmmTableValue("監督", s.doc)
}

func (s *DMMScraper) GetRuntime() string {
	if s.doc == nil {
		return ""
	}
	return getDmmTableValue("収録時間", s.doc)
}

func (s *DMMScraper) GetTags() (tags []string) {
	if s.doc == nil {
		return
	}
	s.doc.Find("table[class=mg-b20] tr").EachWithBreak(
		func(i int, s *goquery.Selection) bool {
			if strings.Contains(s.Text(), "ジャンル") {
				s.Find("td a").Each(func(i int, ss *goquery.Selection) {
					tags = append(tags, ss.Text())
				})
				return false
			}
			return true
		})
	return
}

func (s *DMMScraper) GetMaker() string {
	if s.doc == nil {
		return ""
	}
	return getDmmTableValue("メーカー", s.doc)
}

func (s *DMMScraper) GetActors() (actors []string) {
	if s.doc == nil {
		return
	}
	s.doc.Find("#performer").Each(func(i int, s *goquery.Selection) {
		actors = append(actors, s.Find("a").Text())
	})
	return
}

func (s *DMMScraper) GetLabel() string {
	if s.doc == nil {
		return ""
	}
	return getDmmTableValue("レーベル", s.doc)
}

func (s *DMMScraper) GetNumber() string {
	if s.doc == nil {
		return ""
	}
	return getDmmTableValue("品番", s.doc)
}

func (s *DMMScraper) GetCover() string {
	url, _ := s.doc.Find("#sample-video a").First().Attr("href")
	return url
}

func (s *DMMScraper) GetWebsite() string {
	return s.docUrl
}

func (s *DMMScraper) GetPremiered() string {
	if s.doc == nil {
		return ""
	}
	return getDmmTableValue("発売日", s.doc)
}

func (s *DMMScraper) GetSeries() string {
	if s.doc == nil {
		return ""
	}
	return getDmmTableValue("シリーズ", s.doc)
}

func (s *DMMScraper) NeedCut() bool {
	return true
}
//
func getDmmTableValue(key string, doc *goquery.Document) (val string) {
	doc.Find("table[class=mg-b20] tr").EachWithBreak(
		func(i int, s *goquery.Selection) bool {
			if strings.Contains(s.Text(), key) {
				val = s.Find("td a").Text()
				if val == "" {
					val = s.Find("td").Last().Text()
				}
				if val == "----" {
					val = ""
				}
				return false
			}
			return true
		})
	return
}
