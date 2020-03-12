package writer

import (
	"blockcompressor/pkg/command"
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

func NewWriter(w io.Writer) Writer {
	return bufWriter{w: w}
}

func (b bufWriter) Write(wg *sync.WaitGroup, c <-chan *bytes.Reader) {
	wg.Add(1)
	go func() {

		for r := range c {
			s := r.Size()
			length := command.ToByte(s)
			log.Println(s)
			n, err := b.w.Write(length)
			_ = n
			if err != nil {
				log.Fatal(err)
			}

			t, err := io.CopyN(b.w, r, s)
			if t != s {
				log.Fatal("not equal")
			}
			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(err)
			}
		}
		wg.Done()
	}()
}

type decWriter struct {
	w io.Writer
}

func NewDecWriter(w io.Writer) Writer {
	return decWriter{w: w}
}

func (b decWriter) Write(wg *sync.WaitGroup, c <-chan *bytes.Reader) {
	wg.Add(1)
	go func() {
		for r := range c {
			n := r.Size()
			s, err := io.CopyN(b.w, r, n)
			if n != s {
				log.Fatal("not equal")
			}

			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(err)
			}
		}
		wg.Done()
	}()
}
