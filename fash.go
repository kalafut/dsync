package main

import (
	"fmt"
	"hash/fnv"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/davecheney/profile"
	"github.com/kalafut/gosh"
	"github.com/spf13/cobra"
)

const SAMPLE_SIZE = 4096
const WORKERS = 10

var exts = gosh.NewSet(".jpg", ".mp4")

func hashFiles(files <-chan *File, catalog *Catalog) {
	buffer := make([]byte, SAMPLE_SIZE)
	h := fnv.New64a()
	for file := range files {
		h.Reset()
		f, _ := os.Open(file.Path)
		f.Read(buffer)
		h.Write(buffer)
		f.Close()

		file.Hash = h.Sum64()
		catalog.AddFile(file)
	}
}

func traverse(root string) <-chan *File {
	files := make(chan *File)
	go func() {
		defer close(files)
		filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				panic(err)
			}
			if exts.Contains(strings.ToLower(filepath.Ext(path))) {
				files <- &File{Path: path, Size: info.Size(), ModTime: info.ModTime()}
			}
			return nil
		})
	}()
	return files
}

func monitor(catalog *Catalog) {
	for {
		fmt.Println(len(catalog.Files))
		time.Sleep(1 * time.Second)
	}
}

func cli() {
	var HugoCmd = &cobra.Command{
		Use:   "hugo",
		Short: "Hugo is a very fast static site generator",
		Long: `A Fast and Flexible Static Site Generator built with
					            love by spf13 and friends in Go.
								            Complete documentation is available at http://hugo.spf13.com`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("here!!!")
			// Do Stuff Here
		},
	}
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of Hugo",
		Long:  `All software has versions. This is Hugo's`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Hugo Static Site Generator v0.9 -- HEAD")
		},
	}
	HugoCmd.AddCommand(versionCmd)

	var buildCmd = &cobra.Command{
		Use:   "build [root]",
		Short: "Print the version number of Hugo",
		Long:  `All software has versions. This is Hugo's`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				return
			}
			build(args[0])
		},
	}
	HugoCmd.AddCommand(buildCmd)

	var listCmd = &cobra.Command{
		Use:   "catalog",
		Short: "Print the version number of Hugo",
		Long:  `All software has versions. This is Hugo's`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := listCatalog(); err != nil {
				fmt.Println(err)
			}
		},
	}
	HugoCmd.AddCommand(listCmd)

	HugoCmd.Execute()
}

func listCatalog() error {
	catalog, err := LoadCatalog("c.gob")
	if err != nil {
		return err
	}
	catalog.List()

	return nil
}

func build(root string) {
	var wg sync.WaitGroup
	catalog := NewCatalog()

	go monitor(catalog)

	files := traverse(root)

	for w := 0; w < WORKERS; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			hashFiles(files, catalog)
		}()
	}
	wg.Wait()

	catalog.Save("c.gob")
}

func main() {
	cfg := profile.Config{
		MemProfile:     true,
		CPUProfile:     true,
		BlockProfile:   true,
		ProfilePath:    ".",  // store profiles in current directory
		NoShutdownHook: true, // do not hook SIGINT
	}

	// p.Stop() must be called before the program exits to
	// ensure profiling information is written to disk.
	//p := profile.Start(&cfg)
	//defer p.Stop()
	_ = cfg

	cli()
}
