package main

import (
	"archive/zip"
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"sync"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func main() {
	var pathFile string
	flag.StringVar(&pathFile, "file", "", "Path of file to compress")
	flag.Parse()

	if pathFile == "" {
		log.Fatal("file not provided")
	}

	fmt.Println("File to compress ", pathFile)

	fileReader, err := os.Open(pathFile)
	check(err)

	fileNameWithoutExt := strings.TrimSuffix(pathFile, path.Ext(pathFile))

	outputFile := fileNameWithoutExt + ".zip"

	newZipFile, err := os.Create(outputFile)
	check(err)

	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	f, err := zipWriter.Create(pathFile)
	check(err)

	c := make(chan []byte, 3)
	r := bufio.NewReader(fileReader)
	w := bufio.NewWriter(f)

	size := r.Size()
	count := size/blockSize + 1
	var wg sync.WaitGroup
	wg.Add(count)

	for i := 0; i < count; i++ {
		go read(c, r, &wg)
		write(c, w)
	}

	wg.Wait()
}

var blockSize = 1000

func read(c chan<- []byte, r *bufio.Reader, wg *sync.WaitGroup) {
	data := make([]byte, blockSize)
	n, err := r.Read(data)

	if err == io.EOF {
		// close(c)
		fmt.Println("EOF")

		return
	}

	fmt.Printf("read  %d bytes\n", n)

	if n < blockSize {
		data = data[:n]
	}
	c <- data
	wg.Done()
}

func write(c <-chan []byte, w *bufio.Writer) {
	t, ok := (<-c)
	if !ok {
		fmt.Print("Chanel is closed")
		return
	}

	n, err := w.Write(t)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("wrote %d bytes\n", n)
}
