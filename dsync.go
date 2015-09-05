package main

import (
	"bufio"
	"log"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
)

const CATALOG = "catalog.gob"
const CONFIG = ".dsync"

func add(path, name string) {
	c, err := LoadCatalog(CATALOG)
	if err != nil {
		log.Fatal(err)
	}
	c.AddRoot(path, name)
	c.Save(CATALOG)
}

func initCatalog(filename string) {
	c := NewCatalog()
	c.Save(filename)
}

func selectCatalog(config, filename string) {
	f, err := os.Create(config)
	defer f.Close()

	if err != nil {
		log.Fatal(err)
	}

	f.WriteString(filename)

}

func getSelectedCatalog(config string) string {
	_f, err := os.Open(config)
	defer _f.Close()

	f := bufio.NewReader(_f)

	if err != nil {
		log.Fatal(err)
	}

	s, err := f.ReadString('\n')

	return s
}

func main() {
	var (
		app = kingpin.New("dsync", "Directory Synchronizer")

		addCmd   = app.Command("add", "Add a new root")
		add_path = addCmd.Arg("path", "Root path").Required().String()
		add_name = addCmd.Arg("name", "Shortname").String()

		buildCmd = app.Command("build", "Build catalog.")
		path     = buildCmd.Arg("path", "Root path").Required().String()
		name     = buildCmd.Arg("name", "Catalog name").Required().String()

		initCmd      = app.Command("init", "Create a new catalog.")
		initFilename = initCmd.Arg("filename", "Catalog filename").Required().String()

		list      = app.Command("list", "List catalog.")
		list_name = list.Arg("name", "Catalog name").Required().String()

		dedupeCmd = app.Command("dedupe", "Dedupe")
		keep      = dedupeCmd.Arg("keep name", "Keep catalog").Required().String()
		kill      = dedupeCmd.Arg("kill name", "Kill catalog").Required().String()
		doDelete  = dedupeCmd.Flag("delete", "Delete files").Bool()

		selectCmd      = app.Command("select", "Select, and create if necessary, the default catalog.")
		selectFilename = initCmd.Arg("filename", "Catalog filename").Required().String()
	)

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case addCmd.FullCommand():
		add(*add_path, *add_name)

	case buildCmd.FullCommand():
		build(*path, *name)

	// Post message
	case list.FullCommand():
		listCatalog(*list_name)

	case dedupeCmd.FullCommand():
		dedupe(*keep, *kill, *doDelete)

	case initCmd.FullCommand():
		//initCatalog(*initFilename)
		_ = initFilename
		initCatalog(CATALOG)

	case selectCmd.FullCommand():
		//initCatalog(*initFilename)
		_ = initFilename
		selectCatalog(CONFIG, *selectFilename)
	}

}
