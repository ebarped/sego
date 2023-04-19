package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ebarped/game-of-life/pkg/engine"
)

const (
	DOCS_DIR = "linux-kernel-docs"
)

func main() {
	e := engine.New()

	// index documentation
	start := time.Now()
	err := e.Load(DOCS_DIR)
	if err != nil {
		log.Fatalf("error loading: %s\n", err)
	}
	fmt.Println("Time elapsed:", time.Since(start))

	docs := e.Search("memory management")
	fmt.Println("Results:")
	for _, doc := range docs {
		fmt.Println("-", doc)
	}
}
