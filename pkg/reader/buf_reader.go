package reader

import (
	"bufio"
	"fmt"
	"io"
	"log"
)

type readerByte struct {
	r         *bufio.Reader
	blockSize int
}

func NewBufReader(r *bufio.Reader, blockSize int) Reader {
	return readerByte{r: r, blockSize: blockSize}
}

func (r readerByte) Read() <-chan []byte {

	c := make(chan []byte)

	go func() {
		for {
			data := make([]byte, r.blockSize)
			n, err := r.r.Read(data)

			if err == io.EOF {
				fmt.Println("END read")
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
