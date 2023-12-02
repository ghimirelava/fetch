package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
)

type (
	Sitemap struct {
		Loc string `xml:"loc"`
	}

	SitemapIndex struct {
		Sitemaps []Sitemap `xml:"sitemap"`
	}

	UrlIndex struct {
		Urls []Sitemap `xml:"url"`
	}
)

func downloadSiteMap(url string) ([]byte, error) {
	if rsp, err := http.Get(url); err == nil {
		if b, err := io.ReadAll(rsp.Body); err == nil {
			// the scope of rsp, b, err are ONLY inside the "if" clause, not outside
			return b, nil
		} else {
			println("error in i0.ReadAll")
			log.Fatal(err)
		}
	} else {
		println("error in http.Get in downloadSiteMap")
		log.Fatal(err)
	}

	return []byte{}, nil
}

func retrieveSiteMapLinks(url string, dInCh chan string) {
	data, _ := downloadSiteMap(url)
	v := &SitemapIndex{}

	err := xml.Unmarshal([]byte(data), &v)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	for _, val := range v.Sitemaps {

		data, _ := downloadSiteMap(val.Loc)

		h := &UrlIndex{}

		err := xml.Unmarshal([]byte(data), &h)
		if err != nil {
			fmt.Printf("error: %h\n", err)
			return
		}
		for _, val := range h.Urls {
			dInCh <- val.Loc
		}

	}

}
