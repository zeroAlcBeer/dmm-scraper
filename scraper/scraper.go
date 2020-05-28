package scraper

import (
	"better-av-tool/log"
	"bufio"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding/htmlindex"
	"io"
	"net/http"
	"time"
)

type Scraper interface {
	// Remote
	FetchDoc(query, url string) (err error)

	// Local
	GetPlot() string
	GetTitle() string
	GetDirector() string
	GetRuntime() string
	GetTags() []string
	GetMaker() string
	GetActors() []string
	GetLabel() string
	GetNumber() string
	GetCover() string
	GetWebsite() string
	GetPremiered() string
	GetYear() string
	GetSeries() string

	NeedCut() bool
}

var (
	proxyClient *http.Client
)

func init() {
	if proxyClient == nil {
		proxyClient = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
			},
		}
	}
}

func SetHTTPClient(client *http.Client) {
	proxyClient = client
}

func GetConvertDocFromUrl(u string) (*goquery.Document, error) {
	log.Infof("fetching %s", u)
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	res, err := proxyClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		//if res.StatusCode == 410 {
		//	return nil, errors.New(http.StatusText(410))
		//}
		return nil, errors.New(fmt.Sprintf("status code error: %d %s", res.StatusCode, res.Status))
	}
	utfBody, _ := DecodeHTMLBody(res.Body, "Shift_JIS")
	//data, _ := ioutil.ReadAll(res.Body)
	//_ = ioutil.WriteFile(url.QueryEscape(path.Base(u))+".html", data, 0644)
	return goquery.NewDocumentFromReader(utfBody)
}

func GetDocFromUrl(u string) (*goquery.Document, error) {

	log.Infof("fetching %s", u)
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	c := &http.Cookie{
		Name:    "adc",
		Value:   "1",
		Path:    "/",
		Domain:  "mgstage.com",
		Expires: time.Now().Add(1 * time.Hour),
	}
	req.AddCookie(c)
	res, err := proxyClient.Do(req)
	if err != nil {
		return nil, err
	}
	//defer res.Body.Close()
	if res.StatusCode != 200 {
		//if res.StatusCode == 410 {
		//	return nil, errors.New(http.StatusText(410))
		//}
		return nil, errors.New(fmt.Sprintf("status code error: %d %s", res.StatusCode, res.Status))
	}
	//data, _ := ioutil.ReadAll(res.Body)
	//_ = ioutil.WriteFile(url.QueryEscape(path.Base(u))+".html", data, 0644)
	return goquery.NewDocumentFromResponse(res)
}

func detectContentCharset(body io.Reader) string {
	r := bufio.NewReader(body)
	if data, err := r.Peek(1024); err == nil {
		if _, name, ok := charset.DetermineEncoding(data, ""); ok {
			return name
		}
	}
	return "utf-8"
}

func DecodeHTMLBody(body io.Reader, charset string) (io.Reader, error) {
	if charset == "" {
		charset = detectContentCharset(body)
	}
	e, err := htmlindex.Get(charset)
	if err != nil {
		return nil, err
	}
	if name, _ := htmlindex.Name(e); name != "utf-8" {
		body = e.NewDecoder().Reader(body)
	}
	return body, nil
}
