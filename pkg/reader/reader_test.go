package reader

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRead(t *testing.T) {

	t.Run("Should read bytes to chan", func(t *testing.T) {
		// arrange
		blockSize := 3
		text := "Lorem Ipsum is simply dummy  " //text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum."
		exp := []byte(text)

		b := bytes.NewReader(exp)
		r := NewReader(b, blockSize)

		// act
		data := make([]byte, 0, len(exp))
		for s := range r.Read() {
			b := make([]byte, blockSize)
			n, err := s.Read(b)
			if err == io.EOF {
				assert.Fail(t, "EOF")
			}

			data = append(data, b[:n]...)
		}

		// assert
		assert.Equal(t, exp, data)
	})
}
