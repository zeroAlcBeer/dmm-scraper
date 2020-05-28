package scraper

import "strings"

func ParsePath(s Scraper, p string) string {
	p = strings.Replace(p, "{year}", s.GetYear(), 1)
	if len(s.GetActors()) > 0 {
		p = strings.Replace(p, "{actor}", s.GetActors()[0], 1)
	} else {
		p = strings.Replace(p, "{actor}", "", 1)
	}
	p = strings.Replace(p, "{maker}", s.GetMaker(), 1)
	p = strings.Replace(p, "{num}", s.GetNumber(), 1)

	return strings.Replace(p, "//", "/", -1)
}
