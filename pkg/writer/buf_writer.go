package writer

import (
	"bufio"
	"log"
	"sync"
)

type bufWriter struct {
	w *bufio.Writer
}

func NewBufWriter(w *bufio.Writer) Writer {
	return bufWriter{w: w}
}

func (w bufWriter) Write(wg *sync.WaitGroup, c <-chan []byte) {

	wg.Add(1)
	go func() {
		defer wg.Done()
		for data := range c {
			_, err := w.w.Write(data)
			if err != nil {
				log.Fatal(err)
			}
		}
	}()
}
