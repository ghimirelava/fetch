package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type robotData struct {
	body        []byte
	userAgents  []string
	disallow    map[string]bool
	crawlDelay  int
	siteMapLink string
}

func (rd *robotData) extractRobots() {

	stringBody := string(rd.body)

	ua := regexp.MustCompile(`(User-agent:)`)
	dis := regexp.MustCompile(`(Disallow:)`)
	cd := regexp.MustCompile(`(Crawl-delay:)`)
	sm := regexp.MustCompile(`(Sitemap:)`)

	line := strings.Split(string(stringBody), "\n")

	for _, val := range line {
		if ua.MatchString(string(val)) { // check for uer agents
			subString := strings.TrimPrefix(val, "User-agent: ")
			rd.userAgents = append(rd.userAgents, subString)

		} else if dis.MatchString(string(val)) { // check for disallow rules
			subString := strings.TrimPrefix(val, "Disallow: ")
			string := "." + subString
			rd.disallow[string] = true

		} else if sm.MatchString(string(val)) { // check for disallow rules
			subString := strings.TrimPrefix(val, "Sitemap: ")
			rd.siteMapLink = subString

		} else if cd.MatchString(string(val)) { //check crawl delay
			subString := strings.TrimPrefix(val, "Crawl-delay: ")
			s := strings.TrimSpace(subString)

			num, err := strconv.Atoi(s)
			if err != nil {
				fmt.Println("Error during conversion")
				return
			}
			rd.crawlDelay = num
		}
	}
}

func (rd *robotData) downloadRobots(robotURL string) {
	if rsp, err := http.Get(robotURL); err == nil {
		if b, err := io.ReadAll(rsp.Body); err == nil {
			// the scope of rsp, b, err are ONLY inside the "if" clause, not outside
			rd.extractRobots()
			rd.body = b
		} else {
			println("error in i0.ReadAll (dr)")
			log.Fatal(err)
		}
	} else {
		println("error in http.Get in dr")
		log.Fatal(err)
	}
}
