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
	return []string{"title", "p", "h1", "h2", "h3", "pre", "li"}
}

// Document represents an indexed document: stores the path, every distinct
// term and the number of occurrences of the term inside the doc
type Document struct {
	Path       string         `json:"path"`        // path of the document
	WordsIndex map[string]int `json:"words_index"` // maps word -> nº of occurences
}

// New indexes an ".html" documents and returns a representation of it
func New(path string) Document {
	di, err := index(path)
	if err != nil {
		log.Fatalf("error indexing document %s: %s\n", path, err)
	}
	return di
}

// WordCount returns the number of distinct words in the document
func (d Document) WordCount() int {
	return len(d.WordsIndex)
}

// Occurrences return all the occurrences of the term in the doc
func (d Document) Occurrences(term string) int {
	return d.WordsIndex[term]
}

// Contains return true if the term exists in the doc
func (d Document) Contains(term string) bool {
	if _, ok := d.WordsIndex[term]; ok {
		return true
	}
	return false
}

// index indexes the ".html" document in "path"
func index(path string) (Document, error) {
	file, err := os.Open(path)
	if err != nil {
		return Document{}, fmt.Errorf("error opening file %s: %s\n", path, err)
	}
	defer file.Close()

	htmlFile, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		return Document{}, fmt.Errorf("error parsing .html file %s: %s\n", path, err)
	}

	var sb strings.Builder

	for _, tag := range tagsToExplore() {
		htmlFile.Find(tag).Each(func(index int, selection *goquery.Selection) {
			_, err := sb.WriteString(strings.ToLower(selection.Text()) + "\n")
			if err != nil {
				log.Printf("error parsing tag %s from file %s: %s\n", tag, file.Name(), "a")
			}
		})
	}

	removePunctuation := func(r rune) rune {
		if strings.ContainsRune(".,:;()+-", r) {
			return -1
		} else {
			return r
		}
	}

	words := strings.Map(removePunctuation, sb.String())

	doc := Document{
		Path:       path,
		WordsIndex: make(map[string]int),
	}

	for _, term := range strings.Fields(words) {
		if !doc.Contains(term) { // first occurrence ot the term
			doc.WordsIndex[term] = 1
		} else {
			doc.WordsIndex[term] = doc.WordsIndex[term] + 1
		}
	}

	// fmt.Printf("[DEBUG] %v\n", doc.wordsIndex)

	return doc, nil
}

// String allows pretty printing of Document
func (d Document) String() string {
	return fmt.Sprintf("%s words: %d\n", d.Path, d.WordCount())
}
