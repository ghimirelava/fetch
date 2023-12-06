package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sort"
)

func sqlImageSearch(searchTerm, url string) []TFIDFScore {
	db, err := sql.Open("sqlite3", "project04.db")
	if err != nil {
		fmt.Println("Err in sqlImageSearch() open: ", err)
		os.Exit(-1)
	}
	defer db.Close()

	//transatcion
	tx, err := db.Begin()
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}
	defer tx.Commit()

	var termID int

	err = tx.QueryRow("SELECT image_term_id FROM image_terms WHERE word = ?;", searchTerm).Scan(&termID)
	if err != nil {
		fmt.Printf("error in termID select\n")
		log.Fatal(err)
	}

	var numDocs int
	err = tx.QueryRow("SELECT count(*) FROM image_hits WHERE image_term_id = ?;", termID).Scan(&numDocs)
	if err != nil {
		fmt.Printf("error in numDocs select\n")
		log.Fatal(err)
	}

	var docs_in_C int

	err = tx.QueryRow("SELECT COUNT(image_url_id) FROM image_urls;").Scan(&docs_in_C)
	if err != nil {
		fmt.Printf("error in docs_in_C select\n")
		log.Fatal(err)
	}

	var rows *sql.Rows
	var tfidfScoresSlice []TFIDFScore

	rows, err = tx.Query(` SELECT image_terms.word, image_urls.url_title, image_urls.url, image_hits.src_url, image_hits.alt_text, image_urls.word_count, image_hits.image_term_count
			FROM image_hits
			INNER JOIN image_urls ON image_hits.image_url_id = image_urls.image_url_id
			INNER JOIN image_terms on image_terms.image_term_id = image_hits.image_term_id
			WHERE image_hits.image_term_id = ?;`, termID)

	if err != nil {
		println("Error in non-wc query")
		log.Fatal(err)
	}

	for rows.Next() {
		var word, urlTitle, url, src, alt string
		var wordCount, termCount int
		err := rows.Scan(&word, &urlTitle, &url, &src, &alt, &wordCount, &termCount)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("word: %v, wordCount: %v, termCount: %v, numDocs: %v, docs_in_C: %v\n", word, wordCount, termCount, numDocs, docs_in_C)

		tf_idf := sqlCalcTFIDF(float64(termCount), float64(wordCount), float64(numDocs), float64(docs_in_C))
		tfidfScoresSlice = append(tfidfScoresSlice, TFIDFScore{Word: word, Title: urlTitle, URL: url, Source: src, ALT: alt, Score: tf_idf})
	}

	defer rows.Close()
	sort.Sort(Hits(tfidfScoresSlice))

	return tfidfScoresSlice
}
