package scraper

import (
	"regexp"
	"strings"
)

func GetQueryNum(name string) (query, num string, s Scraper) {
	typeSyndrome, _ := regexp.Compile(`Sex\sSyndrome`)
	typeSf, _ := regexp.Compile(`Sex\sFriend|webDL`)
	typeHeyzo, _ := regexp.Compile(`(heyzo|HEYZO)-[0-9]{4}`)
	typeFC2, _ := regexp.Compile(`(fc2|FC2|ppv|PPV)-[0-9]{6,7}`)
	typeMGStage, _ := regexp.Compile(`(siro|SIRO|[0-9]{3,4}[a-zA-Z]{2,5})-[0-9]{3,4}`)
	typeDmm, _ := regexp.Compile(`[a-zA-Z]{2,5}00[0-9]{3,4}`)
	typeDefault, _ := regexp.Compile(`[a-zA-Z]{2,5}-[0-9]{3,4}`)

	switch {
	case typeSyndrome.MatchString(name):
		jaChars := regexp.MustCompile("/[一-龠]+|[ぁ-ゔ]+|[ァ-ヴー]+|[々〆〤]+/u").FindAllString(name, -1)
		query = strings.Join(jaChars, " ")
		s = &GyuttoScraper{}
	case typeSf.MatchString(name):
		query = name
		s = &GyuttoScraper{}
	case typeHeyzo.MatchString(name):
		num = typeHeyzo.FindString(name)
		query = num
		s = &HeyzoScraper{}
	case typeFC2.MatchString(name):
		num = typeFC2.FindString(name)
		query = regexp.MustCompile(`[0-9]{6,7}`).FindString(name)
		s = &Fc2Scraper{}
	case typeMGStage.MatchString(name):
		num = typeMGStage.FindString(name)
		num = strings.ToUpper(num)
		query = num
		s = &MGStageScraper{}
	case typeDmm.MatchString(name):
		num = typeDmm.FindString(name)
		num = strings.Replace(num, "00", "-", 1)
		num = strings.ToUpper(num)
		query = num
		s = &DMMScraper{}
	default:
		num = typeDefault.FindString(name)
		num = strings.ToUpper(num)
		query = num
		s = &DMMScraper{}
	}

	return
}
