package reader

import (
	"blockconverter"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

type fileReader struct {
	file string
	bs   BlockInfo
}

func NewFileReader(file string, bs BlockInfo) blockconverter.Reader {
	return fileReader{file, bs}
}

func (f fileReader) Read() <-chan blockconverter.Block {
	file, err := os.Open(f.file)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	r := NewByteReader(&wg, file, f.bs)

	go func() {
		wg.Wait()
		file.Close()
	}()

	return r.Read()
}

type ioReader struct {
	r  io.Reader
	bs BlockInfo
	wg *sync.WaitGroup
}

func NewByteReader(wg *sync.WaitGroup, r io.Reader, bs BlockInfo) blockconverter.Reader {
	return ioReader{r, bs, wg}
}

func (b ioReader) Read() <-chan blockconverter.Block {
	c := make(chan blockconverter.Block)
	go func() {

		for {
			index, size, err := b.bs.Get(b.r)
			if err != nil {
				log.Fatal(err)
			}
			var buf bytes.Buffer
			w := bufio.NewWriter(&buf)
			_, err = io.CopyN(w, b.r, size)
			c <- blockconverter.Block{Index: index, B: buf.Bytes()}
			if err == io.EOF {
				break
			}

		}

		close(c)
		b.wg.Done()
	}()

	return c
}

type BlockInfo interface {
	Get(io.Reader) (index int64, size int64, err error)
}

const StartBlockIndex = int64(1)

type blockInfo struct {
	index        int64
	maxBlockSize int64
}

func NewBlockInfo(startBlockIndex int64, maxBlockSize int64) BlockInfo {
	return &blockInfo{
		index:        startBlockIndex,
		maxBlockSize: maxBlockSize,
	}
}

func (b *blockInfo) Get(r io.Reader) (index int64, size int64, err error) {
	index = b.index
	size = b.maxBlockSize
	b.index++
	return
}

type compressedBlockInfo struct{}

func NewCompressedBlockInfo() BlockInfo {
	return compressedBlockInfo{}
}

func (b compressedBlockInfo) Get(r io.Reader) (index int64, size int64, err error) {
	_, err = fmt.Fscanf(r, "%b %b", &index, &size)
	if err != nil {
		return 0, 0, err
	}

	return
}
