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
	http.Handle("/form.css", http.FileServer(http.Dir("./static")))

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

		} else if image != "" { //image search
			stemmed, err := snowball.Stem(searchTerm, "english", true)
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}

			tfidfScores = sqlImageSearch(stemmed, url)

			t, err := template.ParseFiles("./static/imageResult.html")
			if err != nil {
				log.Fatalln("ParseFiles: ", err)
			}

			w.Header().Set("Content-Type", "text/html")
			err = t.Execute(w, tfidfScores)

		} else if len(splitTerm) > 1 { //bigram search
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

			t, err := template.ParseFiles("./static/uniResult.html")
			if err != nil {
				log.Fatalln("ParseFiles: ", err)
			}

			w.Header().Set("Content-Type", "text/html")
			err = t.Execute(w, tfidfScores)

		} else { //uni search
			stemmed, err := snowball.Stem(searchTerm, "english", true)
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
			tfidfScores = sqlSearch(stemmed, url, wildcard)

			t, err := template.ParseFiles("./static/uniResult.html")
			if err != nil {
				log.Fatalln("ParseFiles: ", err)
			}

			w.Header().Set("Content-Type", "text/html")
			err = t.Execute(w, tfidfScores)
		}
	})

	return tfidfScores

}
