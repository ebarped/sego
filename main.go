package main

import (
	"log"
	"path/filepath"

	"github.com/ebarped/game-of-life/pkg/engine"
)

const (
	DOCS_DIR = "linux-core-api-docs"
)

func main() {
	e := engine.New()

	err := filepath.WalkDir(DOCS_DIR, e.Load())
	if err != nil {
		log.Fatalf("fatal error traversing docs: %s\n", err)
	}

	//fmt.Printf("Engine status: %s\n", e)
}
