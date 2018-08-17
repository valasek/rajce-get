package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type remoteFile struct {
	url  string
	file string
}

func main() {
	fmt.Println("starting ...")
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		usage()
	}

	//img31.rajce.idnes.cz/d3103/15/15400/15400849_2b8b5e522f4cde0cd7d89efe79cb0a63/images/MOV_0180.jpg
	urlPattern, err := regexp.Compile(`(img[0-9]{2}.rajce.idnes.cz/.*/)(.*)`)
	if err != nil {
		fmt.Printf("implementation error: regexp can not compile: %s\n", urlPattern.String())
		os.Exit(2)
	}

	var URLs []remoteFile
	doc, err := goquery.NewDocument(flag.Args()[0])
	if err != nil {
		fmt.Printf("there is a problem getting a page %s: %s\n", flag.Args()[0], err)
		os.Exit(2)
	}

	doc.Find("a[href]").Each(func(index int, item *goquery.Selection) {
		href, _ := item.Attr("href")
		res := urlPattern.FindAllStringSubmatch(href, -1)
		for i := range res {
			URLs = append(URLs, remoteFile{url: res[i][1], file: res[i][2]})
		}
	})

	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}
	for _, remFile := range URLs {
		downloaded, err := downloadFile(netClient, remFile.file, "HTTPS://"+remFile.url+"/"+remFile.file)
		if err != nil {
			fmt.Printf("downloading of %s failed due to: %v\n", remFile.file, err)
		} else {
			fmt.Printf("%s dowloaded, %d bytes\n", remFile.file, downloaded)
		}
	}

	fmt.Printf("finished, downloaded %d items", len(URLs))
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: rajce-get <url>\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func downloadFile(netClient *http.Client, filePath string, url string) (int64, error) {

	out, err := os.Create(filePath)

	if err != nil {
		return 0, err
	}
	defer out.Close()

	resp, err := netClient.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("get request to %s returned HTTP %d, expected HTTP %d", url, resp.StatusCode, http.StatusOK)
	}

	return io.Copy(out, resp.Body)
}
