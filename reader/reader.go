package reader

import (
	"blockconverter"
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"os"
	"sync"
)

type fileReader struct {
	file string
	bi   BlockInfo
}

func NewFileReader(file string, bi BlockInfo) blockconverter.Reader {
	return fileReader{file, bi}
}

func (f fileReader) Read() <-chan blockconverter.Block {
	file, err := os.Open(f.file)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	r := NewReader(&wg, file, f.bi)

	go func() {
		wg.Wait()
		file.Close()
	}()

	return r.Read()
}

type ioReader struct {
	r  io.Reader
	bi BlockInfo
	wg *sync.WaitGroup
}

func NewReader(wg *sync.WaitGroup, r io.Reader, bi BlockInfo) blockconverter.Reader {
	return ioReader{r, bi, wg}
}

func (b ioReader) Read() <-chan blockconverter.Block {
	c := make(chan blockconverter.Block)
	go func() {

		for {
			index, size, err := b.bi.Get(b.r)
			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(err)
			}
			var buf bytes.Buffer
			w := bufio.NewWriter(&buf)
			_, err = io.CopyN(w, b.r, size)
			w.Flush()
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
	i := make([]byte, 8)
	_, err = r.Read(i)
	if err != nil {
		return 0, 0, err
	}

	index = ToInt(i)

	s := make([]byte, 8)
	_, err = r.Read(s)
	if err != nil {
		return 0, 0, err
	}

	size = ToInt(s)
	if err != nil {
		return 0, 0, err
	}

	return
}

func ToInt(b []byte) int64 {
	i := int64(binary.BigEndian.Uint64(b))
	return i
}
