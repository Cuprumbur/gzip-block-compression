package writer

import "sync"

type Writer interface {
	Write(wg *sync.WaitGroup, c <-chan []byte)
}
