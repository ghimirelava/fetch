package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func makeImageQueries() {
	//create database
	db, err := sql.Open("sqlite3", "project04.db")
	if err != nil {
		fmt.Println("Err makeImageQueries() open: ", err)
		os.Exit(-1)
	}
	defer db.Close()

	// operate on the database

	//imageTerms table
	_, err = db.Exec(`CREATE TABLE image_terms(
			term_id integer PRIMARY KEY,
			word text UNIQUE
		);
	`)
	if err != nil {
		fmt.Println("Err creating image_terms table: ", err)
		os.Exit(-1)
	}

	//url table
	_, err = db.Exec(`CREATE TABLE image_urls( 
			url_id integer PRIMARY KEY,
			url text UNIQUE,
			url_title text,
			word_count integer NOT NULL
		);
	`)
	if err != nil {
		fmt.Println("Err creating urls table: ", err)
		os.Exit(-1)
	}

	//image_hits table
	_, err = db.Exec(`CREATE TABLE image_hits(
			image_hits_id integer PRIMARY KEY,
			term_id integer,
			url_id integer,
			src_url text,
			alt_text text,
			term_count integer NOT NULL,
			FOREIGN KEY(term_id) references terms(term_id),
			FOREIGN KEY(url_id) references urls(url_id)
		);
	`)
	if err != nil {
		fmt.Println("Err creating image_hits table: ", err)
		os.Exit(-1)
	}

}

func popImageTable(wordString, link, linkTitle, src, alt string) {
	db, err := sql.Open("sqlite3", "project04.db")
	if err != nil {
		fmt.Println("Err in popTermsTable() open: ", err)
		os.Exit(-1)
	}
	defer db.Close()

	//transaction
	tx, err := db.Begin()
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}
	defer tx.Commit()

	//image_terms table
	_, err2 := tx.Exec(`INSERT OR IGNORE INTO
		image_terms(word)
		values(?)
	`, wordString)
	if err2 != nil {
		fmt.Println("Err in insert into image_terms table : ", err2)
		os.Exit(-1)
	}

	//urls table
	_, err3 := tx.Exec(`INSERT OR IGNORE INTO
	image_urls(url, url_title, word_count)
		values(?, ?, ?)
	`, link, linkTitle, 0)
	if err3 != nil {
		fmt.Println("Err in insert into image_urls: ", err3)
		os.Exit(-1)
	}

	_, err4 := tx.Exec(`UPDATE image_urls SET word_count = word_count+1
		WHERE url = ?
	`, link)
	if err4 != nil {
		fmt.Println("Err in update image_urls table: ", err4)
		os.Exit(-1)
	}

	//image_hits table
	var temp_term_id int
	err = tx.QueryRow("SELECT term_id FROM image_terms WHERE word = ?;", wordString).Scan(&temp_term_id)

	var temp_url_id int
	err = tx.QueryRow("SELECT url_id FROM image_urls WHERE url = ?;", link).Scan(&temp_url_id)

	var temp_term_count int
	err = tx.QueryRow("SELECT term_count FROM image_hits WHERE url_id = ? AND term_id = ?;", temp_url_id, temp_term_id).Scan(&temp_term_count)

	if err == sql.ErrNoRows {
		_, err5 := tx.Exec(`INSERT INTO
			image_hits(term_id, url_id, src_url, alt_text, term_count)
			VALUES (?, ?, ?, ?, ?)
		`, temp_term_id, temp_url_id, src, alt, 1)

		if err5 != nil {
			fmt.Println("Err in insert into image_hits table: ", err5)
			os.Exit(-1)
		}
	} else if err == nil {
		_, err6 := tx.Exec(`UPDATE image_hits SET term_count = term_count+1
			WHERE term_id = ? AND url_id = ? AND src_url = ? AND alt_text = ?
		`, temp_term_id, temp_url_id, src, alt)

		if err6 != nil {
			fmt.Println("Err in update into image_hits table: ", err6)
			os.Exit(-1)
		}
	}

}
