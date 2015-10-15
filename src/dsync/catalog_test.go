package dsync

import (
	"testing"

	"gopkg.in/tylerb/is.v1"
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

func TestAddRoot(t *testing.T) {
	is := is.New(t)

	c = NewCatalog()

	c.AddRoot("/path/a", "A")
	c.AddRoot("/another/path/b", "B")
	c.AddRoot("\\backslashed\\foo", "BS")

	is.True(contains("A", c.RootNames()))
	is.True(contains("B", c.RootNames()))
	is.False(contains("C", c.RootNames()))
	is.Equal(3, len(c.RootNames()))

	is.Equal(c.RootPath("A"), "/path/a")
	is.Equal(c.RootPath("B"), "/another/path/b")
	is.Equal(c.RootPath("BS"), "/backslashed/foo")
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

func TestUpdate(t *testing.T) {
	is := is.New(t)

	c = NewCatalog()
	c.AddRoot("./test_data", "test")
	c.UpdateRoot("test", false)

	testNames := []string{
		"test_data/file1",
		"test_data/file2",
		"test_data/folder1/file3",
		"test_data/folder1/folder2/file4",
	}

	for _, name := range testNames {
		found := false
		for _, f := range c.List() {
			if f.Path == name {
				found = true
				break
			}
		}
		is.True(found)
	}

	is.Equal(len(testNames), len(c.List()))

}
