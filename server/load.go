package main

import (
	"log"
	"org/server/index"
	"time"
)

var IX = &index.Index{}

func loadIndex() error {
	start := time.Now()
	IX.Root = ROOT
	err := IX.Read()
	if err != nil {
		return err
	}
	log.Printf("Read %d documents in %v", len(IX.Files), time.Since(start))

	start = time.Now()
	IX.Tokenize()
	log.Printf("Tokenized %d documents in %v", len(IX.Files), time.Since(start))

	start = time.Now()
	IX.ParseTags()
	log.Printf("Extracted tags in %v", time.Since(start))

	return nil
}
