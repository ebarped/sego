package engine

import (
	"testing"
)

const (
	DOCS_DIR        = "../../linux-kernel-docs"
	SEARCH_QUERY    = "memory management"
	SAVE_STATE_PATH = "/tmp/sego_index.json"
)

var e *Engine

func init() {
	e = New()
	e.Load(DOCS_DIR)
}

func BenchmarkSearch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		e.Search(SEARCH_QUERY)
	}
}

func BenchmarkLoad(b *testing.B) {
	e := New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.Load(DOCS_DIR)
	}
}

func BenchmarkSaveState(b *testing.B) {
	for i := 0; i < b.N; i++ {
		e.SaveState(SAVE_STATE_PATH)
	}
}
