package files

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"g.rg-s.com/org/go/helper/path"
	"g.rg-s.com/org/go/helper/reqerr"
	"g.rg-s.com/org/go/index"
)

func CopyFile(ix *index.Index, w http.ResponseWriter, r *http.Request) *reqerr.Err {
	p := ix.NewPath(r.URL.Path[len("/api/copy"):])
	e := reqerr.New("CopyFile", p.Path)

	newPath, err := getBodyPath(r)
	if err != nil {
		return e.Set(err, 500)
	}

	newP := p.New(newPath)

	if strings.Contains(newP.Path, "/public/") {
		dir := newP.Parent()
		err := os.MkdirAll(dir.Abs(), 0755)
		if err != nil {
			return e.Set(err, 500)
		}
		err = createInfo(dir)
		if err != nil {
			return e.Set(err, 500)
		}
	}

	err = copyFileFunc(p, newP)
	if err != nil {
		return e.Set(fmt.Errorf("faulty target path: %v", newPath), 500)
	}

	return nil
}

func copyFileFunc(oldpath, newpath *path.Path) error {
	b, err := os.ReadFile(oldpath.Abs())
	if err != nil {
		return err
	}

	return os.WriteFile(newpath.Abs(), b, 0644)
}
