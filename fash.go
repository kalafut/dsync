package main

import (
	"encoding/gob"
	"hash/fnv"
	"os"
	"path/filepath"
	"sync"
)

const SAMPLE_SIZE = 4096
const WORKERS = 10

var EXT = map[string]bool{
	".jpg": true,
	".mp4": true,
}

type File struct {
	Path string
	Size int64
	Hash uint64
}

type Catalog struct {
	mutex *sync.Mutex
	Files []File
}

func g(files <-chan File, results *Catalog) {
	buffer := make([]byte, SAMPLE_SIZE)
	h := fnv.New64a()
	for file := range files {
		h.Reset()
		f, _ := os.Open(file.Path)
		f.Read(buffer)
		h.Write(buffer)

		file.Hash = h.Sum64()
		results.mutex.Lock()
		results.Files = append(results.Files, file)
		results.mutex.Unlock()
	}
}

func traverse(root string) <-chan File {
	files := make(chan File)
	go func() {
		defer close(files)
		filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if true || EXT[filepath.Ext(path)] {
				files <- File{Path: path, Size: info.Size()}
			}
			return nil
		})
	}()
	return files
}

func main() {
	var wg sync.WaitGroup
	catalog := Catalog{mutex: &sync.Mutex{}}

	files := traverse(".")

	for w := 0; w < WORKERS; w++ {
		wg.Add(1)
		go func() {
			g(files, &catalog)
			wg.Done()
		}()
	}
	wg.Wait()

	f, _ := os.Create("out.gob")
	defer f.Close()
	enc := gob.NewEncoder(f)
	for _, file := range catalog.Files {
		enc.Encode(file)
	}
}
