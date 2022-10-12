package main

import (
	"errors"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		eletrize, err := NewEletrizeByFileInCW()
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				fmt.Printf("eletrize: file %q in current directory: %s\n", validFileNames, err.Error())
				os.Exit(1)
			}

			panic(err)
		}

		eletrize.Start()
		os.Exit(0)
	}

	eletrize, err := NewEletrize(os.Args[1])
	if err != nil {
		panic(err)
	}

	eletrize.Start()
}
