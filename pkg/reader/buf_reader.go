package reader

import (
	"io"
	"log"
)

type readerByte struct {
	t         io.Reader
	blockSize int
}

func NewReader(r io.Reader, blockSize int) Reader {
	return readerByte{t: r, blockSize: blockSize}
}

func (r readerByte) Read() <-chan []byte {

	c := make(chan []byte)

	go func() {
		for {
			data := make([]byte, r.blockSize)
			n, err := r.t.Read(data)

			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(err)
			}

			if n < r.blockSize {
				data = data[:n]
			}

			c <- data
		}

		close(c)
	}()
	return c
}
