package main

import (
	"os"

	"github.com/liampulles/proverb-gen/internal/wire"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	os.Exit(wire.Run(wd))
}
