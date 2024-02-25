package scraper

import (
	"dmm-scraper/pkg/config"
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
