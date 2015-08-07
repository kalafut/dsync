package main

import (
	"fmt"
	"hash/fnv"
	"os"
	"path/filepath"
	"sync"
	"time"

	"bitbucket.org/kalafut/gosh"

	"gopkg.in/alecthomas/kingpin.v2"
)

const SAMPLE_SIZE = 4096
const FULL_HASH_LIMIT = 3 * SAMPLE_SIZE

const WORKERS = 10

var exts = gosh.NewSet(".jpg", ".mp4")

func smartHash(file string) (uint64, error) {
	h := fnv.New64a()

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
	h := fnv.New64a()
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
			//if exts.Contains(strings.ToLower(filepath.Ext(path))) {
			files <- &File{Path: path, Size: info.Size(), ModTime: info.ModTime()}
			//}
			return nil
		})
	}()
	return files
}

func monitor(catalog *Catalog) {
	for {
		fmt.Println(len(catalog.Files))
		time.Sleep(1 * time.Second)
	}
}

func listCatalog() error {
	catalog, err := LoadCatalog("c.gob")
	if err != nil {
		return err
	}
	catalog.List()

	return nil
}

func build(root string) {
	var wg sync.WaitGroup
	catalog := NewCatalog()

	go monitor(catalog)

	files := traverse(root)

	for w := 0; w < WORKERS; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			hashFiles(files, catalog)
		}()
	}
	wg.Wait()

	catalog.Save("c.gob")
}

func main() {
	var (
		app = kingpin.New("dsync", "Directory Synchronizer")

		buildCmd = app.Command("build", "Build catalog.")
		path     = buildCmd.Arg("path", "Root path").Required().String()
		//registerName = register.Arg("name", "Name of user.").Required().String()

		list = app.Command("list", "List catalog.")
		//postImage = post.Flag("image", "Image to post.").File()
		//postChannel = post.Arg("channel", "Channel to post to.").Required().String()
		//postText = post.Arg("text", "Text to post.").Strings()
	)

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case buildCmd.FullCommand():
		build(*path)

	// Post message
	case list.FullCommand():
		listCatalog()
	}

}
