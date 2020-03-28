package blockconverter

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
)

type Block struct {
	Index int64
	Data  []byte
}

func (b *Block) InitFrom(r io.Reader) error {
	i := make([]byte, sizeArray)
	_, err := r.Read(i)
	if err != nil {
		return err
	}

	s := make([]byte, sizeArray)
	_, err = r.Read(s)
	if err != nil {
		return err
	}
	length := toInt64(s)

	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)
	_, err = io.CopyN(w, r, length)

	w.Flush()

	b.Index = toInt64(i)
	b.Data = buf.Bytes()

	return nil
}

func (b *Block) Init(r io.Reader, index int64, maxSizeBlock int64) error {
	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)
	_, err := io.CopyN(w, r, maxSizeBlock)

	b.Index = index
	b.Data = buf.Bytes()

	if err != nil {
		return err
	}
	err = w.Flush()
	if err != nil {
		return err
	}

	return nil
}

func (b Block) WriteToWithBlockInfo(w io.Writer) error {
	_, err := w.Write(toByte(b.Index))
	if err != nil {
		return err
	}

	_, err = w.Write(toByte(int64(len(b.Data))))
	if err != nil {
		return err
	}
	_, err = w.Write(b.Data)

	return err
}

const sizeArray = 8

func toInt64(b []byte) int64 {
	i := int64(binary.BigEndian.Uint64(b))
	return i
}

func toByte(i int64) []byte {
	b := make([]byte, sizeArray)
	binary.BigEndian.PutUint64(b, uint64(i))
	return b
}
