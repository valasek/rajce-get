package main

import (
	"io"
	"flag"
	"fmt"
	"net/http"
	"github.com/PuerkitoBio/goquery"
	"os"
	"time"
	"regexp"
)

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
	type remoteFile struct {
		url string
		file string
	}
	var URLs []remoteFile
	doc, err := goquery.NewDocument(flag.Args()[0])
	if err != nil {
        fmt.Println("There is a problem getting a page.", err)
		os.Exit(2)
	}
	downloaded := 0
	doc.Find("a[href]").Each(func(index int, item *goquery.Selection) {
		href, _ := item.Attr("href")
		res := urlPattern.FindAllStringSubmatch(href, -1)
		for i := range res {
			//fmt.Println(res[i][1], res[i][2], "\n")
			URLs = append(URLs, remoteFile{url: res[i][1], file: res[i][2]})
			downloaded++
		}
	})
	for _, remFile := range URLs {
		downloaded, err := downloadFile(remFile.file, "HTTPS://"+remFile.url+"/"+remFile.file)
		if err != nil {
			fmt.Printf("Downloading failed (%s). Error: %v", remFile.file, remFile.url)
		}
		time.Sleep(1 * time.Second)
		fmt.Printf("Dowloaded %d bytes (%s)\n", downloaded, remFile.file)
	}

	fmt.Printf("Finished, downloaded %d items", downloaded)
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: rajce-get <url>\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func downloadFile(filePath string, url string) (downloaded int64, error error) {
	out, err := os.Create(filePath)
	if err != nil {
     	return 0, err
	}
	defer out.Close()
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	n, err := io.Copy(out, resp.Body)
	if err != nil {
		return 0, err
	}
	return n, nil
}
