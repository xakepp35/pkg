package main

import (
	"fmt"
	"os"

	"google.golang.org/protobuf/compiler/protogen"
)

func main() {
	if len(os.Args) == 2 && os.Args[1] == "--version" {
		fmt.Println("version:", version)
		return
	}

	if len(os.Args) == 2 && (os.Args[1] == "--help" || os.Args[1] == "-h") {
		flags.Usage()
		return
	}

	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(generateFile)
}
