package blockchain

import (
	"fmt"
	"sync"

	"github.com/hellotect2022go/nomadcoin/db"
	"github.com/hellotect2022go/nomadcoin/utils"
)

const (
	defaultDifficulty  int = 2 //  hash 앞에 오게될 0개의 n 갯수로 난이도 조절
	difficultyInterval int = 5 //  height 5개마다 난이도 조절
	blockInterval      int = 2 // 2분마다 블록이 생성되길 원함
	allowedRange       int = 2 // 허용범위
)

type blockChain struct {
	NewestHash        string `json:"newestHash"`
	Height            int    `json:"height"`
	CurrentDifficulty int    `json:"currentDifficulty"`
}

func (bc *blockChain) persist() {
	db.SaveBlockchain(utils.ToBytes(bc))
}

func (bc *blockChain) restore(data []byte) {
	utils.FromBytes(data, bc)
}

// Single Pattern 으로 만들기
var bc *blockChain
var once sync.Once // 몇개의 채널이 있던 한번만 실행되도록 하기

func (bc *blockChain) AddBlock(data string) {
	block := createBlock(data, bc.NewestHash, bc.Height+1)
	bc.NewestHash = block.Hash
	bc.Height = block.Height
	bc.CurrentDifficulty = block.Difficulty // 블록 의 난이도 설정후 체인 난이도 업데이트
	bc.persist()
}

func (bc *blockChain) Blocks() []*Block {
	var blocks []*Block
	hashCursor := bc.NewestHash

	for {
		block, _ := FindBlock(hashCursor)
		blocks = append(blocks, block)
		if block.PrevHash != "" {
			hashCursor = block.PrevHash
		} else {
			break
		}
	}
	return blocks
}

func (bc *blockChain) recalculateDifficulty() int {
	allBlocks := bc.Blocks()
	newestBlock := allBlocks[0]                              // 가장 최신의 블럭
	lastRecalculatedBlock := allBlocks[difficultyInterval-1] // 이전 5번째의 블럭
	actualTime := (newestBlock.Timestamp / 60) - (lastRecalculatedBlock.Timestamp / 60)
	expectedTime := difficultyInterval * blockInterval // 2분에 1개씩 5개의 시간 10분

	if actualTime <= (expectedTime - allowedRange) {
		return bc.CurrentDifficulty + 1
	} else if actualTime >= (expectedTime + allowedRange) {
		return bc.CurrentDifficulty - 1
	}
	return bc.CurrentDifficulty
}

func (bc *blockChain) difficulty() int {
	if bc.Height == 0 {
		return defaultDifficulty
	} else if bc.Height%difficultyInterval == 0 {
		return bc.recalculateDifficulty()
	} else {
		return bc.CurrentDifficulty
	}
}

func GetBlockChain() *blockChain {
	if bc == nil {
		once.Do(func() {
			bc = &blockChain{
				Height: 0,
			}
			//search for checkpoint on the db
			// restore bc from byte
			checkpoint := db.Checkpoint()
			if checkpoint == nil {
				bc.AddBlock("Genesis Block")
			} else {
				fmt.Println("Restoring....")
				bc.restore(checkpoint)
			}
		})
	}
	fmt.Printf("NewestHash: %s\nHeight: %d\n\n", bc.NewestHash, bc.Height)
	return bc
}
