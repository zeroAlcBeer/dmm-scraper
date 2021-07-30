package scraper

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	dmmSearchUrl  = "https://www.dmm.co.jp/mono/dvd/-/search/=/searchstr=%s/"
	dmmSearchUrl2 = "https://www.dmm.co.jp/search/=/searchstr=%s/n1=FgRCTw9VBA4GAVhfWkIHWw__/"
)

type DMMScraper struct {
	doc *goquery.Document
	//isArchive bool
}

func (s *DMMScraper) Cookie() *http.Cookie {
	return &http.Cookie{
		Name:    "age_check_done",
		Value:   "1",
		Path:    "/",
		Domain:  "dmm.co.jp",
		Expires: time.Now().Add(1 * time.Hour),
	}
}

// FetchDoc search once or twice to get a detail page
func (s *DMMScraper) FetchDoc(query string) (err error) {
	cookie = s.Cookie()

	// dmm 搜索页
	if strings.Contains(query, "-") {
		query = strings.Replace(query, "-", "00", 1)
	}
	s.doc, err = GetDocFromURL(fmt.Sprintf(dmmSearchUrl, query))
	if err != nil {
		return err
	}
	// 二次搜索
	if s.doc.Find("#list li").Length() == 0 {
		s.doc, err = GetDocFromURL(fmt.Sprintf(dmmSearchUrl2, query))
		if err != nil {
			return err
		}
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

	s.doc, err = GetDocFromURL(detail)
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
	s.doc.Find("#performer a").Each(func(i int, s *goquery.Selection) {
		actors = append(actors, s.Text())
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
	if s.doc == nil {
		return ""
	}
	img, _ := s.doc.Find("#sample-video img").First().Attr("src")
	//if strings.Contains(img, "web.archive.org") {
	//	img = fmt.Sprintf("http://pics.dmm.co.jp/digital/video/%s/%spl.jpg", s.GetNumber(), s.GetNumber())
	//}
	if strings.Contains(img, "ps.jpg") {
		img = strings.Replace(img, "ps.jpg", "pl.jpg", 1)
	}
	//
	//if s.isArchive {
	//	img = img[strings.LastIndex(img, "http"):]
	//	log.Info(img)
	//}
	return img
}

func (s *DMMScraper) GetWebsite() string {
	if s.doc == nil {
		return ""
	}
	return s.doc.Url.String()
}

func (s *DMMScraper) GetPremiered() (rel string) {
	if s.doc == nil {
		return ""
	}
	rel = getDmmTableValue("発売日", s.doc)
	if rel == "" {
		rel = getDmmTableValue("配信開始日", s.doc)
	}
	return strings.Replace(rel, "/", "-", -1)
}

func (s *DMMScraper) GetYear() (rel string) {
	if s.doc == nil {
		return ""
	}
	return regexp.MustCompile(`\d{4}`).FindString(s.GetPremiered())
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
				val = strings.TrimSpace(val)
				return false
			}
			return true
		})
	return
}

func getDmmTableValue2(x int, doc *goquery.Document) (val string) {
	//log.Info(doc.Find("table[class=mg-b20] td[width]").Html())
	return doc.Find("table[class=mg-b20] td").Eq(x - 1).Text()
}
