package scraper

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	dmmSearchUrl = "https://www.dmm.co.jp/mono/dvd/-/search/=/searchstr=%s/"
)

type DMMScraper struct {
	doc        *goquery.Document
	docUrl     string
	HTTPClient *http.Client
}

func (d *DMMScraper) FetchDoc(num string) error {
	if d.HTTPClient == nil {
		d.HTTPClient = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
			},
		}
	}
	res, err := d.HTTPClient.Get(fmt.Sprintf(dmmSearchUrl, num))
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
	if itemCount > 2 {
		return errors.New("multi records, make number specific")
	}

	// 排除 特典\特卖
	firstUrl, _ := listDoc.Find("#list li").First().Find(".tmb a").Attr("href")
	lastUrl, _ := listDoc.Find("#list li").Last().Find(".tmb a").Attr("href")
	if len(firstUrl) < len(lastUrl) {
		d.docUrl = firstUrl
	} else {
		d.docUrl = lastUrl
	}

	resDetail, err := d.HTTPClient.Get(d.docUrl)
	if err != nil {
		return err
	}
	defer resDetail.Body.Close()
	if resDetail.StatusCode != 200 {
		return errors.New(fmt.Sprintf("status code error: %d %s", resDetail.StatusCode, resDetail.Status))
	}

	d.doc, err = goquery.NewDocumentFromReader(resDetail.Body)
	return err
}

func (d *DMMScraper) GetPlot() string {
	if d.doc == nil {
		return ""
	}
	return d.doc.Find("p[class=mg-b20]").Children().Remove().End().Text()
}

func (d *DMMScraper) GetTitle() string {
	if d.doc == nil {
		return ""
	}
	return d.doc.Find("#title").Text()
}

func (d *DMMScraper) GetDirector() string {
	if d.doc == nil {
		return ""
	}
	return getDmmTableValue("監督", d.doc)
}

func (d *DMMScraper) GetRuntime() string {
	if d.doc == nil {
		return ""
	}
	return getDmmTableValue("収録時間", d.doc)
}

func (d *DMMScraper) GetTags() (tags []string) {
	if d.doc == nil {
		return
	}
	d.doc.Find("table[class=mg-b20] tr").EachWithBreak(
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

func (d *DMMScraper) GetStudio() string {
	if d.doc == nil {
		return ""
	}
	return getDmmTableValue("メーカー", d.doc)
}

func (d *DMMScraper) GetMaker() string {
	if d.doc == nil {
		return ""
	}
	return getDmmTableValue("メーカー", d.doc)
}

func (d *DMMScraper) GetOutline() string {
	return ""
}

func (d *DMMScraper) GetActors() (actors []string) {
	if d.doc == nil {
		return
	}
	d.doc.Find("#performer").Each(func(i int, s *goquery.Selection) {
		actors = append(actors, s.Find("a").Text())
	})
	return
}

func (d *DMMScraper) GetLabel() string {
	if d.doc == nil {
		return ""
	}
	return getDmmTableValue("レーベル", d.doc)
}

func (d *DMMScraper) GetNumber() string {
	if d.doc == nil {
		return ""
	}
	return getDmmTableValue("品番", d.doc)
}

func (d *DMMScraper) GetCover() string {
	url, _ := d.doc.Find("#sample-video a").First().Attr("href")
	return url
}

func (d *DMMScraper) GetWebsite() string {
	return d.docUrl
}

func (d *DMMScraper) GetPremiered() string {
	if d.doc == nil {
		return ""
	}
	return getDmmTableValue("発売日", d.doc)
}

func (d *DMMScraper) GetSeries() string {
	if d.doc == nil {
		return ""
	}
	return getDmmTableValue("シリーズ", d.doc)
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
				return false
			}
			return true
		})
	return
}
