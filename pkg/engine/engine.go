package engine

import (
	"fmt"
	"io/fs"
	"math"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ebarped/game-of-life/pkg/document"
)

const topN = 5 // return only topN documents as result of the search

// Engine has a slice of indexed documents
type Engine struct {
	index []document.Document
}

// New creates a new instance of the search engine
func New() *Engine {
	return &Engine{} // maybe preallocate index with an estimate doc count
}

// documentCount returns the number of document that the engine has loaded
func (e Engine) documentCount() int {
	return len(e.index)
}

// Add adds a document.Document to the engine
func (e *Engine) add(doc document.Document) {
	e.index = append(e.index, doc)
}

// Load will traverse the "path" folders locating docs that ends in ".html", indexing & storing them in the engine
func (e *Engine) Load(path string) error {
	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		fileName := d.Name()
		if !d.IsDir() && strings.HasSuffix(fileName, ".html") {
			fmt.Printf("Indexing %q\n", path)
			doc := document.New(path)
			if err != nil {
				return err
			}
			e.add(doc)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// tf calculates the term frequency of a term in the indexed document
// there is one value of tf for each term in each document
func tf(term string, d document.Document) float64 {
	numerator := float64(d.Occurrences(term)) // Number of times term t appears in a document
	denominator := float64(d.WordCount())     // Total number of terms in the document
	// fmt.Printf("[DEBUG] term: %s, doc: %s, tf-num=%.12f, tf-den=%.12f, tf=%.12f\n", term, d.Path(), numerator, denominator, numerator/denominator)
	return numerator / denominator
}

// idf calculates the inverse document frequency of a term in the set of documents
// there is one value of idf for each term in all the documents
func idf(term string, e Engine) float64 {
	numerator := float64(e.documentCount())                      // Total number of documents
	denominator := 1 + float64(e.countDocsThatContainTerm(term)) // Number of documents with term t in it
	//fmt.Printf("[DEBUG] term: %s, idf-num=%.12f, idf-den=%.12f, idf=%.12f\n", term, numerator, denominator, math.Log10(numerator/denominator))
	return math.Log10(numerator/denominator) + 1
}

// countDocsThatContainTerm returns the count of docs in the engine that contains the term "term"
func (e Engine) countDocsThatContainTerm(term string) int {
	var result int
	for _, d := range e.index {
		if d.Contains(term) {
			result++
		}
	}
	return result
}

// Search return an array of strings with the documents that are more relevant to show info about the term, ordered
func (e Engine) Search(query string) []string {
	query = strings.ToLower(query)
	queryTerms := strings.Fields(query)

	// docRanked stores the doc path and the tf-idf value
	// we use this representation to be able to order the result before returning it
	type docRanked struct {
		path  string
		value float64
	}

	var ranking []docRanked

	// precalculate IDF of each term of the query
	termsIDF := make(map[string]float64)
	for _, t := range queryTerms {
		termsIDF[t] = idf(t, e)
	}

	// for each document, calculate the tf-idf of each term of the query and sum them
	for _, doc := range e.index {
		dr := docRanked{
			path: doc.Path(),
		}
		for _, term := range queryTerms {
			dr.value += tf(term, doc) * termsIDF[term]
		}
		ranking = append(ranking, dr)
	}

	sort.SliceStable(ranking, func(i, j int) bool {
		return ranking[i].value > ranking[j].value
	})

	var result []string

	// return only topN docs
	fmt.Printf("[DEBUG] Results\n")
	for i := 0; i < topN && topN < len(ranking); i++ {
		fmt.Printf("[DEBUG] - %s -> %.12f\n", ranking[i].path, ranking[i].value)
		result = append(result, ranking[i].path)
		// fmt.Printf("[DEBUG] Ranking of %s: %f\n", ranking[i].path, ranking[i].value)
	}

	return result
}

// String enables pretty printing of the engine
func (e Engine) String() string {
	var result string
	fmt.Printf("Engine has %d documents loaded\n", e.documentCount())
	for _, doc := range e.index {
		result += fmt.Sprintf("%s", doc)
	}
	return result
}
