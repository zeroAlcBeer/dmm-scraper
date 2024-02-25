package scraper

import (
	"errors"
	"strings"

	"dmm-scraper/third_party/dmm-go-sdk/api"
)

type DMMApiDigitalScraper struct {
	DMMApiScraper
}

func (s *DMMApiDigitalScraper) GetType() string {
	return "DMMApiDigitalScraper"
}

// FetchDoc
// ...
func (s *DMMApiDigitalScraper) FetchDoc(query string) (err error) {
	query = strings.Replace(query, "-", "00", 1)

	dmmProductService.SetSite(api.SiteAdult)
	dmmProductService.SetService("digital")
	//dmmapi.SetFloor("videoa")
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
