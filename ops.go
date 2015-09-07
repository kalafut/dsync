package dsync

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
)

const CATALOG = "catalog.gob"

func Add(path, name string) {
	c, err := LoadCatalog(CATALOG)
	if err != nil {
		log.Fatal(err)
	}
	c.AddRoot(path, name)
	c.Save(CATALOG)
}

func CatalogStatus(filename string) {
	c, err := LoadCatalog(filename)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(c.RootNames())
}

func InitCatalog(filename string) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		c := NewCatalog()
		c.Save(filename)
	}
}

func SelectCatalog(config, filename string) {
	f, err := os.Create(config)
	defer f.Close()

	if err != nil {
		log.Fatal(err)
	}

	f.WriteString(filename)

}

func GetSelectedCatalog(config string) string {
	_f, err := os.Open(config)
	defer _f.Close()

	f := bufio.NewReader(_f)

	if err != nil {
		log.Fatal(err)
	}

	s, err := f.ReadString('\n')

	return s
}

func Build(root string, name string) {
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

	catalog.Save(name)
}
