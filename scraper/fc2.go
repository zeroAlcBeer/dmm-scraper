package scraper

import (
	"better-av-tool/archive"
	"better-av-tool/log"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"regexp"
	"strings"
)

type Fc2Scraper struct {
	doc        *goquery.Document
	docUrl     string
	HTTPClient *http.Client
	isArchive  bool
}

const (
	fc2Url = "https://adult.contents.fc2.com/article_search.php?id=%s"
)

func (s *Fc2Scraper) FetchDoc(num string) error {
	if s.HTTPClient == nil {
		s.HTTPClient = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
			},
		}
	}

	typeFc2, _ := regexp.Compile(`[0-9]{6,7}`)
	num = typeFc2.FindString(num)

	if s.docUrl == "" {
		s.docUrl = fmt.Sprintf(fc2Url, num)
	}
	res, err := s.HTTPClient.Get(s.docUrl)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 && res.StatusCode != 410 {
		return errors.New(fmt.Sprintf("status code error: %d %s", res.StatusCode, res.Status))
	}

	if res.StatusCode == 410 {
		s.docUrl, err = archive.GetAvailableUrl(s.docUrl, s.HTTPClient)
		if err != nil {
			return err
		}
		res, err = s.HTTPClient.Get(s.docUrl)
		if err != nil {
			return err
		}
		if res.StatusCode != 200 {
			return errors.New(fmt.Sprintf("status code error: %d %s", res.StatusCode, res.Status))
		}
		s.isArchive = true
	}

	s.doc, err = goquery.NewDocumentFromReader(res.Body)
	return err
}

func (s *Fc2Scraper) SetHTTPClient(client *http.Client) {
	s.HTTPClient = client
}

func (s *Fc2Scraper) GetPlot() string {
	if s.doc == nil {
		return ""
	}
	tempDoc := s.doc.Find("section[class=explain] p").Clone()
	explain := tempDoc.Children().Remove().End().Text()
	return strings.TrimSpace(explain)
}

func (s *Fc2Scraper) GetTitle() string {
	if s.doc == nil {
		return ""
	}
	title := s.doc.Find("h2[class=title_bar]").Text()
	return strings.TrimSpace(title)
}

func (s *Fc2Scraper) GetDirector() string {
	if s.doc == nil {
		return ""
	}
	a := s.doc.Find("a[class=analyticsLinkClick_toUserPage1]").Text()
	return strings.TrimSpace(a)
}

func (s *Fc2Scraper) GetRuntime() (runtime string) {
	if s.doc == nil {
		return ""
	}
	s.doc.Find("div[class=main_info_block] dl dd").Each(func(i int, ss *goquery.Selection) {
		if i == 5 {
			runtime = strings.TrimSpace(ss.Text())
		}
	})
	return
}

func (s *Fc2Scraper) GetTags() (tags []string) {
	if s.doc == nil {
		return
	}
	s.doc.Find(".incident_tags a").Each(func(i int, ss *goquery.Selection) {
		tags = append(tags, strings.TrimSpace(ss.Text()))
	})
	return
}

func (s *Fc2Scraper) GetMaker() string {
	if s.doc == nil {
		return ""
	}
	a := s.doc.Find("a[class=analyticsLinkClick_toUserPage1]").Text()
	return strings.TrimSpace(a)
}

func (s *Fc2Scraper) GetActors() []string {
	return []string{}
}

func (s *Fc2Scraper) GetLabel() string {
	return ""
}

func (s *Fc2Scraper) GetNumber() string {
	if s.doc == nil {
		return ""
	}
	id, _ := s.doc.Find("#reviews").Attr("data-id")
	return strings.TrimSpace(id)
}

func (s *Fc2Scraper) GetCover() string {
	if s.doc == nil {
		return ""
	}
	img, _ := s.doc.Find("section[class=explain] img").First().Attr("src")
	if s.isArchive {
		img = img[strings.LastIndex(img, "http"):]
		log.Info(img)
	}
	return img
}

func (s *Fc2Scraper) GetWebsite() string {
	return s.docUrl
}

func (s *Fc2Scraper) GetPremiered() (rel string) {
	if s.doc == nil {
		return
	}
	s.doc.Find("div[class=main_info_block] dl dd").Each(func(i int, ss *goquery.Selection) {
		if i == 3 {
			rel = strings.TrimSpace(ss.Text())
		}
	})
	return
}

func (s *Fc2Scraper) GetSeries() string {
	return ""
}
