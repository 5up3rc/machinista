// integrity.go
package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

//explore tree and return a map of file
func walk_tree(input_path string, skip_name string) map[string]string {
	files := make(map[string]string)
	skip := strings.Split(skip_name, " ")
	for i := range skip {
		skip[i] = filepath.Join(input_path, "/", skip[i])
	}
	n := 0
	VisitFile := func(fp string, fi os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err) // can't walk here,
			return nil       // but continue walking elsewhere
		}
		if fi.IsDir() {
			if stringInSlice(fp, skip) && len(skip) > 1 {
				fmt.Println("skip: ", fp)
				return filepath.SkipDir
			} else {
				return nil // not a file.  ignore.
			}
		}
		f, err := os.Open(fp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		}
		defer f.Close()
		h := md5.New()
		io.Copy(h, f)
		sum := h.Sum(nil)
		//fmt.Println(fp, fmt.Sprintf("%x", sum))
		files[fp] = fmt.Sprintf("%x", sum)
		n += 1
		return nil
	}
	filepath.Walk(input_path, VisitFile)
	fmt.Println("total files scanned: ", n)
	return files

}

func check_tree(loadedTree map[string]string, SnapShotTree map[string]string) string {
	var output string
	eq := reflect.DeepEqual(SnapShotTree, loadedTree) //compara due map
	if eq {
		output += "The file system is immutated"
	} else {
		for key, value := range SnapShotTree {
			_, ok := loadedTree[key]
			if !ok {
				output += fmt.Sprintln(key, "Created new file\n")
			} else {
				if value != loadedTree[key] {
					output += fmt.Sprintln(key, value, loadedTree[key], "Modified\n")
				}
			}
		}
	}
	return output
}

func main() {
	fmt.Println("Hello World!")
}
