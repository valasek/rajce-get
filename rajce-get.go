package main

import (
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"os"
	"regexp"
	"time"
)

type remoteFile struct {
	url  string
	file string
}

func main() {
	fmt.Println("Starting ...")
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		usage()
	}

	//img31.rajce.idnes.cz/d3103/15/15400/15400849_2b8b5e522f4cde0cd7d89efe79cb0a63/images/MOV_0180.jpg
	urlPattern, err := regexp.Compile(`(img[0-9]{2}.rajce.idnes.cz/.*/)(.*)`)
	if err != nil {
		fmt.Printf("There is a problem with your regexp.\n")
		os.Exit(2)
	}

	var URLs []remoteFile
	doc, err := goquery.NewDocument(flag.Args()[0])
	if err != nil {
		fmt.Println("There is a problem getting a page.", err)
		os.Exit(2)
	}

	doc.Find("a[href]").Each(func(index int, item *goquery.Selection) {
		href, _ := item.Attr("href")
		res := urlPattern.FindAllStringSubmatch(href, -1)
		for i := range res {
			URLs = append(URLs, remoteFile{url: res[i][1], file: res[i][2]})
		}
	})

	for _, remFile := range URLs {
		downloaded, err := downloadFile(remFile.file, "HTTPS://"+remFile.url+"/"+remFile.file)
		if err != nil {
			fmt.Printf("Downloading failed (%s). Error: %v", remFile.file, remFile.url)
		}
		fmt.Printf("Dowloaded %d bytes (%s)\n", downloaded, remFile.file)
	}

	fmt.Printf("Finished, downloaded %d items", len(URLs))
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: rajce-get <url>\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func downloadFile(filePath string, url string) (int64, error) {

	out, err := os.Create(filePath)

	if err != nil {
		return 0, err
	}
	defer out.Close()

	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := netClient.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	return io.Copy(out, resp.Body)
}
