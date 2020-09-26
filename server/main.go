package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"github.com/gorilla/mux"
)

func main() {
	http.Handle("/", routes())
	http.ListenAndServe(":8334", nil)
}

func routes() *mux.Router {
	r := mux.NewRouter()

	r.PathPrefix("/file/").HandlerFunc(serveStatic)

	api := r.PathPrefix("/api/").Subrouter()
	api.Methods("GET").Queries("listing", "true").HandlerFunc(dirListing)
	api.Methods("GET").HandlerFunc(view)
	api.Methods("POST").HandlerFunc(writeSwitch)
	api.Methods("DELETE").HandlerFunc(deleteFile)

	return r
}

type Err struct {
	Func string
	Path string
	Err  error
}

var ROOT = "org"

func serveStatic(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len("/file"):]

	if fileType(path) == "text" {
		b, err := ioutil.ReadFile(ROOT+path)
		if err != nil {
			http.Error(w, err.Error(), 500)
			log.Println(err)
			return
		}
		fmt.Fprintf(w, "%s", removeNewLine(b))  
		return
	}
	http.ServeFile(w, r, ROOT+path)
}

func deleteFile(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len("/api"):]

	err := os.Remove(ROOT + path)
	if err != nil {
		err = fmt.Errorf("deleteFile: %v", err.Error())
		http.Error(w, err.Error(), 500)
		log.Println(err)
		return
	}
}

func writeSwitch(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len("/api"):]

	if strings.Contains(path, ".") || filepath.Base(path) == "info" {
		writeFile(w, r)
		return
	}
	createDir(w, r)
}

func createDir(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len("/api"):]
	fi, err := os.Stat(ROOT + filepath.Dir(path))
	if err != nil {
		err = fmt.Errorf("createDir: %v", err.Error())
		http.Error(w, err.Error(), 500)
		log.Println(err)
		return
	}
	if !fi.IsDir() {
		http.Error(w, "Can’t create dir in non-dir.", 500)
		return
	}
	err = os.Mkdir(ROOT+path, 0755)
	if err != nil {
		err = fmt.Errorf("createDir: %v", err.Error())
		http.Error(w, err.Error(), 500)
		log.Println(err)
		return
	}
	return
}

func writeFile(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len("/api"):]

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err = fmt.Errorf("writeFile: %v", err.Error())
		http.Error(w, err.Error(), 500)
		log.Println(err)
		return
	}

	body = removeMultipleNewLines(body)
	body = addNewLine(body)

	err = ioutil.WriteFile(ROOT+path, body, 0664)
	if err != nil {
		err = fmt.Errorf("writeFile: %v", err.Error())
		http.Error(w, err.Error(), 500)
		log.Println(err)
	}
	log.Printf("writeFile:\n{%s}\n", body)
}

func view(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len("/api"):]

	fi, err := os.Stat(ROOT + path)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	v := &View{
		File: &File{
			Path: path,
			Type: getFileType(path, fi.IsDir()),
		},
		Parent: filepath.Dir(path),
	}

	err = json.NewEncoder(w).Encode(v)
	if err != nil {
		err = fmt.Errorf("view: %v", err.Error())
		http.Error(w, err.Error(), 500)
		log.Println(err)
		return
	}
}
