package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/kljensen/snowball"
)

func sqlHandleSearch(url string) []TFIDFScore {
	//var d Datasql
	var tfidfScores []TFIDFScore

	db, err := sql.Open("sqlite3", "project04.db")
	if err != nil {
		fmt.Println("Err in popTermsTable() open: ", err)
		os.Exit(-1)
	}
	defer db.Close()

	//transatcion
	tx, err := db.Begin()
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}
	//calling crawl with url to init the inverted index
	//d.sqlCrawl(url)

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.Handle("/project06.css", http.FileServer(http.Dir("./static")))

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {

		//obatin the search term and wlidcard from request
		searchTerm := r.URL.Query().Get("term")
		wildcard := r.URL.Query().Get("wildcard")
		image := r.URL.Query().Get("image")

		//split the search term
		splitTerm := strings.Split(searchTerm, " ")

		// if no search term is obtained then pront error
		if searchTerm == "" {
			log.Fatalln("No search term in handleSearch()")
		} else if image != "" {
			stemmed, err := snowball.Stem(searchTerm, "english", true)
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
			tfidfScores = sqlImageSearch(stemmed, url)
		} else if len(splitTerm) > 1 { //if it is a bigram stem and search in bigram specfifc function
			stemmedOne, err := snowball.Stem(splitTerm[0], "english", true)
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
			stemmedTwo, err := snowball.Stem(splitTerm[1], "english", true)
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
			tfidfScores = sqlSearchBiGram(stemmedOne, stemmedTwo, url)
		} else {
			stemmed, err := snowball.Stem(searchTerm, "english", true)
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
			tfidfScores = sqlSearch(stemmed, url, wildcard)
		}
		t, err := template.ParseFiles("./static/result.html")
		if err != nil {
			log.Fatalln("ParseFiles: ", err)
		}

		w.Header().Set("Content-Type", "text/html")
		err = t.Execute(w, tfidfScores)

		/*for _, val := range tfidfScores {
			if image != "" {
				fmt.Fprintf(w, "%s\n %s\n %s : %v\n\n", val.Word, val.Title, val.URL, val.Score)
				fmt.Fprintf(w, "%v\n %v\n\n", val.Source, val.ALT)
			} else {
				fmt.Fprintf(w, "%s\n %s\n %s : %v\n\n", val.Word, val.Title, val.URL, val.Score)
			}
		}*/
	})

	return tfidfScores

}
