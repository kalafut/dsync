package main

import (
	"testing"
	"time"

	"gopkg.in/tylerb/is.v1"
)

const TEST_DATA = "./test_data"

var (
	c   *Catalog
	r   *Root
	err error
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

	f1 := &File{Path: "f1", Size: 50, Hash: 42, ModTime: time.Now()}
	f2 := &File{Path: "f2", Size: 150, Hash: 95, ModTime: time.Now()}
	f3 := &File{Path: "f3", Size: 50, Hash: 42, ModTime: time.Now()}

	c = NewCatalog()
	r1, _ := c.AddRoot("r1")
	r2, _ := c.AddRoot("r2")

	r1.AddFile(f1)
	r1.AddFile(f2)
	r1.AddFile(f3)

	r2.AddFile(f2)
	r2.AddFile(f3)

	is.Equal(r1.Files, []*File{f1, f2, f3})
	is.Equal(r2.Files, []*File{f2, f3})

	h := c.Hashes[95]
	hf := h[0]
	is.Equal(hf.Root, r1)
	is.Equal(hf.File, f2)

hf=h[1]
is.Equal(hf.Root, r2)
	is.Equal(hf.File, f2)

h = c.Hashes[42]
	hf = h[0]
	is.Equal(hf.Root, r1)
	is.Equal(hf.File, f1)

hf=h[1]
is.Equal(hf.Root, r1)
	is.Equal(hf.File, f3)

	hf = h[2]
	is.Equal(hf.Root, r2)
	is.Equal(hf.File, f3)


}
