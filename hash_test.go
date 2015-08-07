package main

import (
	"io/ioutil"
	"os"
	"testing"

	"gopkg.in/tylerb/is.v1"
)

func TestHash(t *testing.T) {
	is := is.New(t)

	h, err := smartHash("not_found.txt")
	is.Err(err)

	WriteSample("sample.txt", 100)
	h, err = smartHash("sample.txt")
	os.Remove("sample.txt")
	is.NotErr(err)

	is.Equal(h, uint64(10489554043305584505))

	WriteSample("sample.txt", 10000000)
	h, err = smartHash("sample.txt")
	os.Remove("sample.txt")
	is.NotErr(err)

	is.Equal(h, uint64(4524196217030972197))

	WriteSample("sample.txt", 10000001)
	h, err = smartHash("sample.txt")
	os.Remove("sample.txt")
	is.NotErr(err)

	is.Equal(h, uint64(4524196217030972197))
}

func WriteSample(name string, size int) {
	data := make([]byte, size)

	for i := 0; i < size; i++ {
		data[i] = 'A'
	}

	ioutil.WriteFile(name, data, 0666)
}
