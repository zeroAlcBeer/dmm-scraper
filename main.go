package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	myclient "better-av-tool/internal/client"
	"better-av-tool/internal/configs"
	"better-av-tool/internal/img"
	"better-av-tool/internal/logger"
	"better-av-tool/internal/metadata"
	"better-av-tool/scraper"
)

var (
	conf         *configs.Configs
	client       myclient.Client
	imgOperation img.Operation
	posterWidth  = 378
	log          logger.Logger

	outputPath string
)

func init() {
	log = logger.New()

	var err error
	conf, err = configs.NewLoader().LoadFile("config")
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	client = myclient.New()
	if conf.Proxy.Enable {
		err = client.SetProxyUrl(conf.Proxy.Socket)
		if err != nil {
			log.Fatalf("Error load client proxy, %s", err)
		}
	}

	imgOperation = img.NewOperation()

	scraper.Init(client, log)
}

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
			err := s.FetchDoc(query)
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

			// 目录生成
			outputPath = metadata.NewOutputPath(s, conf.Output.Path)
			log.Infof("Making output path: %s", outputPath)
			err = os.MkdirAll(outputPath, 0700)
			if err != nil && !os.IsExist(err) {
				log.Errorf("Making output path err: %s", err)
				continue
			}

			// 做 nfo
			m := metadata.NewMovieNfo(s)

			// 下载封面
			poster := fmt.Sprintf("%s.jpg", num)
			posterPath := path.Join(outputPath, poster)
			err = client.Download(s.GetCover(), posterPath, myclient.DefaultProgress())
			if err != nil {
				log.Error(err)
				continue
			}

			m.SetPoster(poster)

			// 封面裁剪
			if s.NeedCut() {
				err := imgOperation.CropAndSave(posterPath, posterPath, posterWidth, 0)
				if err != nil {
					log.Error(err)
				}
			}

			// 写 nfo
			nfo := path.Join(outputPath, fmt.Sprintf("%s.nfo", num))
			err = m.Save(nfo)
			if err != nil {
				log.Error(err)
				continue
			}

			// 移动影片
			videoName := strings.ToUpper(num) + filepath.Ext(f.Name())
			err = os.Rename(f.Name(), path.Join(outputPath, videoName))
			if err != nil {
				log.Error(err)
			}
		}
	}
}
