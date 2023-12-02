package main

import (
	"fmt"
	"log"
	"strings"
	"unicode"

	"golang.org/x/net/html"
)

type ExtractResult struct {
	wordSlice, hrefSlice []string
	url, title           string
}

func extract(dOutCh chan DownloadResult, exoutC chan ExtractResult) {

	for dlStruct := range dOutCh {
		var ex ExtractResult
		f := func(c rune) bool {
			return !unicode.IsLetter(c) && !unicode.IsNumber(c)
		}

		//produces a tree of node from the from the string casted body
		tree, err := html.Parse(strings.NewReader(string(dlStruct.body)))
		if err != nil {
			log.Fatal(err)
		}

		//initializing an anonymous function
		var extractTree func(*html.Node)

		extractTree = func(n *html.Node) {

			//check if current node is a HTML element and if tag is an achor "a"
			if n.Type == html.ElementNode && n.Data == "a" {
				//loop through the attributes, if the key is "href" then append a.val to hrefs
				for _, a := range n.Attr {
					if a.Key == "href" {
						ex.hrefSlice = append(ex.hrefSlice, a.Val)
						break
					}
				}
			} else if n.Type == html.ElementNode && n.Data == "title" {
				ex.title = n.FirstChild.Data
				fmt.Printf("%s\n", n.FirstChild.Data)
			} else if n.Type == html.TextNode { //if the current node is a text node, append to words slice
				sliceWords := strings.Fields(n.Data)
				for _, everyWord := range sliceWords {
					// fields gets rid of random spacing
					word := strings.FieldsFunc(everyWord, f)
					ex.wordSlice = append(ex.wordSlice, word...)
				}
			} else if n.Type == html.ElementNode && (n.Data == "style" || n.Data == "script") {
				return
			}
			//recursivly traversing through all the child node of the tree
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				extractTree(c)
			}
		}
		extractTree(tree)
		exoutC <- ExtractResult{ex.wordSlice, ex.hrefSlice, dlStruct.url, ex.title}
	}
}
