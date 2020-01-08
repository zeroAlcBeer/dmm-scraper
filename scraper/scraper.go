package scraper

type Scraper interface {
	// Remote
	FetchDoc(num string) error

	// Local
	GetPlot() string
	GetTitle() string
	GetDirector() string
	GetRuntime() string
	GetTags() []string
	GetMaker() string
	//GetOutline() string
	GetActors() []string
	GetLabel() string
	GetNumber() string
	GetCover() string
	GetWebsite() string
	GetPremiered() string
	GetSeries() string
}
