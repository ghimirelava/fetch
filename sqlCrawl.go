package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/kljensen/snowball"
)

type Datasql struct {
	prevDocsInCropus float64
	crawlMap         map[string]bool
	biGramMap        map[string]int //term and term_id
	snippetMap       map[string]int
}

func (d *Datasql) sqlPopulateIndex(seedURL string, eOutCh chan ExtractResult, dInCh chan string) {

	var urlID int
	for er := range eOutCh { // hrefs and words

		d.prevDocsInCropus++

		d.biGramMap = make(map[string]int)
		d.snippetMap = make(map[string]int)
		cleanedWords := []string{}

		stopWordsMap := stopWords() //get the map of stop words
		for _, word := range er.wordSlice {
			stemmed, err := snowball.Stem(word, "english", true)
			if err == nil {
				if _, ok := stopWordsMap[stemmed]; !ok {
					termID, url_id := popTables(stemmed, er.url, er.title) //return term_id to put into map

					urlID = url_id
					d.biGramMap[stemmed] = termID
					cleanedWords = append(cleanedWords, stemmed)
				}
			} else {
				log.Fatalln("Does not stem properly!")
			}
		}
		//image tables
		for _, word := range er.altSlice {
			stemmed, err := snowball.Stem(word, "english", true)
			if err == nil {
				if _, ok := stopWordsMap[stemmed]; !ok {
					for key, val := range er.imgInfoMap {
						if strings.Contains(val, word) {
							popImageTable(stemmed, er.url, er.title, key, val)
						}
					}
				}
			} else {
				log.Fatalln("Does not stem properly!")
			}
		}
		//populate the bigrams table by looping through the array of cleaned urls
		//and getting the term id from the map created for each url
		for i := 0; i < (len(cleanedWords) - 1); i++ {
			term_id_one := d.biGramMap[cleanedWords[i]]
			term_id_two := d.biGramMap[cleanedWords[i+1]]
			popBiGram(term_id_one, term_id_two, urlID)
		}
	}
}

func (d *Datasql) sqlCrawl(seedURL string) error {

	d.crawlMap = make(map[string]bool) //making map that tracks urls that have been crawled
	d.crawlMap[seedURL] = true         //initialize map with seed url

	var rd robotData
	rd.userAgents = []string{}
	rd.disallow = make(map[string]bool)

	rd.downloadRobots(seedURL)
	rd.extractRobots()

	dInCh := make(chan string, 150)
	defer close(dInCh)
	dOutCh := make(chan DownloadResult, 160)
	defer close(dOutCh)
	eOutCh := make(chan ExtractResult, 20)
	defer close(eOutCh)
	quitC := make(chan struct{}, 1) //changed to struct for idomatic go

	retrieveSiteMapLinks(rd.siteMapLink, dInCh)

	go download(dInCh, dOutCh, &rd)               //download body
	go extract(dOutCh, eOutCh)                    //extract wordSlice and hrefSlice form body
	go d.sqlPopulateIndex(seedURL, eOutCh, dInCh) //words and hrefs

	go func() {
		for {
			time.Sleep(10 * time.Second)
			waitForCh := false

			if len(dInCh) == 150 {
				waitForCh = true
			}

			if len(dOutCh) == 160 {
				waitForCh = true
			}

			if len(eOutCh) == 20 {
				waitForCh = true
			}

			if d.prevDocsInCropus == float64(len(d.crawlMap)) && !waitForCh {
				fmt.Println("Crawl finished.")
				quitC <- struct{}{}
				break
			} else {
				d.prevDocsInCropus = float64(len(d.crawlMap))
			}
		}

	}()

outer:
	for {
		select {
		case <-quitC:
			break outer
		}
	}
	return nil

}
