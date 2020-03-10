package compress

type CompressStrategy interface {
	Run(jobs <-chan []byte) <-chan []byte
}

type ContextCompress interface {
	GetCompressStrategy(algorithm string) CompressStrategy
}

type contextCompress struct {
	strategies map[string]CompressStrategy
}

func NewContextCompress(strategies map[string]CompressStrategy) ContextCompress {
	return contextCompress{strategies}
}

func (c contextCompress) GetCompressStrategy(algorithm string) CompressStrategy {
	return c.strategies[algorithm]
}
