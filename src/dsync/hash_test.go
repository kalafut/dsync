package dsync

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/tylerb/is.v1"
)

var tempDir string

func TestMain(m *testing.M) {
	flag.Parse()

	tempDir, _ = ioutil.TempDir(os.TempDir(), "dsync_test_data")
	ret := m.Run()
	os.RemoveAll(tempDir)
	os.Exit(ret)
}

func TestHash(t *testing.T) {
	SAMPLE_FILE := filepath.Join(tempDir, "sample.txt")

	is := is.New(t)

	h, err := smartHash("not_found.txt")
	is.Err(err)

	WriteSample(SAMPLE_FILE, 100)
	h, err = smartHash(SAMPLE_FILE)
	os.Remove(SAMPLE_FILE)
	is.NotErr(err)

	is.Equal(h, uint64(0xf98540d8f8a71e22))

	WriteSample(SAMPLE_FILE, 10000000)
	h, err = smartHash(SAMPLE_FILE)
	os.Remove(SAMPLE_FILE)
	is.NotErr(err)

	is.Equal(h, uint64(0x93686806171fdb95))

	WriteSample(SAMPLE_FILE, 10000001)
	h, err = smartHash(SAMPLE_FILE)
	os.Remove(SAMPLE_FILE)
	is.NotErr(err)

	is.Equal(h, uint64(0x93686806171fdb95))
}

func WriteSample(name string, size int) {
	os.MkdirAll(filepath.Dir(name), 0666)
	data := make([]byte, size)

	for i := 0; i < size; i++ {
		data[i] = 'A'
	}

	err := ioutil.WriteFile(name, data, 0666)
	if err != nil {
		log.Fatal(err)
	}
}
