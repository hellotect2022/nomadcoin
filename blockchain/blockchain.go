package blockchain

type block struct {
	data     string
	hash     string
	prevHash string
}

type blockChain struct {
	blocks []block
}

// Single Pattern 으로 만들기
var b *blockChain

func GetBlockChain() *blockChain {
	if b == nil {
		b = &blockChain{}
	}
	return b
}

//
