package blockcompressor

import (
	"blockcompressor/pkg/compress"
	"blockcompressor/pkg/reader"
	"blockcompressor/pkg/writer"
	"sync"
)

type App struct {
	workers int
	r       reader.Reader
	c       compress.ContextCompress
	w       writer.Writer
}

func NewApp(r reader.Reader, c compress.ContextCompress, w writer.Writer, workers int) App {
	return App{r: r, c: c, w: w, workers: workers}
}

func (b App) Run(algorithm string) {

	jobs := b.r.Read()

	compressStrategy := b.c.GetCompressStrategy(algorithm)
	results := make([]<-chan []byte, 0, b.workers)

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
		go func(ch <-chan []byte) {
			defer wg.Done()
			for d := range ch {
				merged <- d
			}
		}(c)
	}

	go func() {
		wg.Wait()
		close(merged)
	}()

	return merged
}
