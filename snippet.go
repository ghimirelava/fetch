package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
)

func (d *Datasql) getSnippet(searchTerm string) string {
	for key, _ := range d.snippetMap {
		if strings.Contains(key, searchTerm) {
			return key
		}
	}
	return "no string found"
}

func (d *Datasql) popSnippet(s []string) {
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

	for _, val := range s {
		println(val)
		d.snippetMap[val]++
		//terms table
		_, err2 := tx.Exec(`INSERT OR IGNORE INTO
			snips(word)
			values(?)
			`, val)
		if err2 != nil {
			fmt.Println("Err in insert into snips table : ", err2)
			os.Exit(-1)
		}
	}
}
