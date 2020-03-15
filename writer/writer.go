package writer

import (
	"blockconverter"
	"blockconverter/reader"
	"bytes"
	"io"
	"log"
	"os"
)

type fileWriter struct {
	file string
	c    string
}

func NewFileWriter(file string, c string) blockconverter.Writer {
	return fileWriter{file, c}
}

func (b fileWriter) Write(c <-chan blockconverter.Block) <-chan struct{} {
	file, err := os.Create(b.file)
	if err != nil {
		log.Fatal(err)
	}
	done := make(chan struct{})
	w := NewWriterOrder(file, b.c)

	go func() {
		<-w.Write(c)

		file.Close()

		done <- struct{}{}
	}()

	return done
}

type ff struct {
	w io.Writer
	c string
}

func NewWriterOrder(w io.Writer, c string) blockconverter.Writer {
	return ff{w, c}
}

func (f ff) Write(blocks <-chan blockconverter.Block) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		if f.c == "d" {
			b := orderByIndex(blocks)
			for i := range b {

				r := bytes.NewReader(i)

				n, err := r.WriteTo(f.w)
				_ = n
				if err != nil {
					log.Fatal(err)
				}
			}

		}
		for b := range blocks {

			if f.c == "c" {
				r := bytes.NewReader(b.B)
				err := blockconverter.SetData(f.w, b.Index, r.Size())
				if err != nil {
					log.Fatal(err)
				}

				n, err := r.WriteTo(f.w)
				_ = n
				if err != nil {
					log.Fatal(err)
				}
			}

		}
		close(done)
	}()

	return done
}

func orderByIndex(blocks <-chan blockconverter.Block) <-chan []byte {
	res := make(chan []byte)
	go func() {
		i := reader.StartBlockIndex
		m := make(map[int64][]byte)
		for block := range blocks {
			m[block.Index] = block.B

			if v, ok := m[i]; ok {
				res <- v
				delete(m, i)

				i++
			}
		}
		for {
			if v, ok := m[i]; ok {
				res <- v
				delete(m, i)
				i++
			} else {
				break
			}
		}

		close(res)
	}()

	return res
}
