package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"bitbucket.org/kalafut/gosh"

	"github.com/spaolacci/murmur3"
	"gopkg.in/alecthomas/kingpin.v2"
)

const SAMPLE_SIZE = 4096
const FULL_HASH_LIMIT = 3 * SAMPLE_SIZE

const WORKERS = 10

var exts = gosh.NewSet(".jpg", ".mp4")

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
			//if exts.Contains(strings.ToLower(filepath.Ext(path))) {
			if info.Mode().IsRegular() {
				files <- &File{Path: path, Size: info.Size(), ModTime: info.ModTime()}
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

func listCatalog(name string) error {
	catalog, err := LoadCatalog(name)
	if err != nil {
		return err
	}
	catalog.List()

	return nil
}

func build(root string, name string) {
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

	catalog.Save(name)
}

func dedupe(keep string, kill string, doDelete bool) error {
	keepCat, err := LoadCatalog(keep)
	if err != nil {
		return err
	}
	killCat, err := LoadCatalog(kill)
	if err != nil {
		return err
	}
	files := BInA(keepCat, killCat)

	for _, f := range files {
		if doDelete {
			os.Remove(f)
			fmt.Print("Deleting: ")
		}
		fmt.Println(f)
	}

	return nil
}

func main() {
	var (
		app = kingpin.New("dsync", "Directory Synchronizer")

		buildCmd = app.Command("build", "Build catalog.")
		path     = buildCmd.Arg("path", "Root path").Required().String()
		name     = buildCmd.Arg("name", "Catalog name").Required().String()

		list      = app.Command("list", "List catalog.")
		list_name = list.Arg("name", "Catalog name").Required().String()

		dedupeCmd = app.Command("dedupe", "Dedupe")
		keep      = dedupeCmd.Arg("keep name", "Keep catalog").Required().String()
		kill      = dedupeCmd.Arg("kill name", "Kill catalog").Required().String()
		doDelete  = dedupeCmd.Flag("delete", "Delete files").Bool()
	)

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case buildCmd.FullCommand():
		build(*path, *name)

	// Post message
	case list.FullCommand():
		listCatalog(*list_name)

	case dedupeCmd.FullCommand():
		dedupe(*keep, *kill, *doDelete)
	}

}
