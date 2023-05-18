package engine

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"log"
	"math"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/bytedance/sonic/decoder"
	"github.com/bytedance/sonic/encoder"

	"github.com/charlievieth/fastwalk"
	"github.com/ebarped/sego/pkg/document"
)

// Engine has a slice of indexed documents
type Engine struct {
	Index []document.Document `json:"index"`
}

// OptFunc enables constructing a new engine with options
type OptFunc func(Engine) Engine

// New creates a new instance of the search engine
func New(opts ...OptFunc) *Engine {
	var e Engine // maybe preallocate index with an estimate doc count

	// apply opts
	for _, opt := range opts {
		e = opt(e)
	}

	return &e
}

// documentCount returns the number of document that the engine has loaded
func (e Engine) documentCount() int {
	return len(e.Index)
}

// Add adds a document.Document to the engine
func (e *Engine) add(doc document.Document) {
	e.Index = append(e.Index, doc)
}

// Load will traverse decompress the "docs_path" file, traverse files docs that
// ends in ".html", indexing & storing them in the engine
func (e *Engine) Load(docs_path string) error {
	var mu sync.Mutex

	// decompress documentation
	err := untargz(".", docs_path)
	path := strings.TrimSuffix(docs_path, ".tgz")

	// traverse doc files, index them & store them into the engine
	conf := fastwalk.Config{
		Follow: false,
	}
	err = fastwalk.Walk(&conf, path, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".html") {
			fmt.Printf("Indexing %q\n", path)
			doc := document.New(path)

			mu.Lock()
			e.add(doc)
			mu.Unlock()
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// SaveState will save the state of the engine to disk
func (e Engine) SaveState(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("error: cannot create file to save the state on %s: %s\n", path, err)
	}

	encoder.NewStreamEncoder(f).Encode(e)
	if err != nil {
		return fmt.Errorf("error: cannot store engine state as json: %s\n", err)
	}

	return nil
}

// WithState will load the state from the index.json file into the engine
func WithState(path string) OptFunc {
	return func(e Engine) Engine {
		f, err := os.Open(path)
		if err != nil {
			log.Fatalf("error: cannot open file to load state from %s: %s\n", path, err)
		}

		err = decoder.NewStreamDecoder(f).Decode(&e)
		if err != nil {
			log.Fatalf("error: cannot load engine state from json %s: %s\n", path, err)
		}

		return e
	}
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
	// fmt.Printf("[DEBUG] term: %s, idf-num=%.12f, idf-den=%.12f, idf=%.12f\n", term, numerator, denominator, math.Log10(numerator/denominator))
	return math.Log10(numerator/denominator) + 1
}

// countDocsThatContainTerm returns the count of docs in the engine that contains the term "term"
func (e Engine) countDocsThatContainTerm(term string) int {
	var result int
	for _, d := range e.Index {
		if d.Contains(term) {
			result++
		}
	}
	return result
}

// Search return an array of strings with the documents that are more relevant to show info about the term, ordered
func (e Engine) Search(query string, resultCount int) []string {
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
	for _, doc := range e.Index {
		dr := docRanked{
			path: doc.Path,
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
	for i := 0; i < resultCount && resultCount < len(ranking); i++ {
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
	for _, doc := range e.Index {
		result += fmt.Sprintf("%s", doc)
	}
	return result
}

// untargz decompreses a .tar.gz file
func untargz(dst, src string) error {
	gzippedFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("untargz: error opening file %s: %w", src, err)
	}
	defer gzippedFile.Close()

	gzr, err := gzip.NewReader(gzippedFile)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	var header *tar.Header
	for header, err = tr.Next(); err == nil; header, err = tr.Next() {
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(header.Name, 0755); err != nil {
				return fmt.Errorf("untargz: Mkdir() failed: %w", err)
			}
		case tar.TypeReg:
			outFile, err := os.Create(header.Name)
			if err != nil {
				return fmt.Errorf("untargz: Create() failed: %w", err)
			}

			if _, err := io.Copy(outFile, tr); err != nil {
				// outFile.Close error omitted as Copy error is more interesting at this point
				outFile.Close()
				return fmt.Errorf("untargz: Copy() failed: %w", err)
			}
			if err := outFile.Close(); err != nil {
				return fmt.Errorf("untargz: Close() failed: %w", err)
			}
		default:
			return fmt.Errorf("untargz: uknown type: %b in %s", header.Typeflag, header.Name)
		}
	}
	if err != io.EOF {
		return fmt.Errorf("untargz: Next() failed: %w", err)
	}
	return nil
}
