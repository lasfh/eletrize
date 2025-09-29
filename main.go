package main

import (
	"context"
	"errors"
	"fmt"
	"os"
)

func main() {
	if err := execute(); err != nil && !errors.Is(err, context.Canceled) {
		fmt.Printf("eletrize:\n\t%s\n", err.Error())
		os.Exit(1)
	}
}
