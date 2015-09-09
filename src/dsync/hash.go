package dsync

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	//"bitbucket.org/kalafut/gosh"

	//"bitbucket.org/kalafut/gosh"

	"github.com/spaolacci/murmur3"
)

const SAMPLE_SIZE = 4096
const FULL_HASH_LIMIT = 3 * SAMPLE_SIZE

const WORKERS = 10

//var exts = gosh.NewSet(".jpg", ".mp4")
//var excluded = gosh.NewSet(".DS_Store")

func smartHash(file string) (uint64, error) {
	h := murmur3.New64()

	f, err := os.Open(file)
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return 0, err
	}

	if fi.Size() <= FULL_HASH_LIMIT {
		buffer := make([]byte, fi.Size())
		f.Read(buffer)
		h.Write(buffer)
	} else {
		buffer := make([]byte, SAMPLE_SIZE)
		f.Read(buffer)
		h.Write(buffer)
		f.Seek(fi.Size()/2, 0)
		f.Read(buffer)
		h.Write(buffer)
		f.Seek(-SAMPLE_SIZE, 2)
		f.Read(buffer)
		h.Write(buffer)
	}

	return h.Sum64(), nil
}

func hashFiles(files <-chan *File, catalog *Catalog) {
	buffer := make([]byte, SAMPLE_SIZE)
	h := murmur3.New64()
	for file := range files {
		h.Reset()
		f, _ := os.Open(file.Path)
		f.Read(buffer)
		h.Write(buffer)
		f.Close()

		file.Hash = h.Sum64()
		catalog.AddFile(file)
	}
}

func traverse(root string) <-chan *File {
	files := make(chan *File)
	go func() {
		defer close(files)
		filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				panic(err)
			}
			valid := true //!excluded.Contains(filepath.Base(path)) && info.Mode().IsRegular()
			//if exts.Contains(strings.ToLower(filepath.Ext(path))) {

			if valid {
				files <- &File{Path: stdSlash(path), Size: info.Size(), ModTime: info.ModTime()}
			}
			//}
			return nil
		})
	}()
	return files
}

func monitor(catalog *Catalog) {
	for {
		fmt.Println(catalog.count)
		time.Sleep(1 * time.Second)
	}
}
