package index

import (
	"fmt"
	"log"
	"org/go/helper/path"
	fp "path/filepath"
)

type Index struct {
	Files []*File
	Words Words
	Tags  Tags
	Root  string
}

func (ix *Index) NewPath(rel string) *path.Path {
	return &path.Path{
		Root: ix.Root,
		Rel:  rel,
	}
}

func (ix *Index) Read() error {
	if ix.Root == "" {
		return fmt.Errorf("Index.Root undefined")
	}
	files, err := ReadFiles(ix.Root, "/private/graph")
	if err != nil {
		return err
	}
	ix.Files = files
	return nil
}

func (ix *Index) Tokenize() {
	ix.Words = make(Words)
	ix.Words.AddFiles(ix.Files)
}

func (ix *Index) ParseTags() {
	ix.Tags = make(Tags)
	ix.Tags.AddFiles(ix.Files)
}

/*
func (ix *Index) AddFile(abs string) {
	f, err := ReadFile(ix.Root, abs)
	if err != nil {
		log.Println(fmt.Errorf("*index.Index.AddFile: %v", err))
		return
	}
	ix.TokenizeFile(f)
	ix.Files = append(ix.Files, f)
}
*/

func (ix *Index) TokenizeFile(f *File) {
	ix.Words.AddFile(f)
	ix.Tags.AddFile(f)
}

func (ix *Index) UpdateFile(rel string) {
	nf, err := ReadFile(ix.Root, fp.Join(ix.Root, rel))
	if err != nil {
		log.Println(fmt.Errorf("*index.Index.UpdateFile: %v", err))
		return
	}
	ix.TokenizeFile(nf)
	for i, f := range ix.Files {
		if f.Path == rel {
			ix.Files[i] = nf
			return
		}
	}
	ix.Files = append(ix.Files, nf)
	log.Printf("Couldnt update file %v", rel)
}

func (ix *Index) RemoveFile(path string) {
	for i, f := range ix.Files {
		if f.Path == path {
			ix.Files[i] = ix.Files[len(ix.Files)-1]
			ix.Files = ix.Files[:len(ix.Files)-1]
			return
		}
	}
	//TODO: file is not used from tokenized structures
	log.Printf("Couldnt remove file %v", path)
}
