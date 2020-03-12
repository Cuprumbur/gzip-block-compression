package reader

import (
	"blockcompressor/pkg/command"
	"bufio"
	"bytes"
	"io"
	"log"
)

type Reader interface {
	Read() (blocks <-chan *bytes.Reader)
}

type readerByte struct {
	r            io.Reader
	maxBlockSize int
}

func NewReader(r io.Reader, maxBlockSize int) Reader {
	return readerByte{r: r, maxBlockSize: maxBlockSize}
}

func (r readerByte) Read() <-chan *bytes.Reader {

	c := make(chan *bytes.Reader)
	maxBlockSize := int64(r.maxBlockSize)
	go func() {
		for {
			var buf bytes.Buffer
			w := bufio.NewWriter(&buf)
			_, err := io.CopyN(w, r.r, maxBlockSize)
			c <- bytes.NewReader(buf.Bytes())
			if err == io.EOF {
				break
			}
		}

		close(c)
	}()

	return c
}

type readerGzip struct {
	r            io.Reader
	maxBlockSize int
}

func NewReaderGzip(r io.Reader, maxBlockSize int) Reader {
	return readerGzip{r: r, maxBlockSize: maxBlockSize}
}

func (r readerGzip) Read() <-chan *bytes.Reader {

	c := make(chan *bytes.Reader)
	// maxBlockSize := int64(r.maxBlockSize)
	go func() {
		for {
			sizeB := make([]byte, 8)
			n, err := r.r.Read(sizeB)
			_ = n
			length := command.ToInt(sizeB)
			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(err)
			}
			log.Println(length)

			var buf bytes.Buffer
			w := bufio.NewWriter(&buf)
			t, err := io.CopyN(w, r.r, length)
			if t != length {
				log.Fatal("not equal")
			}

			c <- bytes.NewReader(buf.Bytes())
			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(err)
			}

		}

		close(c)
	}()

	return c
}
