package main

import (
	"encoding/json"
	"log"
	"os"
)

// stopWords creates a map of the stopwords to utlizalie the performance properties of a hash lookup, O(1)
func stopWords() map[string]bool {
	//read the json file
	byteSlice, err := os.ReadFile("./stopwords-en.json")
	if err != nil {
		log.Fatal(err)
	}

	stringSlice := []string{}

	//add bytes slice from json file into a string slice
	err2 := json.Unmarshal(byteSlice, &stringSlice)
	if err2 != nil {
		log.Fatal(err2)
	}

	stopMap := make(map[string]bool)

	//make a map of all the words in the string slice for O(1) look up in crawl()
	for _, val := range stringSlice {
		stopMap[val] = true
	}

	return stopMap

}
