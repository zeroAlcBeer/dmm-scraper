package scraper

// Scraper is interface
type Scraper interface {
	// Remote
	FetchDoc(query string) (err error)

	// Local
	GetPlot() string
	GetTitle() string
	GetDirector() string
	GetRuntime() string
	GetTags() []string
	GetMaker() string
	GetActors() []string
	GetLabel() string
	GetNumber() string
	GetCover() string
	GetWebsite() string
	GetPremiered() string
	GetYear() string
	GetSeries() string

	// Operation
	NeedCut() bool
}
