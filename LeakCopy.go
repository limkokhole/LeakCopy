package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

const LeakBaseMainURL string = "http://siph0n.in/leaks.php?page="

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	var wg sync.WaitGroup
	var addrList []string

	wg.Add(0x14)

	for i := 0x1; i != 0x15; i++ {
		url := LeakBaseMainURL + strconv.Itoa(i)
		go func(funcUrl string, gw *sync.WaitGroup) {
			defer gw.Done()
			doc, err := goquery.NewDocument(funcUrl)
			if err != nil {
				log.Fatal(err.Error())
			}

			doc.Find("td").Each(func(i int, s *goquery.Selection) {
				linkTag := s.Find("a")
				link, _ := linkTag.Attr("href")
				if strings.Contains(string(link), "download") {
					addrList = append(addrList, link)
				}
			})

		}(url, &wg)
	}

	wg.Wait()

	currentPath, _ := os.Getwd()
	savePath := currentPath + "/Leaks/"
	os.Mkdir(savePath, 0700)

	wg.Add(len(addrList))
	for i, uhm := range addrList {
		fmt.Printf("%d - http://siph0n.in/%s\t\t---> Downloading <---\n", i, uhm)
		go func(url string, gw *sync.WaitGroup) {
			defer gw.Done()

			resp, err := http.Get(string("http://siph0n.in/" + url))
			if err != nil {
				fmt.Printf("Err: %s\n", err.Error())
			} else {
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				err = ioutil.WriteFile(savePath+url+".txt", body, 0700)
				if err != nil {
					fmt.Printf("Err: %s\n", err.Error())
				} else {
					fmt.Printf("%s Downloaded!\t\t\n", url)
				}
			}

		}(uhm, &wg)
	}

	wg.Wait()

	fmt.Println("Success!")
}
