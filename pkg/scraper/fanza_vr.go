package scraper

type FanzaVRScraper struct {
	FanzaScraper
}

func (s *FanzaVRScraper) GetType() string {
	return "FanzaVRScraper"
}

func (s *FanzaVRScraper) NeedCut() bool {
	return false
}
