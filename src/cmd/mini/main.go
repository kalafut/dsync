package main

import (
	"flag"
	"fmt"
)

func main() {
	flag.Parse()

	sourceA, sourceB := flag.Arg(0), flag.Arg(1)

	fmt.Printf("%v  %v", sourceA, sourceB)
}
