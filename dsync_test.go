package main

import (
	"reflect"
	"testing"

	"gopkg.in/tylerb/is.v1"
)

const TEST_DATA = "./test_data"

var (
	c   *Catalog
	err error
	f1  = &File{Path: "f1", Size: 50, Hash: 42}
	f2  = &File{Path: "f2", Size: 150, Hash: 95}
	f3  = &File{Path: "f3", Size: 50, Hash: 42}
	f4  = &File{Path: "f4", Size: 150, Hash: 95}
)

func TestAddFile(t *testing.T) {
	is := is.New(t)

	c = NewCatalog()

	c.AddFile(f1)
	c.AddFile(f2)
	c.AddFile(f3)

	// Test that files ended up in hash lists
	is.Equal(c.Hashes[42], []*File{f1, f3})
	is.Equal(c.Hashes[95], []*File{f2})
	is.Nil(c.Hashes[999])

	// Test Removal
	c.RemoveFile(f1)
	is.Equal(c.Hashes[42], []*File{f3})
	c.RemoveFile(f3)
	is.Equal(c.Hashes[42], []*File{})

	is.Equal(c.Hashes[95], []*File{f2})
	c.RemoveFile(f2)
	is.Equal(c.Hashes[95], []*File{})
}

func TestList(t *testing.T) {
	is := is.New(t)

	c = NewCatalog()

	c.AddFile(f1)
	c.AddFile(f2)
	c.AddFile(f3)

	l := c.List()
	is.True(contains(f1, l))
	is.True(contains(f2, l))
	is.True(contains(f3, l))
}

func TestDupes(t *testing.T) {
	is := is.New(t)

	c = NewCatalog()

	c.AddFile(f1)
	c.AddFile(f2)
	c.AddFile(f3)
	c.AddFile(f4)

	dupes := c.Dupes()

	is.Equal(dupes[0], []*File{f1, f3})
	is.Equal(dupes[1], []*File{f2, f4})
}

func TestAddRoot(t *testing.T) {
	is := is.New(t)

	c = NewCatalog()

	c.AddRoot("/path/a", "A")
	c.AddRoot("/another/path/b", "B")

	is.True(contains("A", c.RootNames()))
	is.True(contains("B", c.RootNames()))
	is.False(contains("C", c.RootNames()))
	is.Equal(2, len(c.RootNames()))

	is.Equal(c.GetPath("A"), "/path/a")
	is.Equal(c.GetPath("B"), "/another/path/b")
}

func contains(elem interface{}, list interface{}) bool {
	v := reflect.ValueOf(list)
	for i := 0; i < v.Len(); i++ {
		if v.Index(i).Interface() == elem {
			return true
		}
	}
	return false
}

/*
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
*/
