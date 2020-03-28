package main

import (
	"blockconverter"
	"blockconverter/compressor"
	"blockconverter/decompressor"
	"log"
	"os"
	"time"
)

func main() {

	if len(os.Args) != 4 {
		log.Fatal("should provided 3 arguments")
	}
	command := os.Args[1]
	pathFileIn := os.Args[2]
	pathFileOut := os.Args[3]

	if command == "" {
		log.Fatal("command not provided")
	}
	if pathFileIn == "" {
		log.Fatal("file not provided")
	}

	if pathFileOut == "" {
		log.Fatal("dest path file not provided")
	}

	workers := 10

	maximumBlockReadSize := int64(1000000)

	m := map[string]blockconverter.Command{
		"compress":   compressor.NewCommand(maximumBlockReadSize),
		"decompress": decompressor.NewCommand(),
	}

	start := time.Now()
	blockconverter.Convert(m[command], pathFileIn, pathFileOut, workers)
	elapsed := time.Since(start)
	log.Printf("elapsed %s", elapsed)
}
