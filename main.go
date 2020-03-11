package main

import (
	blockcompressor "blockcompressor/app"
	"fmt"
	"log"
	"os"
	// _ "net/http/pprof"
)

func main() {

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

	blockcompressor.Run(pathFile, destPathFile, 200000, 1, command)

	// go func() {
	// 	http.ListenAndServe("localhost:8080", nil)
	// }()

	fmt.Printf("done.")
}
