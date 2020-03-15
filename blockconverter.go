package blockconverter

import (
	"encoding/binary"
	"io"
)

type Block struct {
	Index int64
	B     []byte
}

type Reader interface {
	Read() <-chan Block
}

type Converter interface {
	Run(<-chan Block) <-chan Block
}

type Writer interface {
	Write(<-chan Block) <-chan struct{}
}

type App struct {
	r Reader
	c Converter
	w Writer
}

type AppFactory interface {
	NewReader() Reader
	NewConverter() Converter
	NewWriter() Converter
}

func NewApp(r Reader, c Converter, w Writer) App {
	return App{r, c, w}
}

func (a App) Convert() {

	block := a.r.Read()

	convertedBlock := a.c.Run(block)

	<-a.w.Write(convertedBlock)
}

func toByte(i int64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(i))
	return b
}

func SetData(r io.Writer, index int64, size int64) error {

	_, err := r.Write(toByte(index))
	if err != nil {
		return err
	}
	_, err = r.Write(toByte(size))

	if err != nil {
		return err
	}
	return nil
}
