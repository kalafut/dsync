package main

import "sync"

type Catalog struct {
	mutex *sync.Mutex
	Files []*File
}

func (c *Catalog) AddFile(file *File) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.Files = append(c.Files, file)
}

func NewCatalog() *Catalog {
	return &Catalog{mutex: &sync.Mutex{}}
}
