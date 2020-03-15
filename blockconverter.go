package blockconverter

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
