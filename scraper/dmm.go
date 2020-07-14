package scraper

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

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

// 获取刮削的内容
func (s *DMMScraper) FetchDoc(query, url string) (err error) {
	// 如果有url，就直接从url刮
	if url != "" {

		s.doc, err = GetDocFromUrl(url)
		return err
	}

	// dmm 搜索页
	url = fmt.Sprintf(dmmSearchUrl, query)

	s.doc, err = GetDocFromUrl(url)
	if err != nil {
		return err
	}
	// 二次搜索
	if s.doc.Find("#list li").Length() == 0 {
		//log.Infof("fetching %s empty", s.doc.Text())
		url = fmt.Sprintf(dmmSearchUrl2, query)
		s.doc, err = GetDocFromUrl(url)
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
	nums := regexp.MustCompile("[0-9]+").FindAllString(query, -1)
	minLen := 100
	for _, href := range hrefs {
		if strings.Contains(href, nums[len(nums)-1]) && len(href) < minLen {
			url = href
			minLen = len(href)
		}
	}

	s.doc, err = GetDocFromUrl(url)
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
	return s.GetNumber() + " " + s.doc.Find("#title").Text()
}

func (s *DMMScraper) GetDirector() string {
	if s.doc == nil {
		return ""
	}
	return getDmmTableValue2(5, s.doc)
}

func (s *DMMScraper) GetRuntime() string {
	if s.doc == nil {
		return ""
	}
	return getDmmTableValue2(3, s.doc)
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
	return getDmmTableValue2(7, s.doc)
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
	return getDmmTableValue2(8, s.doc)
}

func (s *DMMScraper) GetNumber() string {
	if s.doc == nil {
		return ""
	}
	return getDmmTableValue2(10, s.doc)
}

func (s *DMMScraper) GetCover() string {
	if s.doc == nil {
		return ""
	}
	img, _ := s.doc.Find("#sample-video a").First().Attr("href")
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
	rel = getDmmTableValue2(2, s.doc)
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
	return getDmmTableValue2(6, s.doc)
}

func (s *DMMScraper) NeedCut() bool {
	return true
}

//func getDmmTableValue(key string, doc *goquery.Document) (val string) {
//	doc.Find("table[class=mg-b20] tr").EachWithBreak(
//		func(i int, s *goquery.Selection) bool {
//			if strings.Contains(s.Text(), key) {
//				val = s.Find("td a").Text()
//				if val == "" {
//					val = s.Find("td").Last().Text()
//				}
//				if val == "----" {
//					val = ""
//				}
//				val = strings.TrimSpace(val)
//				return false
//			}
//			return true
//		})
//	return
//}

func getDmmTableValue2(x int, doc *goquery.Document) (val string) {
	//log.Info(doc.Find("table[class=mg-b20] td[width]").Html())
	return doc.Find("table[class=mg-b20] td[width]").Eq(x-1).Text()
}
