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
	mgstageDetailUrl = "https://www.mgstage.com/product/product_detail/%s/"
)

type MGStageScraper struct {
	doc        *goquery.Document
	docUrl     string
	HTTPClient *http.Client
}

func (s *MGStageScraper) FetchDoc(num string) error {
	if s.HTTPClient == nil {
		s.HTTPClient = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
			},
		}
	}
	s.docUrl = fmt.Sprintf(mgstageDetailUrl, strings.ToUpper(num))

	req, err := http.NewRequest("GET", s.docUrl, nil)
	if err != nil {
		return err
	}
	// $.cookie('adc','1',{domain:'mgstage.com',path:'/',expires:dt});
	c := &http.Cookie{
		Name:    "adc",
		Value:   "1",
		Path:    "/",
		Domain:  "mgstage.com",
		Expires: time.Now().Add(1 * time.Hour),
	}
	req.AddCookie(c)
	res, err := s.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return errors.New(fmt.Sprintf("status code error: %d %s", res.StatusCode, res.Status))
	}
	s.doc, err = goquery.NewDocumentFromReader(res.Body)
	return err
}

func (s *MGStageScraper) GetPlot() string {
	if s.doc == nil {
		return ""
	}
	return s.doc.Find("#introduction p").First().Text()
}

func (s *MGStageScraper) GetTitle() string {
	panic("implement me")
}

func (s *MGStageScraper) GetDirector() string {
	panic("implement me")
}

func (s *MGStageScraper) GetRuntime() string {
	panic("implement me")
}

func (s *MGStageScraper) GetTags() []string {
	panic("implement me")
}

func (s *MGStageScraper) GetStudio() string {
	panic("implement me")
}

func (s *MGStageScraper) GetMaker() string {
	panic("implement me")
}

func (s *MGStageScraper) GetOutline() string {
	panic("implement me")
}

func (s *MGStageScraper) GetActors() []string {
	panic("implement me")
}

func (s *MGStageScraper) GetLabel() string {
	panic("implement me")
}

func (s *MGStageScraper) GetNumber() string {
	panic("implement me")
}

func (s *MGStageScraper) GetCover() string {
	panic("implement me")
}

func (s *MGStageScraper) GetWebsite() string {
	panic("implement me")
}

func (s *MGStageScraper) GetPremiered() string {
	panic("implement me")
}

func (s *MGStageScraper) GetSeries() string {
	panic("implement me")
}
