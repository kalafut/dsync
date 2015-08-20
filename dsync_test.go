package main

import (
	"testing"

	"gopkg.in/tylerb/is.v1"
)

const TEST_DATA = "./test_data"

var (
	c   *Catalog
	r   *Root
	err error
	f1  = &File{Path: "f1", Size: 50, Hash: 42}
	f2  = &File{Path: "f2", Size: 150, Hash: 95}
	f3  = &File{Path: "f3", Size: 50, Hash: 42}
)

//Test that duplicate root aren't allowed
func TestDupeRoot(t *testing.T) {
	is := is.New(t)

	c = NewCatalog()
	r, err = c.AddRoot("s1")
	is.NotNil(r)
	is.NotErr(err)
	is.Equal(r.Path, "s1")

	r, err = c.AddRoot("s1")
	is.Nil(r)
	is.Err(err)
}

func TestRootAddFile(t *testing.T) {
	is := is.New(t)

	c = NewCatalog()
	r1, _ := c.AddRoot("r1")
	r2, _ := c.AddRoot("r2")

	r1.AddFile(f1)
	r1.AddFile(f2)
	r1.AddFile(f3)

	r2.AddFile(f2)
	r2.AddFile(f3)

	// Test that files ended up in files lists
	is.Equal(r1.Files, []*File{f1, f2, f3})
	is.Equal(r2.Files, []*File{f2, f3})

	// Test that files ended up in hash lists
	h := c.Hashes[95]
	is.Equal(h[0], RF{r1, f2})
	is.Equal(h[1], RF{r2, f2})

	h = c.Hashes[42]
	is.Equal(h[0], RF{r1, f1})
	is.Equal(h[1], RF{r1, f3})
	is.Equal(h[2], RF{r2, f3})
}

func TestRootRemoveFile(t *testing.T) {
	is := is.New(t)

	c = NewCatalog()
	r1, _ := c.AddRoot("r1")

	r1.AddFile(f1)
	r1.AddFile(f2)
	is.Equal(r1.Files, []*File{f1, f2})
	is.NotNil(c.Hashes[95])

	r1.RemoveFile(f1.Path)

	is.Equal(r1.Files, []*File{f2})
	is.Equal(c.Hashes[42], []RF{})
}
