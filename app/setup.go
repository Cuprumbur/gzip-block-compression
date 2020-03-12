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
	var br reader.Reader
	if algorithm == "c" {
		br = reader.NewReader(file, blockSize)
	} else if algorithm == "d" {

		br = reader.NewReaderGzip(file, blockSize)
	} else {
		log.Fatal("no reader")
	}

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

	var bw writer.Writer
	if algorithm == "c" {
		bw = writer.NewWriter(outFile)
	} else if algorithm == "d" {
		bw = writer.NewDecWriter(outFile)
	}

	app := NewApp(br, context, bw, workers)
	app.Start(algorithm)

}
