package scraper

import (
	"testing"

	"github.com/PuerkitoBio/goquery"
)

var dmmTests map[string]*DMMScraper
var fc2Tests map[string]*Fc2Scraper
var mgsTests map[string]*MGStageScraper
var heyzoTests map[string]*HeyzoScraper
var gyuttoTests map[string]*GyuttoScraper

type fields struct {
	doc *goquery.Document
}

type args struct {
	query string
	url   string
}

type testCase struct {
	name    string
	fields  fields
	args    args
	wantErr bool
	want    string
}

func TestMain(m *testing.M) {
	fc2Tests = make(map[string]*Fc2Scraper)
	mgsTests = make(map[string]*MGStageScraper)
	heyzoTests = make(map[string]*HeyzoScraper)
	gyuttoTests = make(map[string]*GyuttoScraper)
	dmmTests = make(map[string]*DMMScraper)

	m.Run()
}
