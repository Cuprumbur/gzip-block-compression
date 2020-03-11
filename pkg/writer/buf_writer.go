package writer

import (
	"io"
	"log"
	"sync"
)

type bufWriter struct {
	w io.Writer
}

func NewWriter(w io.Writer) Writer {
	return bufWriter{w: w}
}

func (b bufWriter) Write(wg *sync.WaitGroup, c <-chan []byte) {

	wg.Add(1)
	go func() {
		defer wg.Done()
		for bytes := range c {
			_, err := b.w.Write(bytes)
			if err != nil {
				log.Fatal(err)
			}
		}
	}()
}
