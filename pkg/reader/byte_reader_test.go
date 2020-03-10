package reader

import (
	"bufio"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockReader struct {
	bufio.Reader
	ReadStringFunc func(delim byte) (string, error)
}

func (r *MockReader) ReadString(delim byte) (string, error) {
	if r.ReadStringFunc != nil {
		return r.ReadStringFunc(delim)
	}
	return r.Reader.ReadString(delim)
}

func TestRead(t *testing.T) {

	t.Run("Should read bytes to chan", func(t *testing.T) {
		// arrange
		blockSize := 10
		text := "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum."
		exp := []byte(text)

		sReader := strings.NewReader(text)
		b := bufio.NewReader(sReader)
		r := NewReader(b, blockSize)

		// act
		data := make([]byte, 0, len(exp))
		for s := range r.Read() {
			if len(s) > blockSize {
				assert.Fail(t, "asd")
			}
			data = append(data, s...)
		}

		// assert
		assert.Equal(t, exp, data)
	})
}
