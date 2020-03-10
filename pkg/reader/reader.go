package reader

type Reader interface {
	Read() <-chan []byte
}
