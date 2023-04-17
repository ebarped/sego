package engine

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

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

una vez tengamos un mapeo de termino -> tfidf, al hacer una busqueda solo hay que mostrar la lista de docs de mayor a menor valor de tfidf

*/

// el index deberia mapear word -> struct {filepath: "aa", occurences: 3}

type data struct {
	filePath   string
	occurences int
}

type Engine struct {
	index map[string]data
}

func New() Engine {
	return Engine{
		index: make(map[string]data),
	}
}

func (e Engine) Add(word, path string, occ int) {
	e.index[word] = data{
		filePath:   path,
		occurences: occ,
	}
}

// represents a document, and stores every distinct word and the number of occurrences of it inside the doc
type docIndex struct {
	path       string
	wordsIndex map[string]int
}

func indexDoc(path string) docIndex {
	// solo nos tenemos que quedar con el texto que hay en las secciones <title>, <p>, <h1>, <h2> y <h3> del html
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("error opening file %s: %s\n", path, err)
	}

	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		fmt.Println("No url found")
		log.Fatal(err)
	}

	var words string

	doc.Find("title").Each(func(index int, selection *goquery.Selection) {
		words += selection.Text() + "\n"
	})

	doc.Find("p").Each(func(index int, selection *goquery.Selection) {
		words += selection.Text() + "\n"
	})

	doc.Find("h1").Each(func(index int, selection *goquery.Selection) {
		words += selection.Text() + "\n"
	})

	doc.Find("h2").Each(func(index int, selection *goquery.Selection) {
		words += selection.Text() + "\n"
	})

	doc.Find("h3").Each(func(index int, selection *goquery.Selection) {
		words += selection.Text() + "\n"
	})

	removePunctuation := func(r rune) rune {
		if strings.ContainsRune(".,:;()+-", r) {
			return -1
		} else {
			return r
		}
	}

	words = strings.Map(removePunctuation, words)

	docIndex := docIndex{
		path:       path,
		wordsIndex: make(map[string]int), //maybe preallocate with len(words)
	}

	for _, word := range strings.Fields(words) {
		if _, ok := docIndex.wordsIndex[word]; !ok {
			docIndex.wordsIndex[word] = 1
		} else {
			docIndex.wordsIndex[word] = docIndex.wordsIndex[word] + 1

		}
	}

	fmt.Printf("indice de %s: %v\n", docIndex.path, docIndex.wordsIndex)

	os.Exit(0)

	return docIndex
}

// Load will traverse the docs and index them into memory
func (e Engine) Load() fs.WalkDirFunc {
	return func(path string, d fs.DirEntry, err error) error {
		//fmt.Printf("Visited %s in path %s\n", d.Name(), path)
		fileName := d.Name()
		if !d.IsDir() && strings.HasSuffix(fileName, ".html") {
			indexDoc(path)
			//fmt.Printf("Indexando %s (%s)\n", fileName, path)
			e.Add(fileName, path, 3)
		}

		return nil
	}
}

func (e Engine) String() string {
	var result string
	for k, v := range e.index {
		result += fmt.Sprintf("%s -> %s (%d)\n", k, v.filePath, v.occurences)
	}
	return result
}
