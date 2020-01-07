package scraper

type Scraper interface{

	// Methods
	GetTitle() string
	GetStudio() float64
	GetOutline() float64
}

//fmt.Println(getTitle(doc))
//fmt.Println(getTableValue("メーカー", doc))
//fmt.Println(getTableValue("発売日", doc))
//fmt.Println(getDesc(doc))
//fmt.Println(getActors(doc))
//fmt.Println(getTags(doc))
//fmt.Println(getCover(doc))
//
//fmt.Println(getTableValue("シリーズ", doc))
//fmt.Println(getTableValue("収録時間", doc))
//fmt.Println(getTableValue("監督", doc))
//fmt.Println(getTableValue("レーベル", doc))
