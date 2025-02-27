package scraper

import (
	"dmm-scraper/pkg/config"
	"testing"
)

type args struct {
	query string
	url   string
}

type testCase struct {
	name    string
	args    args
	wantErr bool
	want    string
}

func BeforeTest() {
	c, err := config.NewLoader().LoadFile("../../config")
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	Setup(c)
}

func dumpS(t *testing.T, s Scraper) {
	got := s.GetNumber()
	t.Logf("GetNumber() = %v", got)

	got = s.GetPlot()
	t.Logf("GetPlot() = %v", got)
	got = s.GetTitle()
	t.Logf("GetTitle() = %v", got)
	got = s.GetCover()
	t.Logf("GetCover() = %v", got)
	got = s.GetDirector()
	t.Logf("GetDirector() = %v", got)
	got = s.GetMaker()
	t.Logf("GetMaker() = %v", got)
	got = s.GetLabel()
	t.Logf("GetLabel() = %v", got)
	got = s.GetRuntime()
	t.Logf("GetRuntime() = %v", got)
	got = s.GetPremiered()
	t.Logf("GetPremiered() = %v", got)
	gots := s.GetTags()
	t.Logf("GetTags() = %v", gots)
	gots = s.GetActors()
	t.Logf("GetActors() = %v", gots)
}
