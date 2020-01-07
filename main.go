package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/proxy"
	"log"
	"net"
	"net/http"
	"strings"
)

var proxyClient *http.Client

func init() {
	proxyClient = &http.Client{}
	dialer, _ := proxy.SOCKS5("tcp", "127.0.0.1:7891", nil, proxy.Direct)
	proxyClient.Transport = &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
			c, e := dialer.Dial(network, addr)
			return c, e
		},
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
}

func GetDocFromUrl(url string) (*goquery.Document, error) {
	res, err := proxyClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		err = errors.New(fmt.Sprintf("status code error: %d %s", res.StatusCode, res.Status))
		return nil, err
	}

	return goquery.NewDocumentFromReader(res.Body)
}

func DmmScrape(searchUrl string) (err error) {
	doc, err := GetDocFromUrl(searchUrl)
	if err != nil {
		return
	}

	// Find the ul id=list items
	itemCount := doc.Find("#list li").Length()
	if itemCount == 0 {
		err = errors.New("record not found")
		return
	}
	if itemCount > 2 {
		err = errors.New("multi records, make number specific")
		return
	}

	detailUrl, exist := doc.Find("#list li").First().Find(".tmb a").Attr("href")

	if !exist {
		err = errors.New("detail link not found")
		return
	}

	doc, err = GetDocFromUrl(detailUrl)
	if err != nil {
		return
	}

	fmt.Println(getTitle(doc))
	fmt.Println(getTableValue("メーカー", doc))
	fmt.Println(getTableValue("発売日", doc))
	fmt.Println(getDesc(doc))
	fmt.Println(getActors(doc))
	fmt.Println(getTags(doc))
	fmt.Println(getCover(doc))

	fmt.Println(getTableValue("シリーズ", doc))
	fmt.Println(getTableValue("収録時間", doc))
	fmt.Println(getTableValue("監督", doc))
	fmt.Println(getTableValue("レーベル", doc))

	//dic = {
	//	'title': str(re.sub('\w+-\d+-', '', getTitle(htmlcode))),
	//	'studio': getStudio(htmlcode),
	//		'year': str(re.search('\d{4}', getYear(htmlcode)).group()),
	//	'outline': getOutline(dww_htmlcode),
	//		'runtime': getRuntime(htmlcode),
	//		'director': getDirector(htmlcode),
	//		'actor': getActor(htmlcode),
	//		'release': getRelease(htmlcode),
	//		'number': getNum(htmlcode),
	//		'cover': getCover(htmlcode),
	//		'imagecut': 1,
	//		'tag': getTag(htmlcode),
	//		'label': getSerise(htmlcode),
	//		'actor_photo': getActorPhoto(htmlcode),
	//		'website': 'https://www.javbus.com/' + number,
	//	'source' : 'javbus.py',
	//}

	return
}

func getActors(doc *goquery.Document) (actors []string) {
	doc.Find("#performer").Each(func(i int, s *goquery.Selection) {
		actors = append(actors, s.Find("a").Text())
	})
	return
}

func getTags(doc *goquery.Document) (tags []string) {
	doc.Find("table[class=mg-b20] tr").EachWithBreak(
		func(i int, s *goquery.Selection) bool {
			if strings.Contains(s.Text(), "ジャンル") {
				s.Find("td a").Each(func(i int, ss *goquery.Selection) {
					tags = append(tags, ss.Text())
				})
				return false
			}
			return true
		})
	return
}

func getTitle(doc *goquery.Document) string {
	return doc.Find("#title").Text()
}

func getDesc(doc *goquery.Document) string {
	return doc.Find("p[class=mg-b20]").Text()
}

func getTableValue(key string, doc *goquery.Document) (val string) {
	doc.Find("table[class=mg-b20] tr").EachWithBreak(
		func(i int, s *goquery.Selection) bool {
			if strings.Contains(s.Text(), key) {
				val = s.Find("td a").Text()
				if val == "" {
					val = s.Find("td").Last().Text()
				}
				return false
			}
			return true
		})
	return
}

func getCover(doc *goquery.Document) string {
	url, _ := doc.Find("#sample-video a").First().Attr("href")
	return url
}

func main() {
	err := DmmScrape("https://www.dmm.co.jp/mono/dvd/-/search/=/searchstr=miaa201/list_type=mono/")
	if err != nil {
		log.Fatal(err)
	}

}
