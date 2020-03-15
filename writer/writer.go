package writer

import (
	"blockconverter"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

type fileWriter struct {
	file string
}

func NewFileWriter(file string) blockconverter.Writer {
	return fileWriter{file}
}

func (b fileWriter) Write(c <-chan blockconverter.Block) <-chan struct{} {
	file, err := os.Create(b.file)
	if err != nil {
		log.Fatal(err)
	}
	done := make(chan struct{})
	w := NewWriterF(file)

	go func() {
		<-w.Write(c)

		file.Close()

		done <- struct{}{}
	}()

	return done
}

type ff struct {
	w io.Writer
}

func NewWriterF(w io.Writer) blockconverter.Writer {
	return ff{w}
}

func SetData(r io.Writer, index int64, size int64) error {

	_, err := fmt.Fprintf(r, "%b %b", index, size)
	if err != nil {
		return err
	}
	return nil
}

func (f ff) Write(blocks <-chan blockconverter.Block) <-chan struct{} {
	done := make(chan struct{})
	go func() {

		for b := range blocks {
			r := bytes.NewReader(b.B)
			bw := bufio.NewWriter(f.w)

			err := SetData(bw, b.Index, r.Size())
			if err != nil {
				log.Fatal(err)
			}
			bw.Flush()

			n, err := r.WriteTo(bw)
			_ = n
			if err != nil {
				log.Fatal(err)
			}
			bw.Flush()
		}
		close(done)
	}()

	return done
}

// type decWriter struct {
// 	w io.Writer
// }

// func NewDecWriter(w io.Writer) Writer {
// 	return decWriter{w: w}
// }

// func (b decWriter) Write(wg *sync.WaitGroup, c <-chan *bytes.Reader) {
// 	wg.Add(1)
// 	go func() {
// 		for r := range c {
// 			n := r.Size()
// 			s, err := io.CopyN(b.w, r, n)
// 			if n != s {
// 				log.Fatal("not equal")
// 			}

// 			if err == io.EOF {
// 				break
// 			} else if err != nil {
// 				log.Fatal(err)
// 			}
// 		}
// 		wg.Done()
// 	}()
// }
