package scraper

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	gyuttoSearchUrl = "http://gyutto.com/search/search_list.php?category_id=10&mode=search&sub_category_id=21&search_keyword=%s"
)

type GyuttoScraper struct {
	doc *goquery.Document
}

func (s *GyuttoScraper) FetchDoc(query string) (err error) {
	u := fmt.Sprintf(gyuttoSearchUrl, url.QueryEscape(query))
	s.doc, err = GetDocFromURL(u)
	if err != nil {
		return err
	}

	var hrefs []string
	s.doc.Find(".ListBox li").Each(func(i int, se *goquery.Selection) {
		href, _ := se.Find(".DefiPhotoName a").Attr("href")
		hrefs = append(hrefs, href)
	})
	if len(hrefs) == 0 {
		return errors.New("record not found")
	}

	s.doc, err = GetDocFromURL(hrefs[0])
	if err != nil {
		return err
	}
	s.doc.Url, err = url.Parse(hrefs[0])

	//utfDetailBody, _ := DecodeHTMLBody(res.Body, "Shift_JIS")
	return err
}

func (s *GyuttoScraper) GetTitle() string {
	if s.doc == nil {
		return ""
	}
	return strings.TrimSpace(s.doc.Find(".parts_Mds01 h1").Text())
}

func (s *GyuttoScraper) GetPlot() string {
	if s.doc == nil {
		return ""
	}
	return strings.TrimSpace(s.doc.Find(".unit_DetailSummary p").Text())
}

func (s *GyuttoScraper) GetDirector() string {
	if s.doc == nil {
		return ""
	}
	return strings.TrimSpace(s.doc.Find(".BasicInfo").Eq(2).Find("dd a").First().Text())
}

func (s *GyuttoScraper) GetRuntime() string {
	return ""
}

func (s *GyuttoScraper) GetTags() (tags []string) {
	if s.doc == nil {
		return
	}
	s.doc.Find(".BasicInfo").Eq(3).Find("dd a").Each(
		func(i int, s *goquery.Selection) {
			tags = append(tags, s.Text())
		})
	return
}

func (s *GyuttoScraper) GetMaker() string {
	return s.GetDirector()
}

func (s *GyuttoScraper) GetActors() (actors []string) {
	return
}

func (s *GyuttoScraper) GetLabel() string {
	return ""
}

func (s *GyuttoScraper) GetNumber() string {
	if s.doc == nil {
		return ""
	}
	nums := regexp.MustCompile("[0-9]+").FindAllString(s.doc.Url.String(), -1)
	return nums[len(nums)-1]
}

func (s *GyuttoScraper) GetCover() string {
	if s.doc == nil {
		return ""
	}
	img, _ := s.doc.Find(".highslide").First().Attr("href")
	return "http://image.gyutto.com" + img
}

func (s *GyuttoScraper) GetWebsite() string {
	if s.doc == nil {
		return ""
	}
	return s.doc.Url.String()
}

func (s *GyuttoScraper) GetPremiered() (rel string) {
	if s.doc == nil {
		return ""
	}

	rel = strings.TrimSpace(s.doc.Find(".BasicInfo").Eq(4).Find("dd").Text())
	rel = strings.Replace(rel, "年", "-", 1)
	rel = strings.Replace(rel, "月", "-", 1)
	rel = strings.Replace(rel, "日", "", 1)
	return
}

func (s *GyuttoScraper) GetYear() (rel string) {
	if s.doc == nil {
		return ""
	}
	return regexp.MustCompile(`\d{4}`).FindString(s.GetPremiered())
}

func (s *GyuttoScraper) GetSeries() string {
	return ""
}

func (s *GyuttoScraper) NeedCut() bool {
	return false
}
