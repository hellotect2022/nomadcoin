package blockchain

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hellotect2022go/nomadcoin/db"
	"github.com/hellotect2022go/nomadcoin/utils"
)

type Block struct {
	Transactions []*Tx  `json:"transactions"`
	Hash         string `json:"hash"`
	PrevHash     string `json:"prevHash,omitempty"`
	Height       int    `json:"height"`
	// 자격증명용
	Difficulty int `json:"difficulty"` //  hash 앞에 오게될 0개의 n 갯수로 조절
	Nonce      int `json:"nonce"`      // 블록체인에서 채굴자들이 수정할 수 있는 유일한 값
	Timestamp  int `json:"timestamp"`
}

var ErrNotFound = errors.New("block not found")

func (b *Block) persist() {
	db.SaveBlock(b.Hash, utils.ToBytes(b))
}

func (b *Block) restore(data []byte) {
	utils.FromBytes(data, b)
}

func FindBlock(hash string) (*Block, error) {
	blockBytes := db.Block(hash)
	if blockBytes == nil {
		return nil, ErrNotFound
	}
	block := &Block{}
	block.restore(blockBytes)
	return block, nil
}

func (b *Block) mine() {
	target := strings.Repeat("0", b.Difficulty)
	for {
		b.Timestamp = int(time.Now().Unix())
		hash := utils.Hash(b)
		fmt.Printf("Target: %s\nHash: %s\nNonce: %d\n", target, hash, b.Nonce)
		if strings.HasPrefix(hash, target) {
			b.Hash = hash
			break
		} else {
			b.Nonce++
		}

	}
}

func createBlock(prevHash string, height int) *Block {
	block := &Block{
		Hash:         "",
		PrevHash:     prevHash,
		Height:       height,
		Difficulty:   GetBlockChain().difficulty(),
		Nonce:        0,
		Transactions: []*Tx{makeCoinbaseTx("dhhan")},
	}

	block.mine()

	//payload := block.Data + block.PrevHash + fmt.Sprint(block.Height)
	//block.Hash = fmt.Sprintf("%x", sha256.Sum256([]byte(payload)))
	block.persist()
	return block
}
