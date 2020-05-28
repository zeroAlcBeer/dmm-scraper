package scraper

import (
	"regexp"
	"strings"
)

func GetNum(name string) (num string, s Scraper) {
	typeSyndrome, _ := regexp.Compile(`Sex\sSyndrome`)
	typeSf, _ := regexp.Compile(`Sex\sFriend`)
	typeHeyzo, _ := regexp.Compile(`(heyzo|HEYZO)-[0-9]{4}`)
	typeFc2, _ := regexp.Compile(`(fc2|FC2|ppv|PPV)-[0-9]{6,7}`)
	typeMGStage, _ := regexp.Compile(`(siro|SIRO|[0-9]{3,4}[a-zA-Z]{2,5})-[0-9]{3,4}`)
	typeDmm, _ := regexp.Compile(`[a-zA-Z]{2,5}00[0-9]{3,4}`)
	typeDefault, _ := regexp.Compile(`[a-zA-Z]{2,5}-[0-9]{3,4}`)

	switch {
	case typeSyndrome.MatchString(name):
		num = name
		jaChars := regexp.MustCompile("/[一-龠]+|[ぁ-ゔ]+|[ァ-ヴー]+|[々〆〤]+/u").FindAllString(name, -1)
		num = strings.Join(jaChars, " ")
		s = &GyuttoScraper{}
	case typeSf.MatchString(name):
		num = name
		s = &GyuttoScraper{}
	case typeHeyzo.MatchString(name):
		num = typeHeyzo.FindString(name)
		num = strings.ToUpper(num)
		s = &HeyzoScraper{}
	case typeFc2.MatchString(name):
		num = typeFc2.FindString(name)
		num = regexp.MustCompile(`[0-9]{6,7}`).FindString(num)
		s = &Fc2Scraper{}
	case typeMGStage.MatchString(name):
		num = typeMGStage.FindString(name)
		num = strings.ToUpper(num)
		s = &MGStageScraper{}
	case typeDmm.MatchString(name):
		num = typeDmm.FindString(name)
		num = strings.Replace(num, "00", "-", 1)
		num = strings.ToUpper(num)
		s = &DMMScraper{}
	default:
		num = typeDefault.FindString(name)
		num = strings.ToUpper(num)
		s = &DMMScraper{}
	}

	return
}
