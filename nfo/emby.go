package nfo

import (
	"better-av-tool/scraper"
	"encoding/xml"
	"strings"
)

type EmbyMovie struct {
	XMLName   xml.Name `xml:"movie"`
	Text      string   `xml:",chardata"`
	Plot      string   `xml:"plot"`
	Outline   string   `xml:"outline"`
	Title     string   `xml:"title"`
	Director  string   `xml:"director"`
	Year      string   `xml:"year"`
	Premiered string   `xml:"premiered"`
	//ReleaseDate string   `xml:"releasedate"`
	Runtime string           `xml:"runtime"`
	Genre   []string         `xml:"genre"`
	Studio  string           `xml:"studio"`
	Tag     []string         `xml:"tag"`
	Actor   []EmbyMovieActor `xml:"actor"`
	Poster  string           `xml:"poster"`
	Fanart  []EmbyMovieThumb `xml:"fanart"`
	Label   string           `xml:"label"`
	Num     string           `xml:"num"`
	Cover   string           `xml:"cover"`
	Website string           `xml:"website"`
}

type EmbyMovieActor struct {
	Text string `xml:",chardata"`
	Name string `xml:"name"`
}

type EmbyMovieThumb struct {
	Text  string `xml:",chardata"`
	Thumb string `xml:"thumb"`
}

func Build(s scraper.Scraper, num string) ([]byte, error) {
	//releaseTime, err := time.Parse("2006/01/02", )
	//if err != nil {
	//	return
	//}
	var actors []EmbyMovieActor
	for _, name := range s.GetActors() {
		actors = append(actors, EmbyMovieActor{Name: name})
	}

	m := &EmbyMovie{
		Plot:      s.GetPlot(),
		Title:     num + " " + s.GetTitle(),
		Director:  s.GetDirector(),
		Year:      "",
		Premiered: strings.Replace(s.GetPremiered(), "/", "-", -1),
		Runtime:   s.GetRuntime(),
		Genre:     append(s.GetTags(), s.GetSeries(), s.GetLabel()),
		Tag:       append(s.GetTags(), s.GetSeries(), s.GetLabel()),
		Studio:    s.GetMaker(),
		Label:     s.GetLabel(),
		Actor:     actors,
		Cover:     s.GetCover(),
		Num:       s.GetNumber(),
		Website:   s.GetWebsite(),
		Fanart:    []EmbyMovieThumb{{Thumb: num + ".jpg"}},
		Poster:    num + ".png",
	}

	if len(s.GetPremiered()) >= 4 {
		m.Year = s.GetPremiered()[:4]
	}
	x, err := xml.MarshalIndent(m, "", "  ")
	if err != nil {
		return x, err
	}
	x = []byte(xml.Header + string(x))
	return x, nil
}
