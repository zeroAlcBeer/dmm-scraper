package main

import (
	"better-av-tool/log"
	"better-av-tool/nfo"
	"better-av-tool/scraper"
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/cavaliercoder/grab"
	"github.com/oliamb/cutter"
	"golang.org/x/net/proxy"
)

var (
	scanPath   string
	outputPath string
	proxyUrl   string

	proxyClient *http.Client
	grabClient  *grab.Client
)

func init() {
	flag.StringVar(&scanPath, "path", "", "set scan path")
	flag.StringVar(&outputPath, "output", "", "set output path")
	flag.StringVar(&proxyUrl, "url", "", "set proxy url")
	flag.Parse()

}

func main() {
	if proxyUrl != "" {
		url, err := url.Parse(proxyUrl)
		if err != nil {
			log.Fatal(err)
		}
		dialer, err := proxy.FromURL(url, proxy.Direct)
		if err != nil {
			log.Fatal(err)
		}
		proxyClient = &http.Client{}
		proxyClient.Transport = &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
				c, e := dialer.Dial(network, addr)
				return c, e
			},
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		grabClient = grab.NewClient()
		if proxyClient != nil {
			grabClient.HTTPClient = proxyClient
		}
		grabClient.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.88 Safari/537.36"

	}

	if outputPath == "" {
		log.Fatal("output path empty")
	}
	if err := ensureDir(outputPath); err != nil {
		log.Fatal(err)
	}

	items := make([]DownloadItem, 0)
	var ff = func(pathX string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		//find out if it's a dir or file, if file, print info
		if !info.IsDir() {
			var num string
			var s scraper.Scraper


			typeMGStage, _ := regexp.Compile(`[0-9]{3,4}[a-zA-Z]{2,5}-[0-9]{3,4}`)
			typeDefault, _ := regexp.Compile(`[a-zA-Z]{2,5}-[0-9]{3,4}`)
			if typeMGStage.MatchString(info.Name()) {
				num = typeMGStage.FindString(info.Name())
				s = &scraper.MGStageScraper{}
			} else {
				num = typeDefault.FindString(info.Name())
				s = &scraper.DMMScraper{}
			}
			if num != "" {
				num = strings.ToUpper(num)
				log.Infof("num %s match!", num)
				b, err := GetNfo(s, num)
				if err != nil {
					log.Error(err)
					return nil
				}
				nfoName := path.Join(outputPath, s.GetPremiered()[:4], fmt.Sprintf("%s.nfo", num))
				jpgName := path.Join(outputPath, s.GetPremiered()[:4], fmt.Sprintf("%s.jpg", num))
				err = ensureDir(filepath.Dir(nfoName))
				if err != nil {
					log.Error(err)
					return nil
				}
				err = ioutil.WriteFile(nfoName, b, 0644)
				if err != nil {
					log.Error(err)
					return nil
				}
				log.Infof("%s built!", nfoName)
				item := DownloadItem{
					Name: jpgName,
					Url:  s.GetCover(),
				}
				if _, err := os.Stat(item.Name); os.IsNotExist(err) {
					items = append(items, item)
				}

				newPath := strings.ToUpper(num) + filepath.Ext(pathX)
				err = os.Rename(pathX, path.Join(outputPath, s.GetPremiered()[:4], newPath))
				if err != nil {
					return err
				}
			}
		}

		return nil
	}

	_ = filepath.Walk(scanPath, ff)
	DownloadFiles(items)

}

func GetNfo(s scraper.Scraper, num string) ([]byte, error) {
	s.SetHTTPClient(proxyClient)
	err := s.FetchDoc(num)
	if err != nil {
		return nil, err
	}
	return nfo.Build(s, num)
}

type DownloadItem struct {
	Name string
	Url  string
}

func DownloadFiles(items []DownloadItem) {
	// create grabClient
	grabClient := grab.NewClient()
	grabClient.HTTPClient = proxyClient
	grabClient.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.88 Safari/537.36"

	reqs := make([]*grab.Request, 0)
	for _, item := range items {
		req, _ := grab.NewRequest(item.Name, item.Url)
		reqs = append(reqs, req)
	}
	log.Infof("Downloading %d files...", len(reqs))
	respCh := grabClient.DoBatch(8, reqs...)

	var success int
	for resp := range respCh {
		if err := resp.Err(); err != nil {
			log.Errorf("%s: %v", resp.Filename, err)
		}

		if resp.IsComplete() {
			success++
			_ = CropCover(resp.Filename, strings.Replace(resp.Filename, "jpg", "png", 1))
			log.Infof("Finished %s %d / %d bytes (%d%%)", resp.Filename, resp.BytesComplete(), resp.Size, int(100*resp.Progress()))
		}
	}
	log.Infof("%d files successfully downloaded.", success)
	return
}

func CropCover(jpgName, pngName string) error {
	f, err := os.Open(jpgName)
	if err != nil {
		log.Error("Cannot open file", err)
		return err
	}
	defer f.Close()
	img, err := jpeg.Decode(f)
	if err != nil {
		log.Error("Cannot decode image:", err)
		return err
	}
	srcW := img.Bounds().Dx()
	srcH := img.Bounds().Dy()
	if srcW == 800 {
		img, err = cutter.Crop(img, cutter.Config{
			Height:  srcH,                       // height in pixel or Y ratio(see Ratio Option below)
			Width:   378,                        // width in pixel or X ratio
			Mode:    cutter.TopLeft,             // Accepted Mode: TopLeft, Centered
			Anchor:  image.Point{srcW - 378, 0}, // Position of the top left point
			Options: 0,                          // Accepted Option: Ratio
		})

		if err != nil {
			log.Error("Cannot Crop image:", err)
			return err
		}
	}

	out, err := os.Create(pngName)
	if err != nil {
		log.Error("Cannot create image:", err)
		return err
	}
	defer out.Close()
	err = png.Encode(out, img)

	if err != nil {
		log.Error("Cannot Encode image:", err)
		return err
	}

	return nil
}

func ensureDir(dirName string) error {
	err := os.Mkdir(dirName, os.ModeDir)
	if err == nil || os.IsExist(err) {
		return nil
	} else {
		return err
	}
}
