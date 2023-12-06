package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func dropTables() {
	db, err := sql.Open("sqlite3", "project04.db")
	if err != nil {
		fmt.Println("Err makeQuries() open: ", err)
		os.Exit(-1)
	}
	defer db.Close()

	_, err = db.Exec(` DROP TABLE IF EXISTS terms;
	DROP TABLE IF EXISTS urls;
	DROP TABLE IF EXISTS hits;
	DROP TABLE IF EXISTS bigram_hits;
	DROP TABLE IF EXISTS image_terms;
	DROP TABLE IF EXISTS image_urls;
	DROP TABLE IF EXISTS image_hits;
	DROP TABLE IF EXISTS snips;`)
	if err != nil {
		fmt.Println("Err dropping table: ", err)
		os.Exit(-1)
	}
}

func makeQueries() {
	//create database
	db, err := sql.Open("sqlite3", "project04.db")
	if err != nil {
		fmt.Println("Err makeQuries() open: ", err)
		os.Exit(-1)
	}
	defer db.Close()

	// operate on the database

	//terms table
	_, err = db.Exec(`CREATE TABLE terms(
			term_id integer PRIMARY KEY,
			word text UNIQUE
		);
	`)
	if err != nil {
		fmt.Println("Err creating terms table: ", err)
		os.Exit(-1)
	}

	//snips table
	_, err = db.Exec(`CREATE TABLE snips(
			snip_id integer PRIMARY KEY,
			sentence text UNIQUE
		);
	`)
	if err != nil {
		fmt.Println("Err creating snips table: ", err)
		os.Exit(-1)
	}

	//url table
	_, err = db.Exec(`CREATE TABLE urls( 
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

	//hits table
	_, err = db.Exec(`CREATE TABLE hits(
			hits_id integer PRIMARY KEY,
			term_id integer,
			url_id integer,
			term_count integer NOT NULL,
			snippet_id text,
			FOREIGN KEY(term_id) references terms(term_id),
			FOREIGN KEY(url_id) references urls(url_id),
			FOREIGN KEY(snippet_id) references snips(snippet_id)
		);
	`)
	if err != nil {
		fmt.Println("Err creating hits table: ", err)
		os.Exit(-1)
	}

	//bigram table
	_, err = db.Exec(`CREATE TABLE bigram_hits(
			bigram_id integer PRIMARY KEY,
			term_id_one integer,
			term_id_two integer,
			url_id integer,
			term_count integer NOT NULL,
			snippet_id text,
			UNIQUE (term_id_one, term_id_two, url_id),
			FOREIGN KEY(term_id_one) references terms(term_id),
			FOREIGN KEY(term_id_two) references terms(term_id),
			FOREIGN KEY(url_id) references urls(url_id),
			FOREIGN KEY(snippet_id) references snips(snippet_id)
		);
	`)
	if err != nil {
		fmt.Println("Err creating bigram_hits table: ", err)
		os.Exit(-1)
	}

}

func popBiGram(term_id_one, term_id_two, url_id int) {
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

	_, err5 := tx.Exec(`INSERT OR IGNORE INTO
		bigram_hits(term_id_one, term_id_two, url_id, term_count)
			VALUES (?, ?, ?, ?)
		`, term_id_one, term_id_two, url_id, 0)

	if err5 != nil {
		fmt.Println("Err in insert into hits table: ", err5)
		os.Exit(-1)
	}
	_, err6 := tx.Exec(`UPDATE bigram_hits SET term_count = term_count+1
			WHERE term_id_one = ? AND term_id_two = ? AND url_id = ?
		`, term_id_one, term_id_two, url_id)

	if err6 != nil {
		fmt.Println("Err in update into bigram_hits table: ", err6)
		os.Exit(-1)
	}
}

func popTables(wordString, link, linkTitle, sentence string) (int, int) {
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

	//terms table
	_, err2 := tx.Exec(`INSERT OR IGNORE INTO
		terms(word)
		values(?)
	`, wordString)
	if err2 != nil {
		fmt.Println("Err in insert into terms table : ", err2)
		os.Exit(-1)
	}

	//urls table
	_, err3 := tx.Exec(`INSERT OR IGNORE INTO
		urls(url, url_title, word_count)
		values(?, ?, ?)
	`, link, linkTitle, 0)
	if err3 != nil {
		fmt.Println("Err in insert into urls: ", err3)
		os.Exit(-1)
	}

	_, err4 := tx.Exec(`UPDATE urls SET word_count = word_count+1
		WHERE url = ?
	`, link)
	if err4 != nil {
		fmt.Println("Err in update urls table: ", err4)
		os.Exit(-1)
	}

	//hits table
	var temp_term_id int
	err = tx.QueryRow("SELECT term_id FROM terms WHERE word = ?;", wordString).Scan(&temp_term_id)

	var temp_snip_id int
	err = tx.QueryRow("SELECT snip_id FROM snips WHERE sentence = ?;", sentence).Scan(&temp_snip_id)

	var temp_url_id int
	err = tx.QueryRow("SELECT url_id FROM urls WHERE url = ?;", link).Scan(&temp_url_id)

	var temp_term_count int
	err = tx.QueryRow("SELECT term_count FROM hits WHERE url_id = ? AND term_id = ?;", temp_url_id, temp_term_id).Scan(&temp_term_count)

	if err == sql.ErrNoRows {
		_, err5 := tx.Exec(`INSERT INTO
			hits(term_id, url_id, term_count)
			VALUES (?, ?, ?)
		`, temp_term_id, temp_url_id, 1)

		if err5 != nil {
			fmt.Println("Err in insert into hits table: ", err5)
			os.Exit(-1)
		}
	} else if err == nil {
		_, err6 := tx.Exec(`UPDATE hits SET term_count = term_count+1
			WHERE term_id = ? AND url_id = ?
		`, temp_term_id, temp_url_id)

		if err6 != nil {
			fmt.Println("Err in update into hits table: ", err6)
			os.Exit(-1)
		}
	}

	return temp_term_id, temp_url_id
}
