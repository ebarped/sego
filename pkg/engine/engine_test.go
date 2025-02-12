package engine

import (
	"testing"
)

const (
	DOCS_PATH          = "../../linux-kernel-docs.tgz"
	DOCS_PATH_DST      = "."
	SEARCH_QUERY       = "memory management"
	SAVE_STATE_PATH    = "/tmp/sego_index.json"
	QUERY_RESULT_COUNT = 5
)

var e *Engine

func init() {
	e = New()
	e.Load(DOCS_PATH)
}

func BenchmarkSearch(b *testing.B) {
	for b.Loop() {
		e.Search(SEARCH_QUERY, QUERY_RESULT_COUNT)
	}
}

func BenchmarkLoad(b *testing.B) {
	e := New()

	for b.Loop() {
		e.Load(DOCS_PATH)
	}
}

func BenchmarkSaveState(b *testing.B) {
	for b.Loop() {
		e.SaveState(SAVE_STATE_PATH)
	}
}

func BenchmarkUntargz(b *testing.B) {
	for b.Loop() {
		untargz(DOCS_PATH_DST, DOCS_PATH)
	}
}
