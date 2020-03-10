package main

import (
	"fmt"
	"log"
	"os"
	// _ "net/http/pprof"
)

func main() {
	command := os.Args[0]
	pathFile := os.Args[1]
	destPathFile := os.Args[2]

	if command == "" {
		log.Fatal("command not provided")
	}
	if pathFile == "" {
		log.Fatal("file not provided")
	}

	if destPathFile == "" {
		log.Fatal("dest path file not provided")
	}

	// go func() {
	// 	http.ListenAndServe("localhost:8080", nil)
	// }()

	fmt.Printf("done.")
}
