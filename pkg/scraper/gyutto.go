package scraper

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	gyuttoItemUrl = "http://gyutto.com/i/%s"
)

type GyuttoScraper struct {
	DefaultScraper
}

func (s *GyuttoScraper) GetType() string {
	return "GyuttoScraper"
}

func (s *GyuttoScraper) FetchDoc(query string) (err error) {
	u := fmt.Sprintf(gyuttoItemUrl, query)

	err = s.GetDocFromURL(u)
	if err != nil {
		return err
	}
	s.doc.Url, err = url.Parse(u)

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

func (s *GyuttoScraper) GetNumber() string {
	if s.doc == nil {
		return ""
	}
	nums := regexp.MustCompile("[0-9]+").FindAllString(s.doc.Url.String(), -1)
	return nums[len(nums)-1]
}

func (s *GyuttoScraper) GetFormatNumber() string {
	return strings.ToUpper(fmt.Sprintf("gyutto-%s", s.GetNumber()))
}

func (s *GyuttoScraper) GetCover() string {
	if s.doc == nil {
		return ""
	}
	img, _ := s.doc.Find(".highslide").First().Attr("href")
	return "http://image.gyutto.com" + img
}

func (s *GyuttoScraper) GetPremiered() (rel string) {
	if s.doc == nil {
		return ""
	}

	rel = strings.TrimSpace(s.doc.Find(".BasicInfo").Text())
	rel = regexp.MustCompile(`\d{4}年\d{2}月\d{2}日`).FindString(rel)
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
