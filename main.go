package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"log"
	_ "net/http/pprof"
	"os"
	"path"
	"strings"
	"sync"

)

func main() {
	// println(runtime.GOMAXPROCS(10))
	var pathFile string
	flag.StringVar(&pathFile, "file", "", "Path of file to compress")
	flag.Parse()

	if pathFile == "" {
		log.Fatal("file not provided")
	}

	fmt.Println("File to compress ", pathFile)

	fileReader, err := os.Open(pathFile)
	check(err)

	r := bufio.NewReader(fileReader)

	in := make(chan []byte)
	out := make(chan []byte)
	// workers
	var workerWg sync.WaitGroup
	workers := 1
	for i := 0; i < workers; i++ {
		workerWg.Add(1)
		go compressWorker(&workerWg, in, out, pathFile)
	}

	go read(r, in)

	fmt.Println("closed 'out'")

	fileNameWithoutExt := strings.TrimSuffix(pathFile, path.Ext(pathFile))

	outputFile := fileNameWithoutExt + ".gzip"
	newZipFile, err := os.Create(outputFile)
	check(err)
	defer newZipFile.Close()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for data := range out {

			n, err := newZipFile.Write(data)
			if err != nil {
				log.Fatal(err)
			}
			// _ = n
			fmt.Printf("wrote %d bytes\n", n)
		}

		fmt.Println("stop write")
	}()
	workerWg.Wait()
	close(out)
	wg.Wait()
	// go func() {
	// 	http.ListenAndServe("localhost:8080", nil)

	// }()

	fmt.Printf("done.")
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var blockSize = 1000000

func read(r *bufio.Reader, in chan<- []byte) {
	for {
		data := make([]byte, blockSize)

		n, err := r.Read(data)

		if err == io.EOF {
			fmt.Println("END read")
			break
		}
		check(err)

		fmt.Printf("read %d bytes\n", n)

		if n < blockSize {
			data = data[:n]
		}

		in <- data
	}

	fmt.Println("closed 'in'")
	close(in)
}
func compressWorker(wg *sync.WaitGroup, in <-chan []byte, out chan<- []byte, pathFile string) {
	defer wg.Done()
	for data := range in {
		var buf bytes.Buffer
		w := gzip.NewWriter(&buf)
		defer w.Close()
		w.Name = pathFile
		n, err := w.Write(data)
		if err != nil {
			log.Fatal(err)
		}
		// _ = n
		fmt.Printf("compressed %d bytes\n", n)
		out <- buf.Bytes()
	}
}
