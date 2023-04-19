package document

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// tagsToExplore represents a constant []string that returns the HTML tags from which we will obtain words
func tagsToExplore() []string {
	return []string{"title", "p", "h1", "h2", "h3", "pre"}
}

// Document represents an indexed document: stores the path, every distinct
// term and the number of occurrences of the term inside the doc
type Document struct {
	path       string         // path of the document
	wordsIndex map[string]int // maps word -> nº of occurences
}

// New returns a struct that holds the indexed data of a document
func New(path string) Document {
	di, err := index(path)
	if err != nil {
		log.Fatalf("error indexing document %s: %s\n", path, err)
	}
	return di
}

func (d Document) WordCount() int {
	return len(d.wordsIndex)
}

func (d Document) Path() string {
	return d.path
}

// Occurrences return all the occurrences of the term in the doc
func (d Document) Occurrences(term string) int {
	return d.wordsIndex[term]
}

// Contains return true if the term exists in the doc
func (d Document) Contains(term string) bool {
	if _, ok := d.wordsIndex[term]; ok {
		return true
	}
	return false
}

// index indexes the document in "path"
func index(path string) (Document, error) {
	file, err := os.Open(path)
	if err != nil {
		return Document{}, fmt.Errorf("error opening file %s: %s\n", path, err)
	}
	defer file.Close()

	htmlFile, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		return Document{}, fmt.Errorf("error parsing html file %s: %s\n", path, err)
	}

	var words string

	for _, tag := range tagsToExplore() {
		htmlFile.Find(tag).Each(func(index int, selection *goquery.Selection) {
			words += strings.ToLower(selection.Text()) + "\n"
		})
	}

	removePunctuation := func(r rune) rune {
		if strings.ContainsRune(".,:;()+-", r) {
			return -1
		} else {
			return r
		}
	}

	words = strings.Map(removePunctuation, words)

	doc := Document{
		path:       path,
		wordsIndex: make(map[string]int), // maybe preallocate with len(words)
	}

	for _, word := range strings.Fields(words) {
		if _, ok := doc.wordsIndex[word]; !ok {
			doc.wordsIndex[word] = 1
		} else {
			doc.wordsIndex[word] = doc.wordsIndex[word] + 1
		}
	}

	return doc, nil
}
