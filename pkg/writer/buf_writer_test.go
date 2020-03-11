package writer

import (
	"bytes"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrite(t *testing.T) {

	t.Run("Should write bytes to chan", func(t *testing.T) {
		// arrange
		text := "Lorem Ipsum  is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum."
		exp := []byte(text)
		c := make(chan []byte)
		go func() {
			for _, b := range exp {
				c <- []byte{b}
			}
			close(c)
		}()

		// act
		var buf bytes.Buffer
		w := NewWriter(&buf)
		var wg sync.WaitGroup
		w.Write(&wg, c)
		wg.Wait()

		// assert
		assert.Equal(t, exp, buf.Bytes())
	})
}
