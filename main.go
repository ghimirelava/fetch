package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	fmt.Println("project06")

	dropTables()
	makeQueries()
	makeImageQueries()
	go http.ListenAndServe("localhost:8080", nil)
	url := "https://www.ucsc.edu/robots.txt"
	sqlHandleSearch(url)

	for {
		time.Sleep(2 * time.Second)
	}
}
