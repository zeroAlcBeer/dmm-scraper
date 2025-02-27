package scraper

import (
	"encoding/json"
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
	DMMScraper
	schemaData *SchemaData
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
		strs := strings.Split(query, "-")
		if len(strs) == 2 {
			query = strs[0] + fmt.Sprintf("%05s", strs[1])
		}
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

	err = s.GetDocFromURL(detail)
	if err != nil {
		return err
	}
	return s.parseSchema()
}

func (s *FanzaScraper) GetPlot() string {
	if s.schemaData == nil {
		return ""
	}
	return s.schemaData.Description
}

func (s *FanzaScraper) GetTitle() string {
	if s.schemaData == nil {
		return ""
	}
	return s.schemaData.Name
}

func (s *FanzaScraper) GetTags() (tags []string) {
	if s.schemaData == nil {
		return nil
	}
	return s.schemaData.SubjectOf.Genre
}

func (s *FanzaScraper) GetMaker() string {
	if s.schemaData == nil {
		return ""
	}
	return s.schemaData.Brand.Name
}

func (s *FanzaScraper) GetActors() (actors []string) {
	if s.schemaData == nil {
		return nil
	}
	return []string{s.schemaData.SubjectOf.Actor.Name}
}

func (s *FanzaScraper) GetLabel() string {
	if s.schemaData == nil {
		return ""
	}
	return s.schemaData.Brand.Name
}

func (s *FanzaScraper) GetNumber() string {
	if s.schemaData == nil {
		return ""
	}
	return s.schemaData.SKU
}

func (s *FanzaScraper) GetFormatNumber() string {
	l, i := GetLabelNumber(s.GetNumber())
	if l == "" {
		return fmt.Sprintf("%03d", i)
	}
	return strings.ToUpper(fmt.Sprintf("%s-%03d", l, i))
}

func (s *FanzaScraper) GetCover() string {
	if s.schemaData == nil {
		return ""
	}
	return strings.Replace(s.schemaData.Image, "ps.jpg", "pl.jpg", 1)
}

func (s *FanzaScraper) GetPremiered() (rel string) {
	if s.schemaData == nil {
		return ""
	}
	rel = s.schemaData.SubjectOf.UploadDate
	if s.schemaData.SubjectOf.UploadDate == "" {
		rel = getDmmTableValue("発売日", s.doc)
		if rel == "" {
			rel = getDmmTableValue("配信開始日", s.doc)
		}
	}
	return rel
}

func (s *FanzaScraper) GetYear() (rel string) {
	if s.doc == nil {
		return ""
	}
	return regexp.MustCompile(`\d{4}`).FindString(s.GetPremiered())
}

func (s *FanzaScraper) NeedCut() bool {
	return needCut
}

type SchemaData struct {
	Context     string `json:"@context"`
	Type        string `json:"@type"`
	Name        string `json:"name"`
	Image       string `json:"image"`
	Description string `json:"description"`
	SKU         string `json:"sku"`
	Brand       struct {
		Type string `json:"@type"`
		Name string `json:"name"`
	} `json:"brand"`
	SubjectOf struct {
		Type         string `json:"@type"`
		Name         string `json:"name"`
		Description  string `json:"description"`
		ContentURL   string `json:"contentUrl"`
		ThumbnailURL string `json:"thumbnailUrl"`
		UploadDate   string `json:"uploadDate"`
		Actor        struct {
			Type          string `json:"@type"`
			Name          string `json:"name"`
			AlternateName string `json:"alternateName"`
		} `json:"actor"`
		Genre []string `json:"genre"`
	} `json:"subjectOf"`
}

func (s *FanzaScraper) parseSchema() error {
	if s.doc == nil {
		return nil
	}

	jsonStr := s.doc.Find("script[type='application/ld+json']").Text()
	return json.Unmarshal([]byte(jsonStr), &s.schemaData)
}
