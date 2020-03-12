package main

import (
	blockcompressor "blockcompressor/app"
	"log"
	"os"
	"time"
)

func main() {
	start := time.Now()

	command := os.Args[1]
	pathFile := os.Args[2]
	destPathFile := os.Args[3]

	if command == "" {
		log.Fatal("command not provided")
	}
	if pathFile == "" {
		log.Fatal("file not provided")
	}

	if destPathFile == "" {
		log.Fatal("dest path file not provided")
	}

	blockcompressor.Run(pathFile, destPathFile, 1000000, 1, command)

	elapsed := time.Since(start)
	log.Printf("elapsed %s", elapsed)
}
