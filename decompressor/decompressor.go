package decompressor

import (
	"blockconverter"
	"bufio"
	"bytes"
	"compress/gzip"
	"io"
	"log"
)

func NewCommand() blockconverter.Command {
	return decompressCommand{}
}

type decompressCommand struct{}

func (d decompressCommand) ReadTo(r io.Reader) <-chan blockconverter.Block {
	result := make(chan blockconverter.Block)
	go func() {
		for {
			b := blockconverter.Block{}
			err := b.InitFrom(r)
			if err == io.EOF {
				break
			}
			result <- b
		}
		close(result)
	}()
	return result
}

func (d decompressCommand) Convert(jobs <-chan blockconverter.Block) <-chan blockconverter.Block {
	result := make(chan blockconverter.Block)
	go func() {
		for job := range jobs {

			r := bytes.NewReader(job.Data)
			gz, err := gzip.NewReader(r)
			if err != nil {
				log.Fatal(err)
			}
			var buf bytes.Buffer
			w := bufio.NewWriter(&buf)
			_, err = io.Copy(w, gz)

			if err != nil {
				log.Fatal(err)
			}
			w.Flush()

			result <- blockconverter.Block{Index: job.Index, Data: buf.Bytes()}
		}
		close(result)
	}()
	return result
}

func (d decompressCommand) WriteTo(block <-chan blockconverter.Block, w io.Writer) <-chan struct{} {

	done := make(chan struct{})
	go func() {
		orderedBlocks := orderByIndex(block)

		for b := range orderedBlocks {
			_, err := w.Write(b)
			if err != nil {
				log.Fatal(err)
			}
		}
		close(done)
	}()

	return done
}

func orderByIndex(blocks <-chan blockconverter.Block) <-chan []byte {
	res := make(chan []byte)
	go func() {
		index := int64(0)
		m := make(map[int64][]byte)
		for block := range blocks {
			m[block.Index] = block.Data

			if v, ok := m[index]; ok {
				res <- v
				delete(m, index)

				index++
			}
		}
		for {
			if v, ok := m[index]; ok {
				res <- v
				delete(m, index)
				index++
			} else {
				break
			}
		}

		close(res)
	}()

	return res
}
