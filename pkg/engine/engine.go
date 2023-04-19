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

const topN = 10 // return only topN documents as result of the search

/*
https://en.wikipedia.org/wiki/Tf%E2%80%93idf

- t: termino
- d: documento
- D: conjunto de todos los documentos

para cada termino, hay que calcular:
- term frequency - tf(t,d): relative frequency of term t within document d
- inverse document frequency - idf(t,D): measures how much info the term t provides across all the set of documents D

una vez tengamos estos datos para cada termino, podemos calcular el tf-idf, que nos dice lo relevante que es el termino
t en el documento d teniendo en cuenta todos los docs de D:

tfidf(t,d,D) = tf(t,d) * idf(t,D)

flow:
- user introduce una palabra
- calculamos tf de esa palabra para cada doc
- calculamos idf de esa palabra para todo el set de docs
- calculamos tf-idf
- devolvemos una lista de documentos ordenada por el valor de tf-idf


Si se hace una busqueda de varias palabras? eg:
- busco "linux kernel"
- obtengo 5 docs para linux y 5 para kernel

*/

type Engine struct {
	index []document.Document
}

// New creates a new instance of the search engine
func New() *Engine {
	return &Engine{} // maybe preallocate index with an estimate doc count
}

// String enables pretty printing of the engine
func (e Engine) String() string {
	var result string
	fmt.Printf("Engine has %d documents loaded\n", e.DocumentCount())
	for _, doc := range e.index {
		result += fmt.Sprintf("%s has %d words\n", doc.Path(), doc.WordCount())
	}
	return result
}

// DocumentCount returns the number of document that the engine has loaded
func (e Engine) DocumentCount() int {
	return len(e.index)
}

// Add adds a docIndex struct to the engine
func (e *Engine) add(doc document.Document) {
	e.index = append(e.index, doc)
}

// Load will traverse the docs under "path" which ends in ".html" and index them into memory
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
func tf(term string, d document.Document) float64 {
	numerator := float64(d.Occurrences(term)) // Number of times term t appears in a document
	denominator := float64(d.WordCount())     // Total number of terms in the document
	return numerator / denominator
}

// idf calculates the inverse document frequency of a term in the set of documents
func idf(term string, e Engine) float64 {
	numerator := float64(e.DocumentCount())                  // Total number of documents
	denominator := float64(e.CountDocsThatContainTerm(term)) // Number of documents with term t in it
	return math.Log(numerator / denominator)
}

// tf_idf calculates the tf-idf of a term in a document given a set of documents (loaded in the Engine)
func tf_idf(term string, d document.Document, e Engine) float64 {
	return tf(term, d) * idf(term, e)
}

// CountDocsThatContainTerm returns the count of docs in the engine that contains the term "term"
func (e Engine) CountDocsThatContainTerm(term string) int {
	var result int
	for _, d := range e.index {
		if d.Contains(term) {
			result++
		}
	}
	return result
}

// Search return an array of strings with the documents that are more relevant to show info about the term, ordered
func (e Engine) Search(term string) []string {
	// docRanked stores the doc path and the tf-idf value
	// we use this representation to be able to order the result before returning it
	type docRanked struct {
		path  string
		value float64
	}

	var rankings []docRanked

	for _, d := range e.index {
		dr := docRanked{
			path:  d.Path(),
			value: tf_idf(term, d, e),
		}
		rankings = append(rankings, dr)
	}

	sort.SliceStable(rankings, func(i, j int) bool {
		return rankings[i].value > rankings[j].value
	})

	for _, v := range rankings {
		fmt.Printf("Ranking of %s: %f\n", v.path, v.value)
	}

	var result []string

	// return only topN docs
	for i := 0; i < topN && topN < len(rankings); i++ {
		result = append(result, rankings[i].path)
	}

	return result
}
