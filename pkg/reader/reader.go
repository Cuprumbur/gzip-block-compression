package reader

import (
	"bufio"
	"bytes"
	"io"
)

type Reader interface {
	Read() (blocks <-chan *bytes.Reader)
}

type readerByte struct {
	r         io.Reader
	blockSize int
}

func NewReader(r io.Reader, blockSize int) Reader {
	return readerByte{r: r, blockSize: blockSize}
}

func (r readerByte) Read() <-chan *bytes.Reader {

	c := make(chan *bytes.Reader)
	size := int64(r.blockSize)
	go func() {
		for {
			var buf bytes.Buffer
			w := bufio.NewWriter(&buf)
			_, err := io.CopyN(w, r.r, size)

			c <- bytes.NewReader(buf.Bytes())
			if err == io.EOF {
				break
			}
		}

		close(c)
	}()

	return c
}
