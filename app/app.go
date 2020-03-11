package blockcompressor

import (
	"blockcompressor/pkg/command"
	"blockcompressor/pkg/reader"
	"blockcompressor/pkg/writer"
	"bytes"
	"sync"
)

type App struct {
	workers int
	r       reader.Reader
	c       command.Context
	w       writer.Writer
}

func NewApp(r reader.Reader, c command.Context, w writer.Writer, workers int) App {
	return App{r: r, c: c, w: w, workers: workers}
}

func (b App) Start(command string) {

	jobs := b.r.Read()

	results := runWorkers(jobs, b, command)

	var wg sync.WaitGroup
	b.w.Write(&wg, results)
	wg.Wait()
}

const startIndex = 1

func runWorkers(jobs <-chan *bytes.Reader, app App, instruction string) <-chan *bytes.Reader {
	markedBlocks := mark(jobs)

	strategy := app.c.GetStrategy(instruction)
	results := make([]<-chan command.Block, 0, app.workers)
	for i := 0; i < app.workers; i++ {
		worker := strategy.Run(markedBlocks)
		results = append(results, worker)
	}

	merged := merge(results...)

	res := order(merged)

	return res
}

func mark(jobs <-chan *bytes.Reader) <-chan command.Block {
	c := make(chan command.Block)
	go func() {
		i := startIndex
		for j := range jobs {
			c <- command.Block{R: j, Indx: i}
			i++
		}
		close(c)
	}()
	return c
}

func merge(chans ...<-chan command.Block) <-chan command.Block {
	merged := make(chan command.Block)
	var wg sync.WaitGroup
	wg.Add(len(chans))

	move := func(src <-chan command.Block) {
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

func order(merged <-chan command.Block) <-chan *bytes.Reader {
	res := make(chan *bytes.Reader)
	go func() {
		i := startIndex
		m := make(map[int]*bytes.Reader)
		for data := range merged {
			m[data.Indx] = data.R

			if v, ok := m[i]; ok {
				res <- v
				delete(m, i)
				i++
			}
		}
		for {
			if v, ok := m[i]; ok {
				res <- v
				delete(m, i)
				i++
			} else {
				break
			}
		}

		close(res)
	}()

	return res

}
