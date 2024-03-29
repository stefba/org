package crud

import (
	"io"
	"net/http"
	"org/go/helper"
	"org/go/helper/file"
	"org/go/helper/path"
	"org/go/index"
	"os"
	"path/filepath"
)

func RenameFile(ix *index.Index, w http.ResponseWriter, r *http.Request) *helper.Err {
	p := ix.NewPath(r.URL.Path[len("/api/move"):])

	e := &helper.Err{
		Func: "renameFile",
		Path: p.Rel,
		Code: 500,
	}

	newPath, err := getBodyPath(r)
	if err != nil {
		e.Err = err
		return e
	}

	newP := p.New(newPath)

	// dont like that.
	err = createBot(newP)
	if err != nil {
		e.Err = err
		return e
	}

	err = file.RenameSortEntry(p, newP)
	if err != nil {
		e.Err = err
		return e
	}

	err = os.Rename(p.Abs(), newP.Abs())
	if err != nil {
		e.Err = err
		return e
	}

	return nil
}

func deleteBot(path string) error {
	dir := filepath.Dir(path)
	if filepath.Base(dir) != "bot" {
		return nil
	}

	l, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	if len(l) == 0 {
		return os.Remove(dir)
	}

	return nil
}

func createBot(path *path.Path) error {
	dir := path.Parent()
	if dir.Base() != "bot" {
		return nil
	}

	_, err := os.Stat(dir.Abs())
	if err == nil {
		return nil
	}

	return os.Mkdir(dir.Abs(), 0755)
}

func getBodyPath(r *http.Request) (string, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}

	// TODO: make sure it’s a valid path
	p := string(body)

	return p, nil
}
