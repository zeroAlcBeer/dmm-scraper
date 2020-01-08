package main

import (
	"better-av-tool/log"
	"better-av-tool/nfo"
	"better-av-tool/scraper"
	"context"
	"crypto/tls"
	"flag"
	"golang.org/x/net/proxy"
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
	output   string
	proxyUrl string

	proxyClient *http.Client
)

func init() {
	flag.StringVar(&output, "o", "", "set output path")
	flag.StringVar(&proxyUrl, "p", "", "set proxy url")
	flag.Parse()

}

func main() {
	var err error

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
	}

	if output == "" {
		if _, err := os.Stat("OUTPUT"); os.IsNotExist(err) {
			err = os.Mkdir("OUTPUT", 0644)
			if err != nil {
				log.Fatal(err)
			}
		}
	} else {
		if _, err := os.Stat(output); os.IsNotExist(err) {
			err = os.Mkdir(output, 0644)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	myDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	var ff = func(pathX string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// find out if it's a dir or file, if file, print info
		if !info.IsDir() {
			r, _ := regexp.Compile(`[a-zA-Z]{2,5}-[0-9]{3,4}`)
			num := r.FindString(info.Name())
			if num == "" {
				r, _ := regexp.Compile(`[0-9]{3,4}[a-zA-Z]{2,5}-[0-9]{3,4}`)
				num = r.FindString(info.Name())
			}
			if num != "" {
				num = strings.ToUpper(num)
				log.Infof("  num %s found!", num)

				s := &scraper.MGStageScraper{}
				s.HTTPClient = proxyClient
				err = s.FetchDoc(num)
				if err != nil {
					log.Error(err)
					return nil
				}
				b, err := nfo.Build(s)
				if err != nil {
					log.Error(err)
					return nil
				}
				err = ioutil.WriteFile(path.Join(output, num+".nfo"), b, 0644)
				if err != nil {
					log.Error(err)
					return nil
				}
				//err = os.Rename(info.Name(), path.Join(output, info.Name()))
				//if err != nil {
				//	return err
				//}
			}
		}

		return nil
	}

	_ = filepath.Walk(myDir, ff)

}
