package blockcompressor

import (
	"blockcompressor/pkg/compress"
	"blockcompressor/pkg/reader"
	"blockcompressor/pkg/writer"
	"sync"
)

type BlockCompressor struct {
	workers int
	r       reader.Reader
	c       compress.ContextCompress
	w       writer.Writer
}

func NewBlockCompressor(r reader.Reader, c compress.ContextCompress, w writer.Writer) BlockCompressor {
	return BlockCompressor{r: r, c: c, w: w}
}

func (b BlockCompressor) Run(algorithmCompress string) {
	jobs := b.r.Read()

	compressStrategy := b.c.GetCompressStrategy(algorithmCompress)
	results := make([]<-chan []byte, b.workers)

	for i := 0; i < b.workers; i++ {
		worker := compressStrategy.Run(jobs)
		results = append(results, worker)
	}

	mergedResults := merge(results)

	var wg sync.WaitGroup
	b.w.Write(&wg, mergedResults)
	wg.Wait()
}

func merge(chans []<-chan []byte) <-chan []byte {
	merged := make(chan []byte)
	var wg sync.WaitGroup
	wg.Add(len(chans))
	for _, c := range chans {
		go func() {
			defer wg.Done()
			for d := range c {
				merged <- d
			}
		}()
	}

	go func() {
		wg.Wait()
		close(merged)
	}()

	return merged
}
