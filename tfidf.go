package main

type TFIDFScore struct {
	Word  string
	Title string
	URL   string
	Score float64
}
type Hits []TFIDFScore

func (s Hits) Len() int {
	return len(s)
}

func (s Hits) Less(i, j int) bool {
	if s[i].Score == s[j].Score {
		return s[i].URL > s[j].URL
	}
	return s[i].Score > s[j].Score
}

func (s Hits) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func sqlCalcTFIDF(termCount, termsInDoc, numDocs, docTermOCCint float64) float64 {

	tf := termCount / termsInDoc  //TF = occurrences in doc / number of stemmed words in doc
	df := numDocs / docTermOCCint //DF = number of documents the term occurs in / number of docs --> numDocs --> IDF = 1/DF
	idf := 1 / df
	tf_idf := tf * idf //TF-IDF = TF * IDF

	return tf_idf

}
