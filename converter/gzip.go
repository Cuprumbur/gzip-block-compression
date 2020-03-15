package converter

import (
	"blockconverter"
	"bufio"
	"bytes"
	"compress/gzip"
	"io"
	"log"
	"sync"
)

func NewGzipWorkers(workers int, c blockconverter.Converter) blockconverter.Converter {
	return gzipWorkers{workers, c}
}

type gzipWorkers struct {
	workers int
	c       blockconverter.Converter
}

func (g gzipWorkers) Run(jobs <-chan blockconverter.Block) <-chan blockconverter.Block {

	results := make([]<-chan blockconverter.Block, g.workers)

	for i := 0; i < g.workers; i++ {
		worker := g.c.Run(jobs)
		results[i] = worker
	}

	merged := merge(results...)

	return merged
}

func merge(chans ...<-chan blockconverter.Block) <-chan blockconverter.Block {
	merged := make(chan blockconverter.Block)
	var wg sync.WaitGroup
	wg.Add(len(chans))

	move := func(src <-chan blockconverter.Block) {
		for s := range src {
			merged <- s
		}
		wg.Done()
	}

	for _, c := range chans {
		go move(c)
	}

	go func() {
		wg.Wait()
		close(merged)
	}()

	return merged
}

type gzipCompress struct{}

func NewGzipCompress() blockconverter.Converter {
	return gzipCompress{}
}

func (g gzipCompress) Run(jobs <-chan blockconverter.Block) <-chan blockconverter.Block {
	result := make(chan blockconverter.Block)

	go func() {
		for block := range jobs {
			r := bytes.NewReader(block.B)
			var buf bytes.Buffer
			w := gzip.NewWriter(&buf)
			_, err := io.CopyN(w, r, r.Size())
			if err != nil {
				log.Fatal()
			}

			w.Close()
			result <- blockconverter.Block{Index: block.Index, B: buf.Bytes()}
		}

		close(result)
	}()

	return result
}

type gzipDecompress struct {
}

func NewGzipDecompress() blockconverter.Converter {
	return gzipDecompress{}
}

const sizePort = int64(10000)

func (g gzipDecompress) Run(jobs <-chan blockconverter.Block) <-chan blockconverter.Block {
	result := make(chan blockconverter.Block)

	go func() {
		for block := range jobs {
			br := bytes.NewReader(block.B)
			r, err := gzip.NewReader(br)
			if err != nil {
				log.Fatal(err)
			}

			var buf bytes.Buffer
			w := bufio.NewWriter(&buf)

			for {
				_, err := io.CopyN(w, r, sizePort)
				if err == io.EOF {
					break
				} else if err != nil {
					log.Fatal()
				}

			}
			w.Flush()

			result <- blockconverter.Block{Index: block.Index, B: buf.Bytes()}
		}

		close(result)
	}()

	return result
}
