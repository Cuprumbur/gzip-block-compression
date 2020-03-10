package compress

import (
	"blockcompressor/pkg/compress"
	"bytes"
	"compress/gzip"
	"log"
)

type gzipDecompress struct {
	fileName string
}

func NewGzipDecompress(fileName string) compress.CompressStrategy {
	return gzipDecompress{fileName}
}

func (g gzipDecompress) Run(jobs <-chan []byte) <-chan []byte {
	result := make(chan []byte)

	go func() {
		for data := range jobs {
			var buf bytes.Buffer
			b := bytes.NewReader(data)
			r, err := gzip.NewReader(b)
			if err != nil {
				log.Fatal(err)
			}
			r.Name = g.fileName
			_, err = r.Read(data)
			if err != nil {
				log.Fatal(err)
			}
			r.Close()

			result <- buf.Bytes()
		}

		close(result)
	}()

	return result
}
