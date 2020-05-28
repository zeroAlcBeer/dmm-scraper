package scraper

import (
	"better-av-tool/log"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"html"
	"regexp"
	"strings"
)

type HeyzoScraper struct {
	doc        *goquery.Document
	movieId    string
}

const (
	heyzoDetailUrl = "https://www.heyzo.com/moviepages/%s/index.html"
)

type HeyzoData struct {
	Context  string `json:"@context"`
	Type     string `json:"@type"`
	Name     string `json:"name"`
	Image    string `json:"image"`
	Encoding struct {
		Type           string `json:"@type"`
		EncodingFormat string `json:"encodingFormat"`
	} `json:"encoding"`
	Actor struct {
		Type  string `json:"@type"`
		Name  string `json:"name"`
		Image string `json:"image"`
	} `json:"actor"`
	Description   string `json:"description"`
	Duration      string `json:"duration"`
	DateCreated   string `json:"dateCreated"`
	ReleasedEvent struct {
		Type      string `json:"@type"`
		StartDate string `json:"startDate"`
		Location  struct {
			Type string `json:"@type"`
			Name string `json:"name"`
		} `json:"location"`
	} `json:"releasedEvent"`
	Video struct {
		Type         string `json:"@type"`
		Description  string `json:"description"`
		Duration     string `json:"duration"`
		Name         string `json:"name"`
		Thumbnail    string `json:"thumbnail"`
		ThumbnailURL string `json:"thumbnailUrl"`
		UploadDate   string `json:"uploadDate"`
		Actor        string `json:"actor"`
		Provider     string `json:"provider"`
	} `json:"video"`
	AggregateRating struct {
		Type        string `json:"@type"`
		RatingValue string `json:"ratingValue"`
		BestRating  string `json:"bestRating"`
		ReviewCount string `json:"reviewCount"`
	} `json:"aggregateRating"`
}

func (s *HeyzoScraper) FetchDoc(query, u string) (err error) {
	if u != "" {
		s.doc, err = GetDocFromUrl(u)
		return err
	}
	s.movieId = regexp.MustCompile("[0-9]+").FindString(query)

	u = fmt.Sprintf(heyzoDetailUrl, s.movieId)
	s.doc, err = GetDocFromUrl(u)
	return err
}

func (s *HeyzoScraper) GetPlot() string {
	if s.doc == nil {
		return ""
	}
	p := s.doc.Find("p[class=memo]").Text()
	return strings.TrimSpace(p)
}

func (s *HeyzoScraper) GetTitle() string {
	if s.doc == nil {
		return ""
	}
	j, _ := s.doc.Find("script[type='application/ld+json']").Html()
	d := &HeyzoData{}
	j = strings.ReplaceAll(html.UnescapeString(j), "\n", "")
	err := json.Unmarshal([]byte(j), d)
	if err != nil {
		log.Error(err)
	}
	return d.Name
}

func (s *HeyzoScraper) GetDirector() string {
	return ""
}

func (s *HeyzoScraper) GetRuntime() string {
	return ""
}

func (s *HeyzoScraper) GetTags() (tags []string) {
	if s.doc == nil {
		return
	}
	s.doc.Find(".tag-keyword-list a").Each(func(i int, ss *goquery.Selection) {
		tags = append(tags, strings.TrimSpace(ss.Text()))
	})
	return
}

func (s *HeyzoScraper) GetMaker() string {
	return "HEYZO"
}

func (s *HeyzoScraper) GetActors() (actors []string) {
	if s.doc == nil {
		return
	}
	s.doc.Find(".table-actor a").Each(func(i int, ss *goquery.Selection) {
		actors = append(actors, strings.TrimSpace(ss.Text()))
	})
	return
}

func (s *HeyzoScraper) GetLabel() string {
	return ""
}

func (s *HeyzoScraper) GetNumber() string {
	return s.movieId
}

func (s *HeyzoScraper) GetCover() string {
	if s.doc == nil {
		return ""
	}
	j, _ := s.doc.Find("script[type='application/ld+json']").Html()
	d := &HeyzoData{}
	j = strings.ReplaceAll(html.UnescapeString(j), "\n", "")
	err := json.Unmarshal([]byte(j), d)
	if err != nil {
		log.Error(err)
	}
	return strings.ReplaceAll(d.Image, "//", "https://")
}

func (s *HeyzoScraper) GetWebsite() string {
	return s.doc.Url.String()
}

func (s *HeyzoScraper) GetPremiered() (rel string) {
	if s.doc == nil {
		return ""
	}
	rel = s.doc.Find(".table-release-day td").Last().Text()
	return strings.Replace(strings.TrimSpace(rel), "/", "-", -1)
}

func (s *HeyzoScraper) GetYear() (rel string) {
	if s.doc == nil {
		return ""
	}
	return regexp.MustCompile(`\d{4}`).FindString(s.GetPremiered())
}

func (s *HeyzoScraper) GetSeries() string {
	if s.doc == nil {
		return ""
	}
	p := s.doc.Find(".table-series td").Last().Text()
	if strings.TrimSpace(p) == "-----" {
		return ""
	}
	return strings.TrimSpace(p)
}

func (s *HeyzoScraper) NeedCut() bool {
	return false
}
