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

type RF struct {
	*Root
	*File
}

type Catalog struct {
	mutex  *sync.Mutex
	Files  []*File
	Roots  []*Root
	Hashes map[uint64][]RF
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
	hs := r.catalog.Hashes[file.Hash]
	hs = append(hs, struct {
		*Root
		*File
	}{r, file})
	r.catalog.Hashes[file.Hash] = hs
}

func (r *Root) RemoveFile(path string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	fs := r.Files

	for i, _ := range fs {
		if fs[i].Path == path {
			r.catalog.Unhash(RF{Root: r, File: fs[i]})
			fs[i] = fs[len(fs)-1]
			fs[len(fs)-1] = nil
			r.Files = fs[:len(fs)-1]
			break
		}
	}

	/*
		r.Files = append(r.Files, file)
		hs := r.catalog.Hashes[file.Hash]
		hs = append(hs, struct {
			*Root
			*File
		}{r, file})
		r.catalog.Hashes[file.Hash] = hs
	*/
}

func (c *Catalog) AddFile(file *File) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.Files = append(c.Files, file)
}

func (c *Catalog) Unhash(file RF) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	h, ok := c.Hashes[file.Hash]
	if ok {
		for i, _ := range h {
			if h[i] == file {
				fmt.Println(c.Hashes[file.Hash])
				h[i] = h[len(h)-1]
				c.Hashes[file.Hash] = h[:len(h)-1]
				fmt.Println(c.Hashes[file.Hash])
				break
			}
		}

	}
}

func (c *Catalog) List() {
	for _, f := range c.Files {
		fmt.Printf("%-30s   %d %x\n", f.Path, f.Size, f.Hash)
	}
}

func NewCatalog() *Catalog {
	return &Catalog{mutex: &sync.Mutex{}, Hashes: make(map[uint64][]RF)}
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
