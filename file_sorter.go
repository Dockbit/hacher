package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
)

type File struct {
	Path string
	os.FileInfo
}

type Files []*File

type ByName struct{ Files }
type ByMtime struct{ Files }

func (a Files) Len() int      { return len(a) }
func (a Files) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func (a ByName) Less(i, j int) bool { return a.Files[i].Path < a.Files[j].Path }

func (a ByMtime) Less(i, j int) bool {
	return a.Files[i].ModTime().Before(a.Files[j].ModTime())
}

func (f File) String() string {
	return fmt.Sprintf("{%s: %v %d}", f.Path, f.ModTime(), f.Size())
}

/*
 * Sorts files by descending mtime with an optional regex filtering.
 */
func fileSorter(path string, filter ...string) Files {
	var files Files

	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err == nil && info.Mode().IsRegular() {
				matched := true

				// filter by regex
				if len(filter) > 0 {
					matched, _ = regexp.MatchString(filter[0], filepath.Base(path))
				}
				if matched {
					files = append(files, &File{Path: path, FileInfo: info})
				}
			}
			return nil
		},
	)
	checkError(err)
	sort.Sort(sort.Reverse(ByMtime{files}))

	return files
}
