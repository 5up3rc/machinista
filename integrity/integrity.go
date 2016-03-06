// integrity.go
package main

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

func main() {
	root_dirPtr := flag.String("source", "./", "Source of file")
	data_filePtr := flag.String("dest", "./", "Destination of file data")
	skip_namePtr := flag.String("skip", "", "Skip dirs")
	flag.Parse()
	//-------------------------------------------------------------------------------------

	var loadedTree map[string]string
	response := load(&loadedTree, last_fileName(*data_filePtr), *root_dirPtr, *data_filePtr) //carica la map dal file
	if response == false {                                                                   //se si inizializza il sistema riesegue il caricamento
		fmt.Println("Inizializzo il sistema")
		SnapShotTree := walk_tree(*root_dirPtr, *skip_namePtr)
		store(SnapShotTree, last_fileName(*data_filePtr))
		//panic(err)
		load(&loadedTree, last_fileName(*data_filePtr), *root_dirPtr, *data_filePtr)
	}
	SnapShotTree := walk_tree(*root_dirPtr, *skip_namePtr) //legge e salva lo stato attuale
	fmt.Println("loaded elements: ", len(loadedTree), "Snap elements: ", len(SnapShotTree))

	//-------------------------------------------------------------------------------------
	//Esegue copia di archiviazione del file con le impronte
	_, err := copy(last_fileName(*data_filePtr), new_fileName(*data_filePtr))
	if err != nil {
		fmt.Println("Errore copia di backup file dati")
		panic(err)
	}
	//Salva il nuovo file delle impronte
	store(SnapShotTree, last_fileName(*data_filePtr))
	//-------------------------------------------------------------------------------------

	fmt.Println(check_tree(loadedTree, SnapShotTree))
}

func copy(src, dst string) (int64, error) {
	src_file, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer src_file.Close()

	src_file_stat, err := src_file.Stat()
	if err != nil {
		return 0, err
	}

	if !src_file_stat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	dst_file, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer dst_file.Close()
	return io.Copy(dst_file, src_file)
}

func new_fileName(path string) string {
	return fmt.Sprintf("%s/gofp-%d-%d.txt", path, time.Now().YearDay(), time.Now().Month())
}

func last_fileName(path string) string {
	return fmt.Sprintf("%s/gofp-last.txt", path)
}

//write data
func store(data interface{}, file string) {
	m := new(bytes.Buffer)
	enc := gob.NewEncoder(m)

	err := enc.Encode(data)
	if err != nil {
		fmt.Println("Errore scrittura file hash")
		panic(err)
	}
	err = ioutil.WriteFile(file, m.Bytes(), 0600)
	if err != nil {
		fmt.Println("Errore scrittura file hash 1")
		panic(err)
	}
}

//read data
func load(e interface{}, file string, input_path string, data_path string) bool {
	n, err := ioutil.ReadFile(file)
	if err != nil {
		return false
	}

	p := bytes.NewBuffer(n)
	dec := gob.NewDecoder(p)

	err = dec.Decode(e)
	if err != nil {
		fmt.Println("Errore apertura file hash")
		return false
	}
	return true
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

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
		//fmt.Println(fp, fmt.Sprintf("%x -- %s", sum, fi.Mode().Perm()))
		files[fp] = fmt.Sprintf("%x-- %s", sum, fi.Mode().Perm())
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
