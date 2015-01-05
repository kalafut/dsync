package main

import (
	"encoding/gob"
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

type Catalog struct {
	filename string
	mutex    *sync.Mutex
	Files    []*File
}

func (c *Catalog) AddFile(file *File) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.Files = append(c.Files, file)
}

func (c *Catalog) List() {
	for _, f := range c.Files {
		fmt.Printf("%-30s   %d\n", f.Path, f.Size)
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
