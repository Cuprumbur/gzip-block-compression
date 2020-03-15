package converter

import (
	"blockconverter"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGzipCompress(t *testing.T) {
	t.Run("compress data", func(t *testing.T) {
		gc := NewGzipCompress()
		c := make(chan blockconverter.Block)
		index := int64(7)
		text := "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum."

		b := []byte(text)

		go func() {

			c <- blockconverter.Block{Index: index, B: b}
		}()

		result := gc.Run(c)

		block := <-result
		assert.Equal(t, index, block.Index)
		assert.Greater(t, len(b), len(block.B))
	})

}

func TestGzipDecompress(t *testing.T) {

	t.Run("decompress data", func(t *testing.T) {
		gd := NewGzipDecompress()
		c := make(chan blockconverter.Block)
		index := int64(7)
		b := []byte{31, 139, 8, 0, 0, 0, 0, 0, 0, 255, 76, 144, 177, 142, 27, 49, 12, 68, 127, 101, 186, 52, 7, 195, 87, 36, 64, 202, 148, 7, 92, 128, 252, 2, 189, 26, 123, 137, 149, 40, 65, 164, 236, 236, 223, 7, 43, 39, 136, 75, 130, 195, 153, 199, 249, 172, 157, 5, 31, 205, 71, 129, 58, 92, 75, 203, 59, 210, 40, 101, 71, 240, 119, 160, 94, 17, 43, 209, 186, 90, 168, 221, 32, 150, 16, 123, 163, 51, 230, 172, 150, 134, 71, 223, 79, 120, 181, 90, 197, 113, 33, 109, 222, 254, 147, 124, 113, 120, 136, 37, 233, 233, 53, 129, 119, 118, 184, 218, 194, 169, 126, 255, 122, 62, 251, 27, 30, 43, 13, 98, 24, 182, 89, 125, 216, 19, 128, 29, 81, 235, 6, 193, 77, 114, 230, 62, 233, 246, 198, 73, 229, 75, 151, 114, 201, 76, 208, 64, 84, 20, 217, 8, 121, 238, 189, 113, 209, 66, 195, 165, 214, 237, 132, 143, 152, 132, 62, 250, 93, 239, 76, 176, 26, 168, 150, 119, 92, 245, 78, 44, 180, 24, 93, 233, 111, 184, 140, 128, 100, 175, 147, 44, 83, 26, 212, 162, 130, 153, 75, 244, 106, 186, 188, 118, 241, 134, 206, 34, 106, 71, 45, 116, 167, 133, 74, 206, 59, 134, 45, 171, 216, 141, 105, 6, 63, 196, 209, 106, 27, 89, 186, 250, 1, 251, 44, 233, 253, 251, 183, 179, 227, 161, 177, 206, 177, 51, 83, 156, 199, 131, 159, 140, 46, 206, 128, 175, 100, 56, 150, 106, 241, 55, 229, 181, 242, 38, 238, 114, 59, 160, 143, 50, 74, 237, 135, 201, 241, 74, 222, 159, 182, 137, 190, 69, 109, 104, 227, 146, 213, 215, 227, 222, 235, 53, 30, 210, 137, 172, 27, 241, 35, 167, 225, 248, 37, 55, 254, 148, 141, 29, 106, 75, 30, 233, 208, 221, 217, 93, 171, 249, 196, 249, 159, 121, 250, 19, 0, 0, 255, 255, 30, 8, 40, 133, 62, 2, 0, 0}

		go func() {
			c <- blockconverter.Block{Index: index, B: b}
		}()

		result := gd.Run(c)

		block := <-result
		assert.Equal(t, index, block.Index)
		assert.Less(t, len(b), len(block.B))
	})
}
