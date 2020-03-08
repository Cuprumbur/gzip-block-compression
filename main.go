package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path"
	"strings"
	"sync"

)

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

	r := bufio.NewReader(fileReader)

	in := make(chan []byte, 3)
	out := make(chan []byte, 3)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			data := make([]byte, blockSize)

			n, err := r.Read(data)

			if err == io.EOF {
				fmt.Println("END read")
				close(in)
				fmt.Println("closed 'in'")
				break
			}

			fmt.Printf("read %d bytes\n", n)

			if n < blockSize {
				data = data[:n]
			}

			in <- data
		}

	}()
	wg.Add(1)
	go func() {
		defer wg.Done()

		var wgZip sync.WaitGroup
		for {

			data, ok := (<-in)

			if !ok {
				go func() {
					wgZip.Wait()
					close(out)
					fmt.Println("closed 'out'")
				}()
				break
			}
			wgZip.Add(1)
			go func() {
				defer wgZip.Done()
				var buf bytes.Buffer
				w := gzip.NewWriter(&buf)
				w.Name = pathFile
				n, err := w.Write(data)
				if err != nil {
					log.Fatal(err)
				}
				w.Close()

				fmt.Printf("compressed %d bytes\n", n)

				select {
				case out <- buf.Bytes():
				default:
				}

			}()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		fileNameWithoutExt := strings.TrimSuffix(pathFile, path.Ext(pathFile))

		outputFile := fileNameWithoutExt + ".gzip"
		newZipFile, err := os.Create(outputFile)
		check(err)

		defer newZipFile.Close()

		for {
			data, ok := (<-out)
			if !ok {
				fmt.Println("stop write")
				break
			}

			n, err := newZipFile.Write(data)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("wrote %d bytes\n", n)
		}

	}()
	go func() {
		http.ListenAndServe("localhost:8080", nil)

	}()
	wg.Wait()
	fmt.Printf("done.")
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var blockSize = 1000000
