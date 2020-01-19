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
	flag.StringVar(&proxyUrl, "proxy", "", "set proxy url")
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

	var ff = func(pathX string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !IsValidVideo(filepath.Ext(pathX)) {
			return nil
		}

		var num string
		var s scraper.Scraper

		typeHeyzo, _ := regexp.Compile(`(heyzo|HEYZO)-[0-9]{4}`)
		typeFc2, _ := regexp.Compile(`(fc2|FC2|ppv|PPV)-[0-9]{6,7}`)
		typeMGStage, _ := regexp.Compile(`(siro|SIRO|[0-9]{3,4}[a-zA-Z]{2,5})-[0-9]{3,4}`)
		typeDmm, _ := regexp.Compile(`[a-zA-Z]{2,5}00[0-9]{3,4}`)
		typeDefault, _ := regexp.Compile(`[a-zA-Z]{2,5}-[0-9]{3,4}`)

		switch {
		case typeHeyzo.MatchString(info.Name()):
			num = typeHeyzo.FindString(info.Name())
			s = &scraper.HeyzoScraper{}
		case typeFc2.MatchString(info.Name()):
			num = typeFc2.FindString(info.Name())
			s = &scraper.Fc2Scraper{}
		case typeMGStage.MatchString(info.Name()):
			num = typeMGStage.FindString(info.Name())
			s = &scraper.MGStageScraper{}
		case typeDmm.MatchString(info.Name()):
			num = typeDmm.FindString(info.Name())
			num = strings.Replace(num, "00", "-", 1)
			s = &scraper.DMMScraper{}
		default:
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
				Name:    jpgName,
				Url:     s.GetCover(),
				NeedCut: s.NeedCut(),
			}
			if _, err := os.Stat(item.Name); os.IsNotExist(err) {
				go download(item)
			}

			newPath := strings.ToUpper(num) + filepath.Ext(pathX)
			err = os.Rename(pathX, path.Join(outputPath, s.GetPremiered()[:4], newPath))
			if err != nil {
				return err
			}
		}

		return nil
	}

	_ = filepath.Walk(scanPath, ff)
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
	Name    string
	Url     string
	NeedCut bool
}

func download(item DownloadItem) {
	log.Infof("Downloading from %s", item.Url)

	req, _ := grab.NewRequest(item.Name, item.Url)
	resp := grabClient.Do(req)

	if err := resp.Err(); err != nil {
		log.Errorf("%s: %v", resp.Filename, err)
	}

	if resp.IsComplete() {
		pngName := strings.ReplaceAll(resp.Filename, "jpg", "png")
		_ = cropOrCopy(resp.Filename, pngName, item.NeedCut)
		log.Infof("Finished %s %d / %d bytes (%d%%)", resp.Filename, resp.BytesComplete(), resp.Size, int(100*resp.Progress()))
	}
	return
}

func cropOrCopy(jpgName, pngName string, needCut bool) error {
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
	if needCut {
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

func IsValidVideo(ext string) bool {
	switch ext {
	case
		".wmv",
		".WMV",
		".mp4",
		".MP4",
		".avi",
		".AVI",
		".mkv",
		".MKV":
		return true
	}
	return false
}
