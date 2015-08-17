package main

import (
	"encoding/gob"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"
)

const DEFAULT_NAME = "catalog.gob"

type File struct {
	Path    string
	Size    int64
	Hash    uint64
	ModTime time.Time
}

type Root struct {
	catalog *Catalog
	mutex   *sync.Mutex
	Path    string
	Files   []*File
}

type Catalog struct {
	mutex  *sync.Mutex
	Files  []*File
	Roots  []*Root
	Hashes map[uint64][]struct {
		*Root
		*File
	}
}

func (c *Catalog) AddRoot(path string) (*Root, error) {
	for _, r := range c.Roots {
		if r.Path == path {
			return nil, errors.New("Root " + path + " already exists")
		}
	}

	root := &Root{Path: path, mutex: &sync.Mutex{}, catalog: c}
	c.Roots = append(c.Roots, root)

	return root, nil
}

func (r *Root) AddFile(file *File) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.Files = append(r.Files, file)
	//h := r.catalog.Hashes

}

func (c *Catalog) AddFile(file *File) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.Files = append(c.Files, file)
}

func (c *Catalog) List() {
	for _, f := range c.Files {
		fmt.Printf("%-30s   %d %x\n", f.Path, f.Size, f.Hash)
	}
}

func NewCatalog() *Catalog {
	return &Catalog{mutex: &sync.Mutex{}}
}

func LoadCatalog(filename string) (*Catalog, error) {
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		return nil, err
	}

	dec := gob.NewDecoder(f)

	c := NewCatalog()
	err = dec.Decode(c)

	return c, err
}

func (c *Catalog) Save(filename string) error {
	f, err := os.Create(filename)
	defer f.Close()

	enc := gob.NewEncoder(f)

	enc.Encode(c)

	return err
}

func BInA(a *Catalog, b *Catalog) []string {
	files := []string{}

	for _, bFile := range b.Files {
		for _, aFile := range a.Files {
			if aFile.Size == bFile.Size && aFile.Hash == bFile.Hash {
				files = append(files, bFile.Path)
			}
		}
	}

	return files
}
