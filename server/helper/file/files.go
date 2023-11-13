package file

import (
	"org/server/helper/path"
	"os"
	fp "path/filepath"
	"sort"
)

func ReadFiles(path *path.Path) ([]*File, error) {
	l, err := os.ReadDir(path.Abs())
	if err != nil {
		return nil, err
	}
	files := []*File{}
	for i, fi := range l {
		if fi.Name() == ".DS_Store" {
			continue
		}
		fpath := fp.Join(path.Rel, fi.Name())
		files = append(files, &File{
			Num:  i,
			Name: fi.Name(),
			Path: fpath,
			Type: GetFileType(fpath, fi.IsDir()),
		})
	}
	return files, nil
}

func GetFiles(path *path.Path) ([]*File, bool, error) {
	files, err := ReadFiles(path)
	if err != nil {
		return nil, false, err
	}

	for _, f := range files {
		if f.Type == "text" {
			err = f.Read()
			if err != nil {
				return nil, false, err
			}
		}
	}
	if !hasSort(path.Abs()) {
		return antoSort(files), false, err
	}

	sorted, err := parseSort(path)
	if err != nil {
		return nil, false, err
	}

	fresh := merge(sorted, files)

	return renumerate(fresh), true, nil
}

func renumerate(files []*File) []*File {
	for i, f := range files {
		f.Num = i
	}
	return files
}

func preSort(files []*File) []*File {
	nu := []*File{}
	for _, f := range files {
		base := fp.Base(f.Path)
		if fp.Ext(f.Path) == ".info" && hasFile(files, stripExt(f.Path)) {
			continue
		}
		if base == ".DS_Store" {
			continue
		}
		if base == "info" {
			nu = append([]*File{f}, nu...)
			continue
		}
		nu = append(nu, f)
	}

	return antoSort(nu)
}

func divide(all []*File) (dirs, info, files []*File) {
	sort := []*File{}
	for _, f := range all {
		if f.Type == "dir" {
			dirs = append(dirs, f)
			continue
		}
		if f.Name == "info" {
			info = append(info, f)
			continue
		}
		if f.Name == ".sort" {
			sort = append(sort, f)
			continue
		}
		files = append(files, f)
	}
	return dirs, info, append(files, sort...)
}

func antoSort(all []*File) []*File {
	dirs, info, files := divide(all)

	sort.Sort(Asc(dirs))
	sort.Sort(Asc(files))

	return append(dirs, append(info, files...)...)
}

func separate(all []*File) []*File {
	dirs, info, files := divide(all)
	return append(dirs, append(info, files...)...)
}

func hasFile(files []*File, path string) bool {
	for _, f := range files {
		if f.Path == path {
			return true
		}
	}
	return false
}

type Asc []*File

func (a Asc) Len() int           { return len(a) }
func (a Asc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Asc) Less(i, j int) bool { return a[i].Name < a[j].Name }