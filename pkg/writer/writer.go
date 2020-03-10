package writer

import (
	"log"
	"os"
	"sync"
)

type Writer interface {
	Write(wg *sync.WaitGroup, c <-chan []byte)
}

type fileWriter struct {
	pathFile string
}

func NewFileWriter(pathFile string) Writer {
	return fileWriter{pathFile: pathFile}
}

func (w fileWriter) Write(wg *sync.WaitGroup, c <-chan []byte) {

	file, err := os.Create(w.pathFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for data := range c {
			_, err := file.Write(data)
			if err != nil {
				log.Fatal(err)
			}
		}
	}()
}
