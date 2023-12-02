package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sort"
)

func sqlSearchBiGram(searchTermOne, searchTermTwo, url string) []TFIDFScore {
	db, err := sql.Open("sqlite3", "project04.db")
	if err != nil {
		fmt.Println("Err in popSearchBiGram() open: ", err)
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

	var rows *sql.Rows
	var tfidfScoresSlice []TFIDFScore
	var term_id_one, term_id_two int

	var docs_in_C int
	err = tx.QueryRow("SELECT COUNT(url_id) FROM urls;").Scan(&docs_in_C)
	if err != nil {
		fmt.Printf("error in docs_in_C select\n")
		log.Fatal(err)
	}

	err = tx.QueryRow("SELECT term_id FROM terms WHERE word = ?;", searchTermOne).Scan(&term_id_one)
	if err != nil {
		fmt.Printf("error in term_id_one select\n")
		log.Fatal(err)
	}

	err = tx.QueryRow("SELECT term_id FROM terms WHERE word = ?;", searchTermTwo).Scan(&term_id_two)
	if err != nil {
		fmt.Printf("error in term_id_two select\n")
		log.Fatal(err)
	}

	var numDocs int
	err = tx.QueryRow("SELECT count(*) FROM bigram_hits WHERE term_id_one = ? AND term_id_two = ?;", term_id_one, term_id_two).Scan(&numDocs)
	if err != nil {
		fmt.Printf("error in numDocs select\n")
		log.Fatal(err)
	}

	rows, err = tx.Query(` SELECT urls.url_title, urls.url, urls.word_count, bigram_hits.term_count
			FROM urls 
			INNER JOIN bigram_hits
			WHERE urls.url_id = bigram_hits.url_id AND bigram_hits.term_id_one = ? AND bigram_hits.term_id_two = ?;`, term_id_one, term_id_two)

	if err != nil {
		println("Error in bi query\n")
		log.Fatal(err)
	}

	for rows.Next() {
		var urlTitle, url string
		var wordCount, termCount int
		err := rows.Scan(&urlTitle, &url, &wordCount, &termCount)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("wordCount: %v, termCount: %v, numDocs: %v, docs_in_C: %v\n", wordCount, termCount, numDocs, docs_in_C)

		tf_idf := sqlCalcTFIDF(float64(termCount), float64(wordCount), float64(numDocs), float64(docs_in_C))
		tfidfScoresSlice = append(tfidfScoresSlice, TFIDFScore{Title: urlTitle, URL: url, Score: tf_idf})
	}

	defer rows.Close()
	sort.Sort(Hits(tfidfScoresSlice))

	return tfidfScoresSlice

}

func sqlSearch(searchTerm string, url string, wildcard string) []TFIDFScore {
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
	defer tx.Commit()

	println(searchTerm)
	println(wildcard)

	var termID int

	if wildcard == "" {
		err = tx.QueryRow("SELECT term_id FROM terms WHERE word = ?;", searchTerm).Scan(&termID)
		if err != nil {
			fmt.Printf("error in termID select\n")
			log.Fatal(err)
		}
	}

	var numDocs int
	err = tx.QueryRow("SELECT count(*) FROM hits WHERE term_id = ?;", termID).Scan(&numDocs)
	if err != nil {
		fmt.Printf("error in numDocs select\n")
		log.Fatal(err)
	}

	var docs_in_C int
	err = tx.QueryRow("SELECT COUNT(url_id) FROM urls;").Scan(&docs_in_C)
	if err != nil {
		fmt.Printf("error in docs_in_C select\n")
		log.Fatal(err)
	}

	var rows *sql.Rows
	var tfidfScoresSlice []TFIDFScore

	if wildcard == "wildcard" {

		rows, err = tx.Query(` SELECT terms.word, urls.url_title, urls.url, urls.word_count, hits.term_count, hits.term_id
			FROM terms
			INNER JOIN urls ON hits.url_id = urls.url_id
			INNER JOIN hits ON terms.term_id = hits.term_id
			WHERE terms.word LIKE ?;`, (searchTerm + "%"))
		if err != nil {
			println("Error in wc query")
			log.Fatal(err)
		}

		for rows.Next() {
			var word, urlTitle, url string
			var wordCount, termCount, term_id int
			err := rows.Scan(&word, &urlTitle, &url, &wordCount, &termCount, &term_id)
			if err != nil {
				log.Fatal(err)
			}
			err = tx.QueryRow("SELECT count(*) FROM hits WHERE term_id = ?;", term_id).Scan(&numDocs)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("word: %v, wordCount: %v, termCount: %v, numDocs: %v, docs_in_C: %v\n", word, wordCount, termCount, numDocs, docs_in_C)

			tf_idf := sqlCalcTFIDF(float64(termCount), float64(wordCount), float64(numDocs), float64(docs_in_C))
			tfidfScoresSlice = append(tfidfScoresSlice, TFIDFScore{Word: word, Title: urlTitle, URL: url, Score: tf_idf})
		}

	} else {
		rows, err = tx.Query(` SELECT terms.word, urls.url_title, urls.url, urls.word_count, hits.term_count
			FROM hits
			INNER JOIN urls ON hits.url_id = urls.url_id
			INNER JOIN  terms on terms.term_id = hits.term_id
			WHERE hits.term_id = ?;`, termID)

		if err != nil {
			println("Error in non-wc query")
			log.Fatal(err)
		}

		for rows.Next() {
			var word, urlTitle, url string
			var wordCount, termCount int
			err := rows.Scan(&word, &urlTitle, &url, &wordCount, &termCount)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("word: %v, wordCount: %v, termCount: %v, numDocs: %v, docs_in_C: %v\n", word, wordCount, termCount, numDocs, docs_in_C)

			tf_idf := sqlCalcTFIDF(float64(termCount), float64(wordCount), float64(numDocs), float64(docs_in_C))
			tfidfScoresSlice = append(tfidfScoresSlice, TFIDFScore{Word: word, Title: urlTitle, URL: url, Score: tf_idf})
		}
	}
	defer rows.Close()
	sort.Sort(Hits(tfidfScoresSlice))

	return tfidfScoresSlice
}
