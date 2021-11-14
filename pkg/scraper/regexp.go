package scraper

import (
	"bytes"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/dlclark/regexp2"
)

func GetQuery(name string) (query string, scrapers []Scraper) {
	typeGyutto, _ := regexp2.Compile(`item[0-9]{6,7}`, 0)
	isGyutto, _ := typeGyutto.MatchString(name)
	typeHeyzo, _ := regexp2.Compile(`(heyzo|HEYZO)-[0-9]{4}`, 0)
	isHeyzo, _ := typeHeyzo.MatchString(name)
	typeFC2, _ := regexp2.Compile(`(?<=(fc2|FC2|ppv|PPV)-)[0-9]{6,7}`, 0)
	isFC2, _ := typeFC2.MatchString(name)
	typeMGStage, _ := regexp2.Compile(`([0-9]{3,4}[a-zA-Z]{2,6})-[0-9]{3,4}`, 0)
	isMGStage, _ := typeMGStage.MatchString(name)
	// typeAnime, _ := regexp2.Compile(`(GLOD|JDXA|MJAD|ACRN|ORORE|DPLT)(-|)[0-9]{3,6}`, regexp2.RightToLeft)
	// isAnime, _ := typeDMM.MatchString(name)
	typeDMM, _ := regexp2.Compile(`[a-zA-Z]{2,5}(-|)[0-9]{3,6}`, regexp2.RightToLeft)
	isDMM, _ := typeDMM.MatchString(name)

	switch {
	case isGyutto:
		match, _ := typeGyutto.FindStringMatch(name)
		query = match.String()
		scrapers = append(scrapers, &GyuttoScraper{})
	case isHeyzo:
		match, _ := typeHeyzo.FindStringMatch(name)
		query = match.String()
		scrapers = append(scrapers, &HeyzoScraper{})
	case isFC2:
		match, _ := typeFC2.FindStringMatch(name)
		query = match.String()
		scrapers = append(scrapers, &Fc2Scraper{})
	case isMGStage:
		match, _ := typeMGStage.FindStringMatch(name)
		query = match.String()
		scrapers = append(scrapers, &MGStageScraper{})
	case isDMM:
		match, _ := typeDMM.FindStringMatch(name)
		query = match.String()
		if dmmProductService != nil {
			scrapers = append(scrapers, &DMMApiScraper{}, &DMMApiDigitalScraper{})
		} else {
			scrapers = append(scrapers, &DMMScraper{}, &FanzaScraper{})
		}
	}

	return
}

func GetLabelNumber(s string) (string, int) {
	labelRe, _ := regexp.Compile(`[a-zA-Z]{2,5}(-|)[0-9]{3,7}`)
	numRe, _ := regexp.Compile(`[0-9]{3,7}`)
	l := labelRe.FindString(s)
	if l == "" {
		i, _ := strconv.Atoi(numRe.FindString(s))
		return "", i
	}
	var label bytes.Buffer
	var number bytes.Buffer
	for _, c := range l {
		if unicode.IsLetter(c) {
			label.WriteRune(c)
		}
		if unicode.IsNumber(c) {
			number.WriteRune(c)
		}
	}
	l = strings.ToLower(label.String())
	i, _ := strconv.Atoi(number.String())
	return l, i
}

func isURLMatchQuery(u, query string) bool {
	var labelMatch bool
	var numMatch bool
	cid, _ := regexp.Compile(`cid=([^//]+)`)
	cids := cid.FindStringSubmatch(u)
	if len(cids) > 1 {
		label1, number1 := GetLabelNumber(cids[1])
		label2, number2 := GetLabelNumber(query)
		//fmt.Println(label1, number1)
		//fmt.Println(label2, number2)
		if label1 == label2 {
			labelMatch = true
		}

		if number1 == number2 {
			numMatch = true
		}
	}
	return labelMatch && numMatch
}
