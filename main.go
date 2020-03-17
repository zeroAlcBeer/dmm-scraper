package main

import (
	"better-av-tool/log"
	"better-av-tool/nfo"
	"better-av-tool/scraper"
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/cavaliercoder/grab"
	"github.com/oliamb/cutter"
	"golang.org/x/net/proxy"
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
)

var (
	basePath   = "output"
	outputPath string
	proxyUrl   string

	proxyClient *http.Client
	grabClient  *grab.Client
)

func init() {
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

	if err := ensureDir("output"); err != nil {
		log.Fatal(err)
	}

	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if !IsValidVideo(filepath.Ext(f.Name())) {
			continue
		}

		if num, s := GetNum(f.Name()); num != "" {
			log.Infof("Match num %s!", num)
			b, err := GetNfo(num, s)
			if err != nil {
				log.Error(err)
				continue
			}

			outputPath = basePath
			err = MakeOutputPath(s)
			if err != nil {
				log.Error(err)
				continue
			}

			err = BuildNfo(string(b), num)
			if err != nil {
				log.Error(err)
				continue
			}

			err = DownCover(DownloadItem{
				Name:    path.Join(outputPath, fmt.Sprintf("%s.jpg", num)),
				Url:     s.GetCover(),
				NeedCut: s.NeedCut(),
			})
			if err != nil {
				log.Error(err)
				continue
			}

			err = MoveFile(num, f.Name())
			if err != nil {
				log.Error(err)
				continue
			}
		}
	}
}

func GetNfo(num string, s scraper.Scraper) ([]byte, error) {
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

func DownCover(item DownloadItem) error {
	log.Infof("Downloading from %s", item.Url)

	req, _ := grab.NewRequest(item.Name, item.Url)
	resp := grabClient.Do(req)

	if err := resp.Err(); err != nil {
		return err
	}

	if resp.IsComplete() {
		pngName := strings.ReplaceAll(resp.Filename, "jpg", "png")
		_ = cropOrCopy(resp.Filename, pngName, item.NeedCut)
		log.Infof("Finished %s %d / %d bytes (%d%%)", resp.Filename, resp.BytesComplete(), resp.Size, int(100*resp.Progress()))
	}
	return nil
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
	if _, err := os.Stat(dirName); err == nil {
		return nil
	}
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

func GetNum(name string) (num string, s scraper.Scraper) {
	typeHeyzo, _ := regexp.Compile(`(heyzo|HEYZO)-[0-9]{4}`)
	typeFc2, _ := regexp.Compile(`(fc2|FC2|ppv|PPV)-[0-9]{6,7}`)
	typeMGStage, _ := regexp.Compile(`(siro|SIRO|[0-9]{3,4}[a-zA-Z]{2,5})-[0-9]{3,4}`)
	typeDmm, _ := regexp.Compile(`[a-zA-Z]{2,5}00[0-9]{3,4}`)
	typeDefault, _ := regexp.Compile(`[a-zA-Z]{2,5}-[0-9]{3,4}`)

	switch {
	case typeHeyzo.MatchString(name):
		num = typeHeyzo.FindString(name)
		s = &scraper.HeyzoScraper{}
	case typeFc2.MatchString(name):
		num = typeFc2.FindString(name)
		s = &scraper.Fc2Scraper{}
	case typeMGStage.MatchString(name):
		num = typeMGStage.FindString(name)
		s = &scraper.MGStageScraper{}
	case typeDmm.MatchString(name):
		num = typeDmm.FindString(name)
		num = strings.Replace(num, "00", "-", 1)
		s = &scraper.DMMScraper{}
	default:
		num = typeDefault.FindString(name)
		s = &scraper.DMMScraper{}
	}

	num = strings.ToUpper(num)
	return
}

func MakeOutputPath(s scraper.Scraper) error {
	if len(s.GetPremiered()) >= 4 {
		outputPath = path.Join(outputPath, s.GetPremiered()[:4])
	}
	log.Infof("Making output path %s", outputPath)
	return ensureDir(outputPath)
}

func BuildNfo(b, num string) error {
	nfoName := path.Join(outputPath, fmt.Sprintf("%s.nfo", num))
	return ioutil.WriteFile(nfoName, []byte(b), 0644)
}

func MoveFile(num, name string) error {
	newPath := strings.ToUpper(num) + filepath.Ext(name)
	return os.Rename(name, path.Join(outputPath, newPath))
}
