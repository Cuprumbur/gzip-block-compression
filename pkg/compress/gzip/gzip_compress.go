package compress

import (
	"blockcompressor/pkg/compress"
	"bytes"
	"compress/gzip"
	"log"
)

type gzipCompress struct {
	fileName string
}

func NewGzipCompress(fileName string) compress.CompressStrategy {
	return gzipCompress{fileName}
}

func (g gzipCompress) Run(jobs <-chan []byte) <-chan []byte {
	result := make(chan []byte)

	go func() {
		for data := range jobs {
			var buf bytes.Buffer
			w := gzip.NewWriter(&buf)

			w.Name = g.fileName
			_, err := w.Write(data)
			if err != nil {
				log.Fatal(err)
			}
			w.Close()

			result <- buf.Bytes()
		}

		close(result)
	}()

	return result
}
