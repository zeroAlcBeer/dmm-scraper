package scraper

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

func GetQuery(name string) (query string, s Scraper) {
	typeSyndrome, _ := regexp.Compile(`Sex\sSyndrome`)
	typeSf, _ := regexp.Compile(`Sex\sFriend|webDL`)
	typeHeyzo, _ := regexp.Compile(`(heyzo|HEYZO)-[0-9]{4}`)
	typeFC2, _ := regexp.Compile(`(fc2|FC2|ppv|PPV)-[0-9]{6,7}`)
	typeMGStage, _ := regexp.Compile(`(siro|SIRO|[0-9]{3,4}[a-zA-Z]{2,5})-[0-9]{3,4}`)
	typeDefault, _ := regexp.Compile(`[a-zA-Z]{2,5}(-|)[0-9]{3,5}`)

	switch {
	case typeSyndrome.MatchString(name):
		jaChars := regexp.MustCompile("/[一-龠]+|[ぁ-ゔ]+|[ァ-ヴー]+|[々〆〤]+/u").FindAllString(name, -1)
		query = strings.Join(jaChars, " ")
		s = &GyuttoScraper{}
	case typeSf.MatchString(name):
		query = name
		s = &GyuttoScraper{}
	case typeHeyzo.MatchString(name):
		query = typeHeyzo.FindString(name)
		s = &HeyzoScraper{}
	case typeFC2.MatchString(name):
		query = regexp.MustCompile(`[0-9]{6,7}`).FindString(name)
		s = &Fc2Scraper{}
	case typeMGStage.MatchString(name):
		query = typeMGStage.FindString(name)
		s = &MGStageScraper{}
	default:
		query = typeDefault.FindString(name)
		s = &DMMScraper{}
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

func FormatNum(s string) string {
	l, i := GetLabelNumber(s)
	if l == "" {
		return fmt.Sprintf("%03d", i)
	}
	return strings.ToUpper(l) + "-" + fmt.Sprintf("%03d", i)
}
