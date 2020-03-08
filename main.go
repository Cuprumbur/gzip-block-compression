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

	fileNameWithoutExt := strings.TrimSuffix(pathFile, path.Ext(pathFile))

	outputFile := fileNameWithoutExt + ".zip"

	newZipFile, err := os.Create(outputFile)
	check(err)

	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	f, err := zipWriter.Create(pathFile)
	check(err)
	r := bufio.NewReader(fileReader)
	w := bufio.NewWriter(f)
	size := r.Size()
	i := 0
	blockSize := 100000
	for ; ; i++ {
		data := make([]byte, blockSize)
		numberReadBytes, err := r.Read(data)

		if err == io.EOF {
			fmt.Println("END")
			break
		}

		fmt.Printf("read %d bytes\n", numberReadBytes)

		if numberReadBytes < blockSize {
			data = data[:numberReadBytes]
		}

		numberWriteBytes, err := w.Write(data)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("wrote %d bytes\n", numberWriteBytes)
	}

	fmt.Printf("done. size = %d   i = %d", size, i)

}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// func read(c chan<- []byte, r *bufio.Reader) error {
// 	data := make([]byte, blockSize)
// 	n, err := r.Read(data)

// 	if err == io.EOF {
// 		return err
// 	}

// 	fmt.Printf("read  %d bytes\n", n)

// 	if n < blockSize {
// 		data = data[:n]
// 	}
// 	c <- data
// 	return nil
// }

// func write(c <-chan []byte, w *bufio.Writer, wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	t, ok := (<-c)
// 	if !ok {
// 		fmt.Print("Chanel is closed")
// 		return
// 	}

// 	n, err := w.Write(t)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Printf("wrote %d bytes\n", n)

// }
