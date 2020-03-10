package reader

import (
	"bufio"
	"log"
	"os"
)

func NewFileReader(pathFile string, blockSize int) Reader {
	file, err := os.Open(pathFile)
	if err != nil {
		log.Fatal(err)
	}

	r := bufio.NewReader(file)
	return NewBufReader(r, blockSize)
}
