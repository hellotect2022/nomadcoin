package blockchain

import (
	"errors"
	"time"

	"github.com/hellotect2022go/nomadcoin/utils"
)

const (
	minerReward int = 50
)

type mempool struct {
	// 확정되지 않은 거래내역
	Txs []*Tx
}

var Mempool *mempool = &mempool{}

type Tx struct {
	Id        string   `json:"id"`
	Timestamp int      `json:"timestamp"`
	TxIns     []*TxIn  `json:"txIns"`
	TxOuts    []*TxOut `json:"txOuts"`
}

func (t *Tx) getId() {
	t.Id = utils.Hash(t)
}

type TxIn struct {
	//Amount int
	// 이전 트랜잭션 에서의 값을 참조해서 가져옴
	TxID  string `json:"txId"`
	Index int    `json:"index"`
	Owner string `json:"ownere"`
}

type TxOut struct {
	Owner  string `json:"owner"`
	Amount int    `json:"amount"`
}

type UTxOut struct {
	TxID   string
	Index  int
	Amount int
}

func isOnMempool(UTxOut *UTxOut) bool {
	exists := false
Outer:
	for _, tx := range Mempool.Txs {
		for _, input := range tx.TxIns {
			if input.TxID == UTxOut.TxID && input.Index == UTxOut.Index {
				exists = true
				break Outer
			}

		}
	}
	return exists
}

func makeCoinbaseTx(address string) *Tx {
	// 돈을 주는 사람 COINBASE 채굴보상
	txIns := []*TxIn{
		{"", -1, "COINBASE"},
	}

	// 돈을 받는 사람 (채굴자)
	txOuts := []*TxOut{
		{address, minerReward},
	}

	tx := &Tx{
		Id:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId()

	return tx
}

func makeTx(from, to string, amount int) (*Tx, error) {
	if BalanceByAddress(from, GetBlockChain()) < amount {
		return nil, errors.New("Not enough money")
	}

	var txOuts []*TxOut
	var txIns []*TxIn
	total := 0

	// unspent txOuts
	uTxOuts := UTxOutsByAddress(from, GetBlockChain())
	for _, uTxOut := range uTxOuts {
		if total >= amount {
			break
		}
		txIn := &TxIn{uTxOut.TxID, uTxOut.Index, from}
		txIns = append(txIns, txIn)
		total += uTxOut.Amount
	}
	// 잔돈은 다시 나에게로 전달
	if change := total - amount; change != 0 {
		changeTxOut := &TxOut{from, change}
		txOuts = append(txOuts, changeTxOut)
	}

	txOut := &TxOut{to, amount}
	txOuts = append(txOuts, txOut)

	tx := &Tx{
		Id:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId()
	return tx, nil
}

func (m *mempool) AddTx(to string, amount int) error {
	tx, err := makeTx("dhhan", to, amount)
	if err != nil {
		return err
	}
	m.Txs = append(m.Txs, tx)
	return nil
}

func (m *mempool) TxToConfirm() []*Tx {
	coinbase := makeCoinbaseTx("dhhan")
	txs := m.Txs
	txs = append(txs, coinbase) // mempool 거래내역들에 + 코인 채굴에 대한 거래내역도 포함
	m.Txs = nil                 // mempool 에서 삭제함
	return txs
}
