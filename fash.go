package main

import (
	"hash/fnv"
	"os"
	"path/filepath"
)

const SAMPLE_SIZE = 4096

func g(path string) uint64 {
	buffer := make([]byte, SAMPLE_SIZE)
	h := fnv.New64a()

	f, _ := os.Open(path)
	f.Read(buffer)
	h.Write(buffer)

	return h.Sum64()
}

func main() {
	g("test2")
	traverse("d://tmp//t7")
}

func traverse(root string) {
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		go g(path)
		if filepath.Ext(path) == ".c" {
		}
		return nil
	})
}
