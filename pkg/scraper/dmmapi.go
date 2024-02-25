package scraper

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"dmm-scraper/third_party/dmm-go-sdk/api"
)

type DMMApiScraper struct {
	DefaultScraper
	Item api.Item
}

func (s *DMMApiScraper) GetType() string {
	return "DMMApiScraper"
}

// FetchDoc
// ...
func (s *DMMApiScraper) FetchDoc(query string) (err error) {
	query = strings.Replace(query, "-", "", 1)

	dmmProductService.SetSite(api.SiteAdult)
	dmmProductService.SetService("mono")
	//dmmapi.SetFloor("dvd")
	dmmProductService.SetKeyword(query)
	result, err := dmmProductService.Execute()
	if result.TotalCount == 1 {
		s.Item = result.Items[0]
		return nil
	}
	if result.TotalCount == 0 {
		return errors.New("record not found")
	}
	if result.ResultCount > 1 {
		for _, item := range result.Items {
			if isURLMatchQuery(item.URL, query) {
				s.Item = item
				break
			}
		}
	}
	return err
}

func (s *DMMApiScraper) GetPlot() string {
	return ""
}

func (s *DMMApiScraper) GetTitle() string {
	return s.Item.Title
}

func (s *DMMApiScraper) GetDirector() string {
	if len(s.Item.ItemInformation.Directors) > 0 {
		return s.Item.ItemInformation.Directors[0].Name
	}
	return ""
}

func (s *DMMApiScraper) GetRuntime() string {
	return s.Item.Volume
}

func (s *DMMApiScraper) GetTags() (tags []string) {
	for _, genre := range s.Item.ItemInformation.Genres {
		tags = append(tags, genre.Name)
	}
	return
}

func (s *DMMApiScraper) GetMaker() string {
	if len(s.Item.ItemInformation.Maker) > 0 {
		return s.Item.ItemInformation.Maker[0].Name
	}
	return ""
}

func (s *DMMApiScraper) GetActors() (actors []string) {
	for _, actor := range s.Item.ItemInformation.Actress {
		actors = append(actors, actor.Name)
	}
	return
}

func (s *DMMApiScraper) GetLabel() string {
	if len(s.Item.ItemInformation.Label) > 0 {
		return s.Item.ItemInformation.Label[0].Name
	}
	return ""
}

func (s *DMMApiScraper) GetNumber() string {
	return s.Item.ContentID
}

func (s *DMMApiScraper) GetFormatNumber() string {
	l, i := GetLabelNumber(s.GetNumber())
	if l == "" {
		return fmt.Sprintf("%03d", i)
	}
	return strings.ToUpper(fmt.Sprintf("%s-%03d", l, i))
}

func (s *DMMApiScraper) GetCover() string {
	return s.Item.ImageURL.Large
}

func (s *DMMApiScraper) GetPremiered() (rel string) {
	return s.Item.Date
}

func (s *DMMApiScraper) GetYear() (rel string) {
	return regexp.MustCompile(`\d{4}`).FindString(s.GetPremiered())
}

func (s *DMMApiScraper) GetSeries() string {
	if len(s.Item.ItemInformation.Series) > 0 {
		return s.Item.ItemInformation.Series[0].Name
	}
	return ""
}

func (s *DMMApiScraper) GetWebsite() string {
	return s.Item.URL
}

func (s *DMMApiScraper) NeedCut() bool {
	return needCut
}
