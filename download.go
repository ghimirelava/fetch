package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
)

type DownloadResult struct {
	body []byte
	url  string
}

func (rd *robotData) checkDisallow(url string) bool {
	shouldDownload := true

	for val, _ := range rd.disallow {
		regex, err := regexp.Compile(val)
		if err != nil {
			log.Fatal("Error compiling regular expression:", err)
		}
		if regex.MatchString(url) {
			// URL matches the disallow rule, do not download
			shouldDownload = false
			break
		}
	}
	return shouldDownload
}

func download(inC chan string, outC chan DownloadResult, rd *robotData) {
	for url := range inC {
		shouldDownload := rd.checkDisallow(url)

		if shouldDownload == true {
			if rsp, err := http.Get(url); err == nil {
				//defer rsp.Body.Close()

				if rsp.StatusCode == 200 {
					if b, err := io.ReadAll(rsp.Body); err == nil {
						outC <- DownloadResult{b, url}
					} else {
						fmt.Println("Error reading response body:", err)
					}
				}
			} else {
				fmt.Println(err)
				fmt.Println("Error in http.Get:", err)
			}
		}
	}
}
