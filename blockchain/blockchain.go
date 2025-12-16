package blockchain

import (
	"crypto/sha256"
	"fmt"
	"sync"
)

type Block struct {
	Data     string `json:"data"`
	Hash     string `json:"hash"`
	PrevHash string `json:"prevHash,omitempty"`
	Height   int    `json:"height"`
}

type blockChain struct {
	blocks []*Block
}

// Single Pattern 으로 만들기
var bc *blockChain
var once sync.Once // 몇개의 채널이 있던 한번만 실행되도록 하기

func (b *Block) calculateHash() {
	hash := sha256.Sum256([]byte(b.Data + b.PrevHash))
	b.Hash = fmt.Sprintf("%x", hash) //16진수 : base16 문자열로 전환
}

func getLastHash() string {
	totalBlocks := len(GetBlockChain().blocks)
	if totalBlocks == 0 {
		return ""
	}
	return GetBlockChain().blocks[totalBlocks-1].Hash
}

func createBlock(data string) *Block {
	newBlock := Block{data, "", getLastHash(), len(GetBlockChain().blocks) + 1}
	newBlock.calculateHash()
	return &newBlock
}

func (b *blockChain) AddBlock(data string) {
	b.blocks = append(b.blocks, createBlock(data))
}

func GetBlockChain() *blockChain {
	if bc == nil {
		once.Do(func() {
			bc = &blockChain{}
			bc.AddBlock("Genesis Block")
		})
	}
	return bc
}

//

func AllBlocks() []*Block {
	return GetBlockChain().blocks
}

func (bc *blockChain) GetBlock(height int) *Block {
	return bc.blocks[height-1]
}
