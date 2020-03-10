package reader

import (
	"bufio"
	"log"
	"os"
)

type fileReader struct {
	pathFile  string
	blockSize int
}

func NewFileReader(pathFile string, blockSize int) Reader {
	fileReader, err := os.Open(pathFile)
	if err != nil {
		log.Fatal(err)
	}

	r := bufio.NewReader(fileReader)
	return NewReader(r, blockSize)
}
