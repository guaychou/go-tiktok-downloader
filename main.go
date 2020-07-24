package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"os"

	"os"
)

type VideoProps struct {
	Url []string `json:"urls"`
}

type Data struct {
	VideoProps VideoProps `json:"video"`
	ImagePreview []string `json:"covers"`
	Text string `json:"text"`
}
type ItemInfos struct {
	ItemInfos Data `json:"itemInfos"`
}

type VideoData struct {
	VideoData ItemInfos `json:"videoData"`
}

type PageProps struct {
	PageProps VideoData `json:"pageProps"`
}

type Props struct {
	Props PageProps `json:"props"`
}

func main(){
	copiedLink:=flag.String("link","","Tiktok Copied link \n Example: https://vt.tiktok.com/D8RK6S/")
	flag.Parse()
	if *copiedLink==""{
		flag.PrintDefaults()
		os.Exit(0)
	}
	scraper :=newScraper()
	data,err:=getVideoLink(*copiedLink,scraper)
	checkErr(err)
	log.Println("URL :"+data.Props.PageProps.VideoData.ItemInfos.VideoProps.Url[0])
	download(scraper,data)
}

func newScraper() *colly.Collector {
	c := colly.NewCollector()
	c.UserAgent="Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.89 Safari/537.36"
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept","text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
		r.Headers.Set("Accept-Encoding","gzip, deflate")
		r.Headers.Set("Accept-Language","en-US,en;q=0.9")
	})
	return c
}

func getVideoLink(copiedLink string,scraper *colly.Collector)(Props,error){
	data:=Props{}
	var result string
	scraper.OnHTML("script[id=__NEXT_DATA__]", func(e *colly.HTMLElement) {
		err:=json.Unmarshal([]byte(e.Text),&data)
		checkErr(err)
	})
	err:=scraper.Visit(copiedLink)
	fmt.Print(result)
	checkErr(err)
	return data,nil
}

func checkErr(err error) {
	if err!=nil{
		log.Fatal(err)
	}
}

func download(scraper *colly.Collector, data Props){
	scraper.OnResponse(func(response *colly.Response) {
		video, err:=os.Create(data.Props.PageProps.VideoData.ItemInfos.Text+".mp4")
		checkErr(err)
		defer video.Close()
		n2,err:=video.Write(response.Body)
		fmt.Printf("Wrote %d bytes\n", n2)
		video.Sync()
	})
	err:=scraper.Visit(data.Props.PageProps.VideoData.ItemInfos.VideoProps.Url[0])
	checkErr(err)
}