package engine

import (
	"testing"
)

const DOCS_DIR = "../../linux-kernel-docs"

var e *Engine

func init() {
	e = New()
	e.Load(DOCS_DIR)
}

func BenchmarkSearch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		e.Search("memory management")
	}
}
