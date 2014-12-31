package main

import (
	"fmt"
	"hash/fnv"
	"os"
	"path/filepath"
	"sync"
)

const SAMPLE_SIZE = 4096
const WORKERS = 10

var wg, wg2 sync.WaitGroup

var EXT = map[string]bool{
	".jpg": true,
	".mp4": true,
}

type File struct {
	Path string
	Size int64
	Hash uint64
}

func g(files <-chan File, results chan<- File) {
	buffer := make([]byte, SAMPLE_SIZE)
	h := fnv.New64a()
	for file := range files {
		h.Reset()
		f, _ := os.Open(file.Path)
		f.Read(buffer)
		h.Write(buffer)

		file.Hash = h.Sum64()
		results <- file
	}
	wg.Done()
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
	results := make(chan File, 100)

	files := traverse(".")

	wg2.Add(1)
	go func() {
		for file := range results {
			fmt.Printf("%s %d\n", file.Path, file.Hash)
		}
		wg2.Done()
	}()

	for w := 0; w < WORKERS; w++ {
		wg.Add(1)
		go g(files, results)
	}
	wg.Wait()
	close(results)
	wg2.Wait()

}
