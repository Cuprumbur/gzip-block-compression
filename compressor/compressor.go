package compressor

import (
	"blockconverter"
	"bytes"
	"compress/gzip"
	"io"
	"log"
)

func NewCommand(maximumBlockReadSize int64) blockconverter.Command {
	return compressCommand{maximumBlockReadSize}
}

type compressCommand struct {
	maximumBlockReadSize int64
}

func (c compressCommand) ReadTo(r io.Reader) <-chan blockconverter.Block {
	result := make(chan blockconverter.Block)
	go func() {
		index := int64(0)
		for {
			b := blockconverter.Block{}
			err := b.Init(r, index, c.maximumBlockReadSize)
			result <- b
			if err == io.EOF {
				break
			}
			index++
		}
		close(result)
	}()
	return result
}

func (c compressCommand) Convert(jobs <-chan blockconverter.Block) <-chan blockconverter.Block {
	result := make(chan blockconverter.Block)
	go func() {
		for job := range jobs {
			var buf bytes.Buffer
			w := gzip.NewWriter(&buf)
			r := bytes.NewReader(job.Data)
			_, err := io.CopyN(w, r, r.Size())

			if err != nil {
				log.Fatal(err)
			}
			w.Close()

			result <- blockconverter.Block{Index: job.Index, Data: buf.Bytes()}
		}
		close(result)
	}()
	return result
}

func (c compressCommand) WriteTo(block <-chan blockconverter.Block, w io.Writer) <-chan struct{} {

	done := make(chan struct{})
	go func() {
		for b := range block {
			err := b.WriteToWithBlockInfo(w)
			if err != nil {
				log.Fatal(err)
			}
		}
		close(done)
	}()

	return done
}
