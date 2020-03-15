package writer

import (
	"blockconverter"
	"blockconverter/reader"
	"bufio"
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrite(t *testing.T) {

	t.Run("Should write bytes to chan", func(t *testing.T) {

		// arrange
		text := "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum."
		// text = strings.Repeat(text, 10)
		textBytes := []byte(text)
		mExpected := make(map[int64]int64)
		c := make(chan blockconverter.Block)
		part := 5

		// send test data to chan
		go func() {
			index := int64(1)
			for i := 0; i < len(textBytes); i += part {
				right := i + part

				if right > len(textBytes) {
					right = len(textBytes)
				}
				b := textBytes[i:right]

				mExpected[index] = int64(len(b))

				c <- blockconverter.Block{B: b, Index: index}
				index++
			}
			close(c)
		}()

		// act
		var buf bytes.Buffer
		w := NewWriterF(&buf)
		done := w.Write(c)
		<-done

		// assert

		bi := reader.NewCompressedBlockInfo()
		r := bufio.NewReader(&buf)
		mActual := make(map[int64]int64)
		for {
			index, length, err := bi.Get(r)
			if err == io.EOF {
				break
			} else if err != nil {
				assert.Fail(t, "cannot create err expect io.EOF", err)
			}

			mActual[index] = length
			b := make([]byte, length)
			_, err = r.Read(b)
			if err != nil {
				assert.Fail(t, "")
			}
		}

		assert.Equal(t, mExpected, mActual)
	})
}
