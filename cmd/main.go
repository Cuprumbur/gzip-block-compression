package main

import (
	"blockconverter"
	"blockconverter/converter"
	"blockconverter/reader"
	"blockconverter/writer"
	"log"
	"os"
	"time"
)

func main() {
	start := time.Now()

	command := os.Args[1]
	pathFile := os.Args[2]
	destPathFile := os.Args[3]

	if command == "" {
		log.Fatal("command not provided")
	}
	if pathFile == "" {
		log.Fatal("file not provided")
	}

	if destPathFile == "" {
		log.Fatal("dest path file not provided")
	}
	workers := 10
	maxBlockSize := int64(1000000)

	m := map[string]AppFactory{
		"c": NewCompressFactory(pathFile, destPathFile, workers, maxBlockSize),
		"d": NewDecompressFactory(pathFile, destPathFile, workers),
	}

	f := m[command]

	a := blockconverter.NewApp(f.NewReader(), f.NewConverter(), f.NewWriter())

	a.Convert()

	elapsed := time.Since(start)
	log.Printf("elapsed %s", elapsed)
}

type AppFactory interface {
	NewReader() blockconverter.Reader
	NewConverter() blockconverter.Converter
	NewWriter() blockconverter.Writer
}

type compressFactory struct {
	fileIn       string
	fileOut      string
	workers      int
	maxBlockSize int64
}

func NewCompressFactory(fileIn string, fileOut string, workers int, maxBlockSize int64) AppFactory {
	return compressFactory{fileIn, fileOut, workers, maxBlockSize}
}

func (a compressFactory) NewReader() blockconverter.Reader {
	return reader.NewFileReader(a.fileIn, reader.NewBlockInfo(reader.StartBlockIndex, a.maxBlockSize))
}

func (a compressFactory) NewConverter() blockconverter.Converter {
	return converter.NewGzipWorkers(a.workers, converter.NewGzipCompress())
}

func (a compressFactory) NewWriter() blockconverter.Writer {
	return writer.NewFileWriter(a.fileOut, "c")
}

type decompressFactory struct {
	fileIn  string
	fileOut string
	workers int
}

func NewDecompressFactory(fileIn string, fileOut string, workers int) AppFactory {
	return decompressFactory{fileIn, fileOut, workers}
}

func (a decompressFactory) NewReader() blockconverter.Reader {
	return reader.NewFileReader(a.fileIn, reader.NewCompressedBlockInfo())
}

func (a decompressFactory) NewConverter() blockconverter.Converter {
	return converter.NewGzipWorkers(a.workers, converter.NewGzipDecompress())
}

func (a decompressFactory) NewWriter() blockconverter.Writer {
	return writer.NewFileWriter(a.fileOut, "d")
}
