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
	Owner  string
	Amount int
}

type TxOut struct {
	Owner  string
	Amount int
}

func makeCoinbaseTx(address string) *Tx {
	// 돈을 주는 사람 COINBASE 채굴보상
	txIns := []*TxIn{
		{"COINBASE", minerReward},
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
	if GetBlockChain().BalanceByAddress(from) < amount {
		return nil, errors.New("not enough money")
	}

	var txIns []*TxIn
	var txOuts []*TxOut

	total := 0

	// 거래하기 전의 내 계좌에 있는 돈들을 나타냄
	// ⭐초과 되는 금액은 그후에 거스름돈으로 받음
	oldTxOuts := GetBlockChain().TxOutsByAddress(from)
	for _, txOut := range oldTxOuts {
		if total > amount {
			break
		}
		txIn := &TxIn{txOut.Owner, txOut.Amount}
		txIns = append(txIns, txIn)
		total += txIn.Amount
	}
	change := total - amount // 거슬러 줘야 하는돈
	if change != 0 {
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
