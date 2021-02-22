package metadata

import (
	"better-av-tool/scraper"
	"encoding/xml"
	"os"
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

func (m *EmbyMovie) ToXML() ([]byte, error) {
	x, err := xml.MarshalIndent(m, "", "  ")
	if err != nil {
		return x, err
	}
	x = []byte(xml.Header + string(x))
	return x, nil
}

func (m *EmbyMovie) Save(filename string) error {
	b, err := m.ToXML()
	if err != nil {
		return err
	}
	return os.WriteFile(filename, b, 0644)
}

func (m *EmbyMovie) SetPoster(filename string) {
	m.Fanart = []EmbyMovieThumb{{Thumb: filename}}
	m.Poster = filename
}

func newEmbyMovieActors(names []string) []EmbyMovieActor {
	var actors []EmbyMovieActor
	for _, name := range names {
		actors = append(actors, EmbyMovieActor{Name: name})
	}
	return actors
}

// NewMovieNfo ...
func NewMovieNfo(s scraper.Scraper) MovieNfo {
	return &EmbyMovie{
		Plot:      s.GetPlot(),
		Title:     s.GetTitle(),
		Director:  s.GetDirector(),
		Year:      s.GetYear(),
		Premiered: s.GetPremiered(),
		Runtime:   s.GetRuntime(),
		Genre:     append(s.GetTags(), s.GetSeries(), s.GetLabel()),
		Tag:       append(s.GetTags(), s.GetSeries(), s.GetLabel()),
		Studio:    s.GetMaker(),
		Label:     s.GetLabel(),
		Actor:     newEmbyMovieActors(s.GetActors()),
		Cover:     s.GetCover(),
		Num:       s.GetNumber(),
		Website:   s.GetWebsite(),
		//Fanart:    []EmbyMovieThumb{{Thumb: cover}},
		//Poster:    "",
	}
}

func NewOutputPath(s scraper.Scraper, format string) string {
	p := strings.Replace(format, "{year}", s.GetYear(), 1)
	if len(s.GetActors()) > 0 {
		p = strings.Replace(p, "{actor}", s.GetActors()[0], 1)
	} else {
		p = strings.Replace(p, "{actor}", "", 1)
	}
	p = strings.Replace(p, "{maker}", s.GetMaker(), 1)
	p = strings.Replace(p, "{num}", s.GetNumber(), 1)

	return strings.Replace(p, "//", "/", -1)
}
