package writer

import (
	"bufio"
	"log"
	"os"
)

func NewFileWriter(pathFile string) Writer {
	file, err := os.Create(pathFile)
	if err != nil {
		log.Fatal(err)
	}
	w := bufio.NewWriter(file)
	return NewBufWriter(w)
}
