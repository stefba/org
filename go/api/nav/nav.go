package nav

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"g.rg-s.com/org/go/helper/file"
	"g.rg-s.com/org/go/helper/path"
	"g.rg-s.com/org/go/helper/reqerr"
	"g.rg-s.com/org/go/index"
)

type Nav struct {
	// Path is necessary here to give this object some sort of a life cycle.
	// React’s renderings are much faster than the server-side response.
	// Therefore, it has to know with what object it is dealing with.
	// Path     string  `json:"path"`
	Switcher string       `json:"switcher"`
	Siblings []*file.File `json:"siblings"`
}

func ViewNav(ix *index.Index, w http.ResponseWriter, r *http.Request) *reqerr.Err {
	p := ix.NewPath(r.URL.Path[len("/api/nav"):])
	e := reqerr.New("nav.ViewNav", p.Path)

	nav, err := getNav(p)
	if err != nil {
		return e.Set(err, 404)
	}

	err = json.NewEncoder(w).Encode(nav)
	if err != nil {
		return e.Set(err, 500)
	}

	return nil
}

func getNav(p *path.Path) (*Nav, error) {
	siblings := []*file.File{}

	if strings.Count(p.Path, "/") >= 2 {
		s, err := getSiblings(p)
		if err != nil {
			return nil, err
		}
		siblings = s
	}

	return &Nav{
		Siblings: siblings,
		Switcher: switchPath(p),
	}, nil
}

func getSiblings(p *path.Path) ([]*file.File, error) {
	parent := p.Parent()
	files, _, err := file.GetFiles(parent)
	if err != nil {
		return nil, err
	}

	files = dirsOnly(files)

	c := 0
	for i, f := range files {
		if f.Path == p.Path {
			c = i
			break
		}
	}

	length := len(files)

	start := 0
	end := length

	d := 2

	if c+1+d < length {
		end = c + 1 + d
	}
	if c-d > 0 {
		start = c - d
	}

	return files[start:end], nil
}

// Switch path is the closest corresponding public or private path to a directory.
//
//	/public/graph/20/20-10/10/subdir
//
// -> /private/graph/20/20-10
func switchPath(p *path.Path) string {
	public := false
	if l := len("/public"); len(p.Path) > l {
		public = p.Path[:l] == "/public"
	}

	var find, replace string

	if public {
		find = "public"
		replace = "private"
	} else {
		find = "private"
		replace = "public"
	}

	newP := p.New(strings.Replace(p.Path, find, replace, -1))

	existent := findExistent(newP)

	if existent == "/private" || existent == "/public" {
		return ""
	}
	return existent
}

// Recursive function going upwards the tree until it finds a existent
// directory.
func findExistent(p *path.Path) string {
	if p.Path == "/" || p.Path == "." {
		return p.Path
	}
	_, err := os.Stat(p.Abs())
	if err == nil {
		return p.Path
	}
	return findExistent(p.Parent())
}

func dirsOnly(files []*file.File) []*file.File {
	nu := []*file.File{}
	for _, f := range files {
		if f.Type == "dir" {
			nu = append(nu, f)
		}
	}
	return nu
}
