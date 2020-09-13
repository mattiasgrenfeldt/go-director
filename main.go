package main

import (
	"fmt"
	"os"

	"github.com/mattiasgrenfeldt/go-director/director"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("%s <.dxr>\n", os.Args[0])
		os.Exit(0)
	}

	path := os.Args[1]
	fmt.Printf("%s\n", path)

	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	director.ParseDXR(f)
}
