package scraper

import (
	myclient "better-av-tool/pkg/client"
	"better-av-tool/pkg/config"
	"better-av-tool/pkg/logger"
)

// Scraper is interface
type Scraper interface {
	FetchDoc(query string) (err error)
	GetPlot() string
	GetTitle() string
	GetDirector() string
	GetRuntime() string
	GetTags() []string
	GetMaker() string
	GetActors() []string
	GetLabel() string
	GetNumber() string
	GetFormatNumber() string
	GetCover() string
	GetWebsite() string
	GetPremiered() string
	GetYear() string
	GetSeries() string
	GetType() string
	NeedCut() bool
}

var (
	client myclient.Client
	log    logger.Logger
)

// Setup ...
func Setup(p config.Proxy) {
	log = logger.New()
	client = myclient.New()
	if p.Enable {
		err := client.SetProxyUrl(p.Socket)
		if err != nil {
			log.Errorf("Error parse proxy url, %s, proxy disabled", err)
		}
	}
}
