package main

import (
	"fmt"
	"log"

	"github.com/ebarped/game-of-life/pkg/engine"
)

const (
	DOCS_DIR = "linux-core-api-docs"
)

func main() {
	e := engine.New()

	// index documentation
	err := e.Load(DOCS_DIR)
	if err != nil {
		log.Fatalf("error loading: %s\n", err)
	}

	fmt.Printf("Engine status: \n%s", e)

	docs := e.Search("api")
	for _, doc := range docs {
		fmt.Println("-", doc)
	}
}
