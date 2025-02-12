package document

import (
	"fmt"
	"log"
	"testing"
)

const path = "../../linux-kernel-docs/admin-guide/abi-testing.html"

func BenchmarkIndex(b *testing.B) {
	var di Document
	var err error

	for b.Loop() {
		di, err = index(path)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println(di)
}
