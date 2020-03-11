package command

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"io"
	"log"
)

const copyBlockSize = 1000000

type gzipCompress struct {
	fileName string
}

func NewGzipCompress(fileName string) Command {
	return gzipCompress{fileName}
}

func (g gzipCompress) Run(jobs <-chan Block) <-chan Block {
	result := make(chan Block)

	go func() {
		for block := range jobs {
			var buf bytes.Buffer
			w := gzip.NewWriter(&buf)
			w.Name = g.fileName
			
			for {
				_, err := io.CopyN(w, block.R, copyBlockSize)
				if err == io.EOF {
					break
				} else if err != nil {
					log.Fatal(err)
				}
			}
			w.Close()

			result <- Block{R: bytes.NewReader(buf.Bytes()), Indx: block.Indx}
		}

		close(result)
	}()

	return result
}

type gzipDecompress struct {
	fileName string
}

func NewGzipDecompress(fileName string) Command {
	return gzipDecompress{fileName}
}

func (g gzipDecompress) Run(jobs <-chan Block) <-chan Block {
	result := make(chan Block)

	go func() {
		for block := range jobs {
			r, err := gzip.NewReader(block.R)
			if err != nil {
				log.Fatal(err)
			}
			var buf bytes.Buffer
			w := bufio.NewWriter(&buf)

			for {
				_, err := io.CopyN(w, r, copyBlockSize)
				if err == io.EOF {
					break
				} else if err != nil {
					log.Fatal(err)
				}
			}

			for {
				_, err := io.CopyN(w, block.R, copyBlockSize)
				if err == io.EOF {
					break
				} else if err != nil {
					log.Fatal(err)
				}
			}

			result <- Block{R: bytes.NewReader(buf.Bytes()), Indx: block.Indx}
		}

		close(result)
	}()

	return result
}
