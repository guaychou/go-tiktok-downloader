package main

import (
	"flag"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"encoding/json"
	"os"
)


type Data struct {
	Name string `json:"name"`
	ContentUrl string `json:"contentUrl"`
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
	download(scraper,data)
}

func newScraper() *colly.Collector {
	c := colly.NewCollector()
	c.UserAgent="Mozilla/5.0 (X11; Ubuntu; Linux i686; rv:28.0) Gecko/20100101 Firefox/28.0"
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept","text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
		r.Headers.Set("Accept-Encoding","gzip, deflate")
		r.Headers.Set("Accept-Language","en-US,en;q=0.9")
	})
	return c
}

func getVideoLink(copiedLink string,scraper *colly.Collector)(Data,error){
	data := Data{}
	scraper.OnHTML("script[id=videoObject]", func(e *colly.HTMLElement) {
		err:=json.Unmarshal([]byte(e.Text),&data)
		checkErr(err)
	})
	err:=scraper.Visit(copiedLink)
	checkErr(err)
	return data,nil
}

func checkErr(err error) {
	if err!=nil{
		log.Fatal(err)
	}
}

func download(scraper *colly.Collector, data Data){
	scraper.OnResponse(func(response *colly.Response) {
		video, err:=os.Create(data.Name+".mp4")
		if err!=nil{
			log.Fatal(err)
		}
		defer video.Close()
		n2,err:=video.Write(response.Body)
		fmt.Printf("Wrote %d bytes\n", n2)
		video.Sync()
	})
	err:=scraper.Visit(data.ContentUrl)
	checkErr(err)
}