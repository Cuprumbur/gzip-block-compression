package blockconverter

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBlockInit(t *testing.T) {
	text := "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum."

	t.Run("should init from io.Reader with restriction maxBlockSize", func(t *testing.T) {
		// arrange
		index := int64(3)
		maxBlockSize := int64(10)

		r := bytes.NewReader([]byte(text))
		b := Block{}
		// act
		err := b.Init(r, index, maxBlockSize)

		// assert
		length := int64(len(b.Data))
		assert.Nil(t, err)
		assert.Equal(t, index, b.Index)
		assert.Equal(t, maxBlockSize, length)
	})

	t.Run("should init from io.Reader that provided index and size info", func(t *testing.T) {
		payload := []byte(text)

		index := int64(5)
		blockSize := int64(len(payload))

		dataSavedInBlock := make([]byte, 0)
		dataSavedInBlock = append(dataSavedInBlock, toByte(index)...)
		dataSavedInBlock = append(dataSavedInBlock, toByte(blockSize)...)
		dataSavedInBlock = append(dataSavedInBlock, payload...)

		r := bytes.NewReader(dataSavedInBlock)

		// act
		b := Block{}
		err := b.InitFrom(r)

		// assert
		length := int64(len(b.Data))
		assert.Nil(t, err)
		assert.Equal(t, index, b.Index)
		assert.Equal(t, blockSize, length)

	})

	t.Run("should init payload from io.Reader", func(t *testing.T) {
		payload := []byte(text)

		index := int64(5)
		blockSize := int64(len(payload))

		dataSavedInBlock := make([]byte, 0)
		dataSavedInBlock = append(dataSavedInBlock, toByte(index)...)
		dataSavedInBlock = append(dataSavedInBlock, toByte(blockSize)...)
		dataSavedInBlock = append(dataSavedInBlock, payload...)

		r := bytes.NewReader(dataSavedInBlock)

		// act
		b := Block{}
		err := b.InitFrom(r)

		// assert
		assert.Nil(t, err)
		assert.Equal(t, payload, b.Data)

	})

}
