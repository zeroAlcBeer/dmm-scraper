package scraper

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Fc2Scraper struct {
	DefaultScraper
	fc2data *FC2Data
}

func (s *Fc2Scraper) GetType() string {
	return "Fc2Scraper"
}

const (
	fc2Url  = "https://adult.contents.fc2.com/article/%s/"
	fc2Url2 = "https://adult.contents.fc2.com/article_search.php?id=%s"
)

type FC2Data struct {
	//Type        string `json:"@type"`
	//ID          string `json:"@id"`
	//Sku         string `json:"sku"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Context     string   `json:"@context"`
	Image       FC2Image `json:"image"`
	ProductID   string   `json:"productID"`
	//Brand           Brand           `json:"brand"`
	//AggregateRating AggregateRating `json:"aggregateRating"`
	//PotentialAction PotentialAction `json:"potentialAction"`
	//Offers          Offers          `json:"offers"`
}

type FC2Image struct {
	URL  string `json:"url"`
	Type string `json:"@type"`
}

func (s *Fc2Scraper) FetchDoc(query string) (err error) {

	err = s.GetDocFromURL(fmt.Sprintf(fc2Url, query))
	if err != nil {
		return err
	}
	if len(s.doc.Find(".items_notfound_header").Nodes) != 0 {
		var u string
		u, err = s.GetAvailableUrl(fmt.Sprintf(fc2Url2, query))
		if err != nil {
			return err
		}
		s.isArchive = true
		err = s.GetDocFromURL(u)
		if err != nil {
			return err
		}
	}
	if !s.isArchive {
		err = getFc2Data(s)
	}
	return err
}

func (s *Fc2Scraper) GetTitle() string {
	if s.doc == nil {
		return ""
	}

	if s.isArchive {
		title := s.doc.Find("h2[class=title_bar]").Text()
		return strings.TrimSpace(title)
	}
	//title := s.doc.Find("div[class=items_article_headerInfo] h3").Text()
	return s.fc2data.Name
}

func (s *Fc2Scraper) GetPlot() string {
	if s.doc == nil {
		return ""
	}
	if s.isArchive {
		tempDoc := s.doc.Find("section[class=explain] p").Clone()
		explain := tempDoc.Children().Remove().End().Text()
		return strings.TrimSpace(explain)
	}
	//tempDoc := s.doc.Find("section[class=items_article_Contents] div").Clone()
	//explain := tempDoc.Children().Remove().End().Text()
	return s.fc2data.Description
}

func (s *Fc2Scraper) GetDirector() string {
	if s.doc == nil {
		return ""
	}
	if s.isArchive {
		a := s.doc.Find("a[class=analyticsLinkClick_toUserPage1]").Text()
		return strings.TrimSpace(a)
	}
	a := s.doc.Find("div[class=items_article_headerInfo] ul li").Eq(2).Find("a").Text()
	return strings.TrimSpace(a)
}

func (s *Fc2Scraper) GetRuntime() (runtime string) {
	return ""
}

func (s *Fc2Scraper) GetTags() (tags []string) {
	if s.doc == nil {
		return
	}
	if s.isArchive {
		s.doc.Find(".incident_tags a").Each(func(i int, ss *goquery.Selection) {
			tags = append(tags, strings.TrimSpace(ss.Text()))
		})
		return
	}
	s.doc.Find("a[class=tagTag]").Each(func(i int, ss *goquery.Selection) {
		tags = append(tags, strings.TrimSpace(ss.Text()))
	})
	return
}

func (s *Fc2Scraper) GetMaker() string {
	if s.doc == nil {
		return ""
	}
	return s.GetDirector()
}

func (s *Fc2Scraper) GetActors() []string {
	return []string{}
}

func (s *Fc2Scraper) GetNumber() string {
	if s.doc == nil {
		return ""
	}
	if s.isArchive {
		id, _ := s.doc.Find("#reviews").Attr("data-id")
		return strings.TrimSpace(id)
	}
	//id, _ := s.doc.Find(".items_article_TagArea").Attr("data-id")
	return s.fc2data.ProductID
}

func (s *Fc2Scraper) GetCover() string {
	if s.doc == nil {
		return ""
	}
	if s.isArchive {
		img, _ := s.doc.Find("section[class=explain] img").First().Attr("src")
		return img[strings.LastIndex(img, "http"):]
	}

	return s.fc2data.Image.URL
}

func (s *Fc2Scraper) GetPremiered() (rel string) {
	if s.doc == nil {
		return
	}
	if s.isArchive {
		rel = s.doc.Find("div[class=main_info_block] dl dd").Text()
	} else {
		rel = s.doc.Find(".items_article_Releasedate p").Text()
	}
	rel = regexp.MustCompile(`\d{4}\/(0?[1-9]|1[012])\/(0?[1-9]|[12][0-9]|3[01])`).FindString(rel)
	return strings.Replace(rel, "/", "-", -1)
}

func (s *Fc2Scraper) GetYear() (rel string) {
	if s.doc == nil {
		return ""
	}
	return regexp.MustCompile(`\d{4}`).FindString(s.GetPremiered())
}

func (s *Fc2Scraper) GetFormatNumber() string {
	return strings.ToUpper(fmt.Sprintf("fc2-%s", s.GetNumber()))
}

func getFc2Data(s *Fc2Scraper) error {
	data := s.doc.Find("script[type='application/ld+json']").Text()
	s.fc2data = &FC2Data{}
	err := json.Unmarshal([]byte(data), s.fc2data)
	return err
}
