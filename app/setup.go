package blockcompressor

import (
	"blockcompressor/pkg/command"
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

	m := make(map[string]command.Command)
	cleanPath := path.Clean(filePathIn)
	m["c"] = command.NewGzipCompress(cleanPath)
	m["d"] = command.NewGzipDecompress(cleanPath)
	context := command.NewContext(m)

	outFile, err := os.Create(filePathOut)
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()
	bw := writer.NewWriter(outFile)

	app := NewApp(br, context, bw, workers)
	app.Start(algorithm)
}
