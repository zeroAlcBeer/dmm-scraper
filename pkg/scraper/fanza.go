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
	dmmDigitalSearchUrl = "https://www.dmm.co.jp/digital/-/list/search/=/?searchstr=%s"
)

type FanzaScraper struct {
	DefaultScraper
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

func (s *FanzaScraper) GetPlot() string {
	if s.doc == nil {
		return ""
	}
	tempDoc := s.doc.Find("p[class=mg-b20]")
	return strings.TrimSpace(tempDoc.Children().Remove().End().Text())
}

func (s *FanzaScraper) GetTitle() string {
	if s.doc == nil {
		return ""
	}
	return s.doc.Find("#title").Text()
}

func (s *FanzaScraper) GetDirector() string {
	if s.doc == nil {
		return ""
	}
	return getDmmTableValue("監督", s.doc)
}

func (s *FanzaScraper) GetRuntime() string {
	if s.doc == nil {
		return ""
	}
	return getDmmTableValue("収録時間", s.doc)
}

func (s *FanzaScraper) GetTags() (tags []string) {
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

func (s *FanzaScraper) GetMaker() string {
	if s.doc == nil {
		return ""
	}
	return getDmmTableValue("メーカー", s.doc)
}

func (s *FanzaScraper) GetActors() (actors []string) {
	if s.doc == nil {
		return
	}
	s.doc.Find("#performer a").Each(func(i int, s *goquery.Selection) {
		actors = append(actors, s.Text())
	})
	return
}

func (s *FanzaScraper) GetLabel() string {
	if s.doc == nil {
		return ""
	}
	return getDmmTableValue("レーベル", s.doc)
}

func (s *FanzaScraper) GetNumber() string {
	if s.doc == nil {
		return ""
	}
	return getDmmTableValue("品番", s.doc)
}

func (s *FanzaScraper) GetFormatNumber() string {
	l, i := GetLabelNumber(s.GetNumber())
	if l == "" {
		return fmt.Sprintf("%03d", i)
	}
	return strings.ToUpper(fmt.Sprintf("%s-%03d", l, i))
}

func (s *FanzaScraper) GetCover() string {
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

func (s *FanzaScraper) GetPremiered() (rel string) {
	if s.doc == nil {
		return ""
	}
	rel = getDmmTableValue("発売日", s.doc)
	if rel == "" {
		rel = getDmmTableValue("配信開始日", s.doc)
	}
	return strings.Replace(rel, "/", "-", -1)
}

func (s *FanzaScraper) GetYear() (rel string) {
	if s.doc == nil {
		return ""
	}
	return regexp.MustCompile(`\d{4}`).FindString(s.GetPremiered())
}

func (s *FanzaScraper) GetSeries() string {
	if s.doc == nil {
		return ""
	}
	return getDmmTableValue("シリーズ", s.doc)
}

func (s *FanzaScraper) NeedCut() bool {
	return true
}
