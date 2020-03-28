package blockconverter

import (
	"bufio"
	"io"
	"log"
	"os"
	"sync"
)

type Reader interface {
	ReadTo(r io.Reader) <-chan Block
}
type Worker interface {
	Convert(<-chan Block) <-chan Block
}

type Writer interface {
	WriteTo(<-chan Block, io.Writer) (done <-chan struct{})
}

type Command interface {
	Reader
	Worker
	Writer
}

func Convert(c Command, filePathIn string, filePathOut string, workers int) {
	fileIn, err := os.Open(filePathIn)
	if err != nil {
		log.Fatal(err)
	}
	defer fileIn.Close()
	r := bufio.NewReader(fileIn)
	blocks := c.ReadTo(r)

	convertedBlocks := make([]<-chan Block, workers)

	for i := 0; i < workers; i++ {
		convertedBlocks[i] = c.Convert(blocks)
	}

	mergedConvertedBlocks := merge(convertedBlocks...)

	fileOut, err := os.Create(filePathOut)
	if err != nil {
		log.Fatal(err)
	}
	defer fileOut.Close()
	w := bufio.NewWriter(fileOut)
	<-c.WriteTo(mergedConvertedBlocks, w)
}

func merge(chans ...<-chan Block) <-chan Block {
	merged := make(chan Block)
	var wg sync.WaitGroup
	wg.Add(len(chans))

	move := func(src <-chan Block) {
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
