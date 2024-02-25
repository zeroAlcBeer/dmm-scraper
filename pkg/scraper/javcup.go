package scraper

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

const (
	javcupDetailUrl = "https://javcup.com/movie/%s"
)

type JavCupScraper struct {
	DefaultScraper
	formatQuery string
}

func (s *JavCupScraper) GetType() string {
	return "JavCupScraper"
}

func (s *JavCupScraper) FetchDoc(query string) (err error) {
	//s.cookie = &http.Cookie{
	//	Name:    "adc",
	//	Value:   "1",
	//	Path:    "/",
	//	Domain:  "mgstage.com",
	//	Expires: time.Now().Add(1 * time.Hour),
	//}
	l, i := GetLabelNumber(query)
	if l == "" {
		return fmt.Errorf("unable to GetLabelNumber")
	}
	s.formatQuery = strings.ToUpper(fmt.Sprintf("%s-%03d", l, i))
	u := fmt.Sprintf(javcupDetailUrl, s.formatQuery)

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	timeoutCtx, timeoutCancel := context.WithTimeout(ctx, 10*time.Second)
	defer timeoutCancel()

	// 定义要获取的内容变量
	var response string

	// 执行任务：打开页面，等待某个特定元素加载完成，获取内容
	err = chromedp.Run(timeoutCtx,
		chromedp.Navigate(u),
		// 根据实际情况替换下面的选择器
		chromedp.WaitVisible(`#video`, chromedp.ByID),
		chromedp.OuterHTML(`html`, &response),
	)
	if err != nil {
		return err
	}

	s.doc, err = goquery.NewDocumentFromReader(strings.NewReader(response)) //

	if err != nil {
		return err
	}
	parsedURL, _ := url.Parse(u)
	s.doc.Url = parsedURL

	return nil
}

func (s *JavCupScraper) GetPlot() string {
	if s.doc == nil {
		return ""
	}
	return strings.TrimSpace(s.doc.Find("div[class=movie-description] p").First().Text())
}

func (s *JavCupScraper) GetTitle() string {
	if s.doc == nil {
		return ""
	}
	title := strings.TrimSpace(s.doc.Find("h1[itemprop=name]").First().Text())
	return strings.TrimSpace(strings.Replace(title, s.formatQuery, "", 1))
}

func (s *JavCupScraper) GetDirector() string {
	if s.doc == nil {
		return ""
	}
	return strings.TrimSpace(s.doc.Find("span[itemprop=director]").First().Text())
}

func (s *JavCupScraper) GetRuntime() string {
	return ""
}

func (s *JavCupScraper) GetTags() (tags []string) {
	if s.doc == nil {
		return
	}

	s.doc.Find("p[itemprop=keywords]").Find("span").Each(
		func(i int, ss *goquery.Selection) {
			tags = append(tags, strings.TrimSpace(ss.Text()))
		})
	return
}

func (s *JavCupScraper) GetMaker() string {
	if s.doc == nil {
		return ""
	}
	return strings.TrimSpace(s.doc.Find("span[itemprop=genre]").First().Text())
}

func (s *JavCupScraper) GetActors() (actors []string) {
	if s.doc == nil {
		return
	}
	s.doc.Find("div[class=model-item]").Find("span").Each(
		func(i int, ss *goquery.Selection) {
			actors = append(actors, strings.TrimSpace(ss.Text()))
		})

	return
}

func (s *JavCupScraper) GetLabel() string {
	return ""
}

func (s *JavCupScraper) GetNumber() string {
	return s.formatQuery
}

func (s *JavCupScraper) GetCover() string {
	if s.doc == nil {
		return ""
	}
	u, _ := s.doc.Find("#video").Attr("poster")
	u = strings.Replace(u, "cdn.javcup.com", "pics.dmm.co.jp", 1)
	u = strings.Replace(u, "img.javcup.com", "pics.dmm.co.jp", 1)
	return u
}

func (s *JavCupScraper) GetPremiered() (rel string) {
	if s.doc == nil {
		return ""
	}
	return strings.TrimSpace(s.doc.Find("span[itemprop=datePublished]").First().Text())
}

func (s *JavCupScraper) GetYear() (rel string) {
	if s.doc == nil {
		return ""
	}
	return regexp.MustCompile(`\d{4}`).FindString(s.GetPremiered())
}

func (s *JavCupScraper) GetSeries() string {
	return ""
}

func (s *JavCupScraper) GetFormatNumber() string {
	return s.GetNumber()
}
