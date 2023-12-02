package main

import (
	"net/http"
	"reflect"
	"testing"
)

func BigramTest(t *testing.T) {
	tests := []struct {
		term string
		want Hits
	}{
		{"computer science",
			Hits{
				{"", "A–Z index – UC Santa Cruz", "https://www.ucsc.edu/azindex/", 0.018746338605741066},
			},
		},
	}
	makeQueries()
	go http.ListenAndServe("localhost:8080", nil)

	url := "https://www.ucsc.edu/robots.txt"
	tfidfScores := sqlHandleSearch(url)

	for _, tn := range tests {

		for _, s := range tn.want {
			for _, r := range tfidfScores {
				if r.URL == s.URL {
					if !reflect.DeepEqual(r.Score, s.Score) {
						t.Errorf("Got vs want not the same for %v. \n", tn.term)
						break
					}
				}
			}
		}

	}

}

func WildcardTest(t *testing.T) {
	tests := []struct {
		term string
		want Hits
	}{
		{"scho",
			Hits{
				{"scholar", "Mission and Vision – UC Santa Cruz", "https://www.ucsc.edu/mission-and-vision/", 0.11678832116788321},
				{"school", "UC Santa Cruz – A world-class public research institution comprised of ten residential college communities nestled in the redwood forests and meadows overlooking central California's Monterey Bay.", "https://www.ucsc.edu/ ", 0.032160804020100506},
				{"scholar", "Achievements, Facts, and Figures – UC Santa Cruz", "https://www.ucsc.edu/about/achievements-facts-and-figures/", 0.031746031746031744},
				{"scholarship", "Mission and Vision – UC Santa Cruz", "https://www.ucsc.edu/mission-and-vision/", 0.029197080291970802},
				{"scholarship", "Admissions & Aid – UC Santa Cruz", "https://www.ucsc.edu/admissions/", 0.028368794326241134},
				{"school", "UC + Santa Cruz. Better together – UC Santa Cruz", "https://www.ucsc.edu/better-together/", 0.016623376623376623},
				{"school", "Programs and Units – UC Santa Cruz", "https://www.ucsc.edu/programs-and-units/", 0.014814814814814815},
				{"school", "Academics – UC Santa Cruz", "https://www.ucsc.edu/academics/", 0.01306122448979592},
				{"scholarship", "Achievements, Facts, and Figures – UC Santa Cruz", "https://www.ucsc.edu/about/achievements-facts-and-figures/", 0.011904761904761904},
				{"scholarship", "Overview – UC Santa Cruz", "https://www.ucsc.edu/about/overview/", 0.011594202898550725},
				{"scholarship", "About UC Santa Cruz – UC Santa Cruz", "https://www.ucsc.edu/about/", 0.011235955056179775},
				{"scholarship", "UC + Santa Cruz. Better together – UC Santa Cruz", "https://www.ucsc.edu/better-together/", 0.01038961038961039},
				{"scholarship", "UC Santa Cruz – A world-class public research institution comprised of ten residential college communities nestled in the redwood forests and meadows overlooking central California's Monterey Bay.", "https://www.ucsc.edu/", 0.010050251256281407},
				{"school", "A–Z index – UC Santa Cruz", "https://www.ucsc.edu/azindex/", 0.007498535442296427},
				{"scholar", "A–Z index – UC Santa Cruz", "https://www.ucsc.edu/azindex/", 0.0062487795352470215},
				{"scholarship", "A–Z index – UC Santa Cruz", "https://www.ucsc.edu/azindex/", 0.0023432923257176333},
			},
		},
	}

	for _, tn := range tests {

		for _, s := range tn.want {
			tfidfScores := sqlSearch(tn.term, s.URL, "wildcard")
			for _, r := range tfidfScores {
				if r.URL == s.URL {
					if !reflect.DeepEqual(r.Score, s.Score) {
						t.Errorf("Got vs want not the same for %v. \n", tn.term)
						break
					}
				}
			}
		}

	}

}
