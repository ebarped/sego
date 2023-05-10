package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ebarped/sego/pkg/engine"
	"github.com/ebarped/sego/pkg/server"
)

const (
	DOCS_PATH  = "linux-kernel-docs.tgz"
	STATE_PATH = "index.json"
	API_PORT   = "4000"
)

func main() {
	index := flag.Bool("index", false, "index the documentation")
	serve := flag.Bool("serve", false, "starts server on port "+API_PORT)
	flag.Parse()

	if flag.NFlag() != 1 {
		fmt.Println("Wrong number of flags")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// index documentation
	if *index {
		start := time.Now()
		e := engine.New()
		err := e.Load(DOCS_PATH)
		if err != nil {
			log.Fatalf("error loading: %s\n", err)
		}
		err = e.SaveState(STATE_PATH)
		fmt.Println("Time elapsed:", time.Since(start))
	}

	// search
	if *serve {
		e := engine.New(engine.WithState(STATE_PATH))
		s := server.New(API_PORT, e)
		s.Start()
	}

	os.Exit(0)
}
