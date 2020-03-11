package blockcompressor

import (
	"blockcompressor/pkg/compress"
	g "blockcompressor/pkg/compress/gzip"
	"blockcompressor/pkg/reader"
	"blockcompressor/pkg/writer"
	"log"
	"os"
	"path"
)

func Run(filePathIn string, filePathOut string, blockSize int, workers int, algorithm string) {
	file, err := os.Open(filePathIn)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	br := reader.NewReader(file, blockSize)

	m := make(map[string]compress.CompressStrategy)
	cleanPath := path.Clean(filePathIn)
	m["c"] = g.NewGzipCompress(cleanPath)
	m["d"] = g.NewGzipDecompress(cleanPath)
	context := compress.NewContextCompress(m)

	outFile, err := os.Create(filePathOut)
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()
	bw := writer.NewWriter(outFile)

	app := NewApp(br, context, bw, workers)
	app.Run(algorithm)
}
