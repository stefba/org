package main

import (
	"net/http"
)

var ROOT = "data"
var BUILD = "build"

func main() {
	setup()
	go loadIndex()

	http.Handle("/", routes())
	http.ListenAndServe(":8334", nil)
}
