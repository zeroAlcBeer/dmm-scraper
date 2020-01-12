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

func (s *MGStageScraper) SetHTTPClient(client *http.Client){
	s.HTTPClient = client
}

func (s *MGStageScraper) FetchDoc(num string) error {
	if s.HTTPClient == nil {
		s.HTTPClient = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
			},
		}
	}
	s.HTTPClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
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
	return strings.TrimSpace(s.doc.Find("#introduction p").First().Text())
}

func (s *MGStageScraper) GetTitle() string {
	if s.doc == nil {
		return ""
	}
	return strings.TrimSpace(s.doc.Find("h1[class=tag]").First().Text())
}

func (s *MGStageScraper) GetDirector() string {
	return ""
}

func (s *MGStageScraper) GetRuntime() string {
	if s.doc == nil {
		return ""
	}
	return strings.TrimSpace(getMgstageTableValue("収録時間", s.doc).Find("td").Text())
}

func (s *MGStageScraper) GetTags() (tags []string) {
	if s.doc == nil {
		return
	}
	getMgstageTableValue("ジャンル", s.doc).Find("td a").Each(
		func(i int, ss *goquery.Selection) {
			tags = append(tags, strings.TrimSpace(ss.Text()))
		})
	return
}

func (s *MGStageScraper) GetMaker() string {
	if s.doc == nil {
		return ""
	}
	return strings.TrimSpace(getMgstageTableValue("メーカー", s.doc).Find("td").Text())
}

func (s *MGStageScraper) GetActors() (actors []string) {
	if s.doc == nil {
		return
	}
	t := getMgstageTableValue("出演", s.doc)
	if t != nil {
		t.Find("td a").Each(
			func(i int, ss *goquery.Selection) {
				actors = append(actors, strings.TrimSpace(ss.Text()))
			})
	}

	return
}

func (s *MGStageScraper) GetLabel() string {
	if s.doc == nil {
		return ""
	}
	return strings.TrimSpace(getMgstageTableValue("レーベル", s.doc).Find("td").Text())
}

func (s *MGStageScraper) GetNumber() string {
	if s.doc == nil {
		return ""
	}
	return strings.TrimSpace(getMgstageTableValue("品番", s.doc).Find("td").Text())
}

func (s *MGStageScraper) GetCover() string {
	if s.doc == nil {
		return ""
	}
	u, _ := s.doc.Find("#EnlargeImage").Attr("href")
	return u
}

func (s *MGStageScraper) GetWebsite() string {
	return s.docUrl
}

func (s *MGStageScraper) GetPremiered() string {
	if s.doc == nil {
		return ""
	}
	return strings.TrimSpace(getMgstageTableValue("配信開始日", s.doc).Find("td").Text())
}

func (s *MGStageScraper) GetSeries() string {
	if s.doc == nil {
		return ""
	}
	return strings.TrimSpace(getMgstageTableValue("シリーズ", s.doc).Find("td").Text())
}

func getMgstageTableValue(key string, doc *goquery.Document) (target *goquery.Selection) {
	doc.Find("div[class=detail_data] table").Last().Find("tr").EachWithBreak(
		func(i int, s *goquery.Selection) bool {
			if strings.Contains(s.Text(), key) {
				target = s
				return false
			}
			return true
		})
	return
}
