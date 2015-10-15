package main

import (
	"os"

	"dsync"

	"gopkg.in/alecthomas/kingpin.v2"
)

const CONFIG = ".dsync"

func main() {
	var (
		app = kingpin.New("dsync", "Directory Synchronizer")

		addCmd   = app.Command("add", "Add a new root")
		add_path = addCmd.Arg("path", "Root path").Required().String()
		add_name = addCmd.Arg("name", "Shortname").String()

		stCmd      = app.Command("stat", "Catalog status")
		stFilename = stCmd.Arg("name", "Catalog name").Required().String()

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

		selectCmd      = app.Command("select", "Select the default catalog.")
		selectFilename = selectCmd.Arg("filename", "Catalog filename").Required().String()
	)

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case addCmd.FullCommand():
		dsync.Add(*add_path, *add_name)

	case buildCmd.FullCommand():
		dsync.Build(*path, *name)

	case list.FullCommand():
		dsync.ListCatalog(*list_name)

	case dedupeCmd.FullCommand():
		dsync.Dedupe(*keep, *kill, *doDelete)

	case initCmd.FullCommand():
		//initCatalog(*initFilename)
		_ = initFilename
		dsync.InitCatalog(dsync.CATALOG)

	case selectCmd.FullCommand():
		dsync.SelectCatalog(CONFIG, *selectFilename)

	case stCmd.FullCommand():
		dsync.CatalogStatus(*stFilename)
	}

}
