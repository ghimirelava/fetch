package main

import (
	"log"
	"net/url"
	"strings"
)

// clean the urls found within the extract function
func clean(host string, hrefs []string) (results []string) {
	results = []string{}

	//parsing the host url
	hostURL, err := url.Parse(host)

	if err != nil {
		log.Fatal(err) //error parsing host url
	}

	for _, val := range hrefs {

		// parsing each value from the hrefs slice
		valURL, err := url.Parse(val)

		if err != nil {
			log.Fatal(err) //error parsing val from hrefs
		}

		if !(strings.Contains(val, ".jpg") || strings.Contains(val, ".png")) && !strings.Contains(val, "#") {
			// if scheme is missing fill it in with the host scheme
			if valURL.Scheme == "" {
				valURL.Scheme = hostURL.Scheme
			}

			// if host is missing fill it in with the host url
			if valURL.Host == "" {
				valURL.Host = hostURL.Host
			}

			// if path is empty add "/"
			if valURL.Path == "" {
				valURL.Path = "/"
			}

			if valURL.Host == hostURL.Host {

				// ignoring non http schemes to make sure we dont get mailto or javascript schemes
				if valURL.Scheme == "https" || valURL.Scheme == "http" {
					// append valURL (changed or not changed) to results slice
					results = append(results, valURL.String())
				}
			}

			if len(valURL.Path) == 0 && len(valURL.Fragment) != 0 {
				// Fragment only (e.g. "#foo")
				// take the host path and clear the Fragment
				// The main URL should be crawled or de-duped by the crawler
				valURL.Path = hostURL.Path
				valURL.Fragment = ""
			}
		}

	}
	return results
}
