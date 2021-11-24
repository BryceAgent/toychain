// transaction.go
package main

import (
	"fmt"
)

type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

type TxOutput struct {
	// represents amount of coins in transaction
	Value int
	// public key
	PubKey string
}

type TxInput struct {
	// ID is label for the transaction
	ID []byte
	// index of specific output within transaction
	// for e.g. transaction has 4 outputs, Out field specifies which
	Out int
	// script to add data to outputs' PubKey
	// here identical to PubKey for simplicity
	Sig string
}

const reward = 100

// coinbase (first) transaction
func CoinbaseTx(toAddress, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Coins to %s", toAddress)
	}

	// has no previous output so initialize with no ID = -1
	txIn := TxInput{[]byte{}, -1, data}
	txOut := TxOutput{reward, toAddress}
	tx := Transaction{nil, []TxInput{txIn}, []TxOutput{txOut}}

	return &tx
}

func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(tx)
	Handle(err)

	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

func (in *TxInput) CanUnlock(data string) bool {
	return in.Sig == data
}

func (out *TxOutput) CanBeUnlocked(data string) bool {
	return out.PubKey == data
}

func (tx *Transaction) IsCoinbase() bool {
	// this checks a transaction and returns true if newly minted
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].Out == -1
}

// to initiate transactions
func NewTransaction(from, to string, amount int, chain *BlockChain) *Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	// step 1 gather spendable outputs
	acc, validOutputs := chain.FindSpendableOutputs(from, amount)

	// step 2 check if funds are sufficient for given amount
	if acc < amount {
		log.Panic("Error: Not enough funds!")
	}

	// step 3 point to outputs to spend
	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		Handle(err)

		for _, out := range outs {
			input := TxInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, TxOutput{amount, to})

	// step 4 make new outputs for leftover amount
	if acc > amount {
		outputs = append(outputs, TxOutput{acc - amount, from})
	}

	// step 5 initialize new transaction with inputs and outputs
	tx := Transaction{nil, inputs, outputs}

	// step 6 set new ID and return
	tx.SetID()

	return &tx
}
