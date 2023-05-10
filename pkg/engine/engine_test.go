package engine

import (
	"testing"
)

const (
	DOCS_PATH       = "../../linux-kernel-docs.tgz"
	DOCS_PATH_DST   = "."
	SEARCH_QUERY    = "memory management"
	SAVE_STATE_PATH = "/tmp/sego_index.json"
)

var e *Engine

func init() {
	e = New()
	e.Load(DOCS_PATH)
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
		e.Load(DOCS_PATH)
	}
}

func BenchmarkSaveState(b *testing.B) {
	for i := 0; i < b.N; i++ {
		e.SaveState(SAVE_STATE_PATH)
	}
}

func BenchmarkUntargz(b *testing.B) {
	for i := 0; i < b.N; i++ {
		untargz(DOCS_PATH_DST, DOCS_PATH)
	}
}
