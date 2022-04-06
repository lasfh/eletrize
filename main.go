package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("A [filename.json] argument was expected")
	}

	eletrize, err := NewEletrize(os.Args[1])
	if err != nil {
		log.Panicln(err)
	}

	eletrize.Start()
}
