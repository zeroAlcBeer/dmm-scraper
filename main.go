package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	myclient "better-av-tool/pkg/client"
	"better-av-tool/pkg/config"
	"better-av-tool/pkg/img"
	"better-av-tool/pkg/logger"
	"better-av-tool/pkg/metadata"
	"better-av-tool/pkg/scraper"
)

var (
	posterWidth = 378
	outputPath string
)

func isValidVideo(ext string) bool {
	switch strings.ToLower(ext) {
	case
		".wmv",
		".mp4",
		".avi",
		".mkv":
		return true
	}
	return false
}

func main() {
	var err error
	log := logger.New()

	conf, err := config.NewLoader().LoadFile("config")
	if err != nil {
		log.Errorf("Error reading config file, %s", err)
		log.Warnf("Loading default config")
		conf = config.Default()
	}

	scraper.Setup(conf.Proxy)

	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		ext := filepath.Ext(f.Name())
		if !isValidVideo(ext) {
			continue
		}
		log.Infof("Check file: %s", f.Name())
		name := strings.TrimSuffix(f.Name(), ext)

		// 用正则处理文件名
		if query, s := scraper.GetQuery(name); query != "" {
			log.Infof("Scraper get query: %s", query)

			// 爬取页面
			err = s.FetchDoc(query)
			if err != nil {
				log.Error(err)
				continue
			}
			log.Infof("Scraper get number: %s", s.GetNumber())
			if s.GetNumber() == "" {
				log.Error("scraper get number empty")
				continue
			}
			num := scraper.FormatNum(s.GetNumber())
			log.Infof("Scraper get num format: %s", num)

			// mkdir
			outputPath = metadata.NewOutputPath(s, conf.Output.Path)
			log.Infof("Making output path: %s", outputPath)
			err = os.MkdirAll(outputPath, 0700)
			if err != nil && !os.IsExist(err) {
				log.Errorf("Making output path err: %s", err)
				continue
			}

			// build nfo
			m := metadata.NewMovieNfo(s)

			// download cover
			poster := fmt.Sprintf("%s.jpg", num)
			posterPath := path.Join(outputPath, poster)
			err = scraper.Download(s.GetCover(), posterPath, myclient.DefaultProgress())
			if err != nil {
				log.Error(err)
				continue
			}

			m.SetPoster(poster)

			// cut cover
			if s.NeedCut() {
				imgOperation := img.NewOperation()
				err = imgOperation.CropAndSave(posterPath, posterPath, posterWidth, 0)
				if err != nil {
					log.Error(err)
				}
			}

			// write nfo file
			nfo := path.Join(outputPath, fmt.Sprintf("%s.nfo", num))
			err = m.Save(nfo)
			if err != nil {
				log.Error(err)
				continue
			}

			// move file
			videoName := strings.ToUpper(num) + filepath.Ext(f.Name())
			err = os.Rename(f.Name(), path.Join(outputPath, videoName))
			if err != nil {
				log.Error(err)
			}
		}
	}
}
