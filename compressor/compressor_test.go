package compressor

import (
	"blockconverter"
	"bufio"
	"bytes"
	"compress/gzip"
	"testing"

	"github.com/stretchr/testify/assert"
)

const text = "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum."

func TestReadTo(t *testing.T) {

	t.Run("should read to the channel the whole data", func(t *testing.T) {
		// arrange
		blockSize := int64(5)
		exp := []byte(text)

		r := bytes.NewReader(exp)

		c := NewCommand(blockSize)

		// act
		data := make([]byte, 0, len(exp))
		for s := range c.ReadTo(r) {
			data = append(data, s.Data...)
		}

		// assert
		assert.Equal(t, exp, data)
	})
}

func TestCompress(t *testing.T) {
	t.Run("source data should be greater than compressed data", func(t *testing.T) {

		// assert
		data := []byte(text)
		index := int64(3)

		block := blockconverter.Block{
			Index: index,
			Data:  data,
		}

		job := make(chan blockconverter.Block)
		go func() {
			job <- block
			close(job)
		}()
		c := NewCommand(0)

		// act
		result := c.Convert(job)
		d := <-result

		// arrange
		assert.Equal(t, index, block.Index)
		assert.Greater(t, len(data), len(d.Data))
	})

	t.Run("should has valid gzip header after compress", func(t *testing.T) {

		// arrange
		index := int64(1)

		job := make(chan blockconverter.Block)
		go func() {
			job <- blockconverter.Block{Index: index, Data: []byte(text)}
			close(job)
		}()

		c := NewCommand(0)

		// act
		result := c.Convert(job)

		block := <-result

		// assert
		assert.Equal(t, index, block.Index)

		br := bytes.NewReader(block.Data)
		_, err := gzip.NewReader(br)
		assert.Nil(t, err)
	})

}

func TestWrite(t *testing.T) {
	//	arrange
	job := make(chan blockconverter.Block)
	b := []byte(text)
	sp := make([][]byte, 0)
	lengthPart := 4
	right := 0
	for i := 0; i < len(b); i += lengthPart {
		if i+lengthPart < len(b) {
			right = i + lengthPart
		} else {
			right = len(b)
		}
		sp = append(sp, b[i:right])
	}
	size := 0
	go func() {
		for _, data := range sp {
			b := blockconverter.Block{Index: 1, Data: data}
			size += 8 // index
			size += 8 // length of block
			size += len(b.Data)
			job <- b
		}
		close(job)
	}()

	// act
	c := NewCommand(0)
	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)

	done := c.WriteTo(job, w)
	<-done
	w.Flush()

	assert.Equal(t, size, len(buf.Bytes()))

}
