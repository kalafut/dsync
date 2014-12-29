package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"os"
)

const BLOCK_SIZE = 8192

func g(path string) {
	buffer := make([]byte, BLOCK_SIZE)
	buf := new(bytes.Buffer)

	f, err := os.Open(path)
	fi, err := f.Stat()
	if err != nil {
		// Could not obtain stat, handle error
	}
	size := fi.Size()
	h := fnv.New64a()

	cnt, err := f.Read(buffer)
	h.Write(buffer)

	middle := size / 2

	f.Seek(middle, os.SEEK_SET)
	f.Read(buffer)
	h.Write(buffer)

	f.Seek(-BLOCK_SIZE, os.SEEK_END)
	f.Read(buffer)
	h.Write(buffer)

	binary.Write(buf, binary.LittleEndian, size)
	h.Write(buf.Bytes())
	fmt.Println(h.Sum64())

	_ = size
	_ = cnt
}

func main() {
	//g("fash.go")
	//g("/Volumes/video/Planes (2013)/Planes.2013.720p.BluRay.x264.YIFY.mp4")
	g("test2")
}
