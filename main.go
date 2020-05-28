package main

import (
	"better-av-tool/grabber"
	"better-av-tool/log"
	"better-av-tool/nfo"
	"better-av-tool/scraper"

	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"golang.org/x/net/proxy"
)

var (
	conf        Conf
	outputPath  string
	proxyClient *http.Client
)

type OutputConf struct {
	Path string
}

type ProxyConf struct {
	Enable bool
	Socket string
}

type Conf struct {
	Output OutputConf
	Proxy  ProxyConf
}

func init() {
	initConf()
	initProxy()
	initGrabber()
	initScraper()
}

func initConf() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	err := viper.Unmarshal(&conf)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
}

func initProxy() {
	proxyClient = &http.Client{}
	if conf.Proxy.Enable {
		url, err := url.Parse(conf.Proxy.Socket)
		if err != nil {
			log.Fatal(err)
		}
		dialer, err := proxy.FromURL(url, proxy.Direct)
		if err != nil {
			log.Fatal(err)
		}
		proxyClient.Transport = &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
				c, e := dialer.Dial(network, addr)
				return c, e
			},
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
}

func initGrabber() {
	grabber.SetHTTPClient(proxyClient)
}

func initScraper() {
	scraper.SetHTTPClient(proxyClient)
}

func ensureDir(dirName string) error {
	if _, err := os.Stat(dirName); err == nil {
		return nil
	}
	err := os.MkdirAll(dirName, os.ModeDir)
	if err == nil || os.IsExist(err) {
		return nil
	} else {
		return err
	}
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
	//if err := ensureDir(conf.Output.Path); err != nil {
	//	log.Fatal(err)
	//}

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
		log.Infof("Check file %s", f.Name())
		name := strings.TrimSuffix(f.Name(), ext)

		// 用正则处理文件名
		if num, s := scraper.GetNum(name); num != "" {
			log.Infof("Match num %s!", num)

			// 爬取页面
			err := s.FetchDoc(num, "")
			if err != nil {
				log.Error(err)
				continue
			}

			// 目录生成
			outputPath = scraper.ParsePath(s, conf.Output.Path)
			log.Infof("Making output path %s", outputPath)
			err = ensureDir(outputPath)
			if err != nil {
				log.Error(err)
				continue
			}
			// 做 nfo
			model := nfo.Build(s)

			// 下载封面
			coverSrc, err := grabber.Download(s.GetCover())
			if err != nil {
				log.Error(err)
				continue
			}
			coverDst := path.Join(outputPath, fmt.Sprintf("%s.jpg", num))
			_ = os.Rename(coverSrc, coverDst)
			model.Fanart = []nfo.EmbyMovieThumb{{Thumb: path.Base(coverDst)}}
			model.Poster = path.Base(coverDst)

			// 封面裁剪
			if s.NeedCut() {
				posterDst := strings.Replace(coverDst, "jpg", "png", 1)
				err = grabber.Crop(coverDst, posterDst)
				if err != nil {
					log.Error(err)
					continue
				}
				model.Poster = path.Base(posterDst)
			}

			// 写 nfo
			nfoName := path.Join(outputPath, fmt.Sprintf("%s.nfo", num))
			err = model.WriteFile(nfoName)
			if err != nil {
				log.Error(err)
				continue
			}

			// 影片重命名
			videoName := strings.ToUpper(num) + filepath.Ext(f.Name())
			err = os.Rename(f.Name(), path.Join(outputPath, videoName))
			if err != nil {
				log.Error(err)
			}
		}
	}
}
