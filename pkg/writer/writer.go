package writer

import (
	"bytes"
	"io"
	"log"
	"sync"
)

type Writer interface {
	Write(wg *sync.WaitGroup, c <-chan *bytes.Reader)
}
type bufWriter struct {
	w io.Writer
}

const block = 100000

func NewWriter(w io.Writer) Writer {
	return bufWriter{w: w}
}

func (b bufWriter) Write(wg *sync.WaitGroup, c <-chan *bytes.Reader) {
	wg.Add(1)
	go func() {
		for r := range c {
			for {
				_, err := io.CopyN(b.w, r, block)

				if err == io.EOF {
					break
				} else if err != nil {
					log.Fatal(err)
				}
			}
		}
		wg.Done()
	}()
}
