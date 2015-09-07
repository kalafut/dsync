package dsync

import (
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
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
	Path    string
	Name    string
	Updated time.Time
}

// Catalog is a truly great type.
type Catalog struct {
	mutex  *sync.Mutex
	count  int
	Roots  map[string][]Root
	Hashes map[uint64][]*File
}

// AddFile adds a File to the catalog.
func (c *Catalog) AddFile(file *File) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	hash := file.Hash
	c.Hashes[hash] = append(c.Hashes[hash], file)
	c.count++
}

// RemoveFile adds a File to the catalog.
func (c *Catalog) RemoveFile(file *File) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	h, ok := c.Hashes[file.Hash]
	if ok {
		for i, _ := range h {
			if h[i] == file {
				h[i] = h[len(h)-1]
				c.Hashes[file.Hash] = h[:len(h)-1]
				break
			}
		}

	}
}

func (c *Catalog) List() []*File {
	var files = make([]*File, 0)

	for _, fa := range c.Hashes {
		for _, f := range fa {
			files = append(files, f)
		}
	}

	//for _, f := range files {
	//	fmt.Printf("%-30s   %d %x\n", f.Path, f.Size, f.Hash)
	//}

	return files
}

func (c *Catalog) Dupes() [][]*File {
	var dupes = make([][]*File, 0)

	for _, fa := range c.Hashes {
		if len(fa) > 1 {
			// Don't bother checking file size, for now
			dupes = append(dupes, fa)
		}
	}

	return dupes
}

// AddRoot adds a named path (aka "root") to the catalog.
func (c *Catalog) AddRoot(path, name string) {
	m, _ := os.Hostname()
	c.Roots[m] = append(c.Roots[m], Root{Path: filepath.ToSlash(path), Name: name})
}

func (c *Catalog) RootNames() []string {
	var r []string
	m, _ := os.Hostname()

	for _, root := range c.Roots[m] {
		r = append(r, root.Name)
	}

	return r
}

// GetPath returns the path associated with the named Root.
func (c *Catalog) RootPath(name string) (ret string) {
	m, _ := os.Hostname()
	for _, r := range c.Roots[m] {
		if r.Name == name {
			ret = r.Path
		}
	}

	return ret
}

func NewCatalog() *Catalog {
	return &Catalog{
		mutex:  &sync.Mutex{},
		Hashes: make(map[uint64][]*File),
		Roots:  make(map[string][]Root),
	}
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

	/* Disable for now
	for _, bFile := range b.Files {
		for _, aFile := range a.Files {
			if aFile.Size == bFile.Size && aFile.Hash == bFile.Hash {
				files = append(files, bFile.Path)
			}
		}
	}
	*/

	return files
}

func ListCatalog(name string) error {
	catalog, err := LoadCatalog(name)
	if err != nil {
		return err
	}
	catalog.List()

	return nil
}

func Dedupe(keep string, kill string, doDelete bool) error {
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

// UpdateRoot recursively scans the named root, adding/updating files as necessary.
// Files that are in the catalog but not found on disk will be removed if clean is true.
func (c *Catalog) UpdateRoot(name string, clean bool) {
	files := traverse(c.RootPath(name))
	hashFiles(files, c)
}
