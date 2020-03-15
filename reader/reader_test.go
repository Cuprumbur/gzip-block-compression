package reader

import (
	"blockconverter"
	"blockconverter/mocks"
	"bytes"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func BenchmarkRead(b *testing.B) {
	blockSize := int64(1000000)
	text := "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum."
	text = strings.Repeat(text, 100)
	exp := []byte(text)
	m := mocks.BlockInfo{}
	m.On("Get", mock.Anything).Return(int64(0), blockSize, nil)

	for i := 0; i < b.N; i++ {
		// arrange

		br := bytes.NewReader(exp)

		var wg sync.WaitGroup
		wg.Add(1)
		r := NewReader(&wg, br, &m)

		// act

		data := make([]byte, 0, len(exp))
		for s := range r.Read() {
			data = append(data, s.B...)
		}

		// assert
		assert.Equal(b, exp, data)
	}
}
func TestRead(t *testing.T) {

	t.Run("Should read bytes to chan", func(t *testing.T) {
		// arrange

		blockSize := int64(3)
		text := "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum."
		exp := []byte(text)
		m := mocks.BlockInfo{}
		m.On("Get", mock.Anything).Return(int64(0), blockSize, nil)
		b := bytes.NewReader(exp)

		var wg sync.WaitGroup
		wg.Add(1)
		r := NewReader(&wg, b, &m)

		// act
		data := make([]byte, 0, len(exp))
		for s := range r.Read() {
			data = append(data, s.B...)
		}

		// assert
		assert.Equal(t, exp, data)
	})
}
func TestBlockInfo(t *testing.T) {

	t.Run("Should read block info from io.Reader", func(t *testing.T) {
		// arrange
		index := int64(5)
		blockSize := int64(3)
		var buf bytes.Buffer

		err := blockconverter.SetData(&buf, index, blockSize)
		if err != nil {
			t.Fatal(err)
		}

		c := NewCompressedBlockInfo()
		i, s, err := c.Get(&buf)
		assert.Nil(t, err)
		assert.Equal(t, index, i)
		assert.Equal(t, blockSize, s)

	})

	t.Run("Should read block info that provided in ctor", func(t *testing.T) {
		// arrange
		startIndex := StartBlockIndex
		maxBlockSize := int64(20000)
		var buf bytes.Buffer

		c := NewBlockInfo(startIndex, maxBlockSize)
		i, s, err := c.Get(&buf)
		assert.Nil(t, err)
		assert.Equal(t, startIndex, i)
		assert.Equal(t, maxBlockSize, s)

	})
}
