// tx.go
package main

import (
	"fmt"
)

type TxOutput struct {
	Value  int
	PubKey string
}

// note: each output is indivisible
type TxInput struct {
	ID  []byte
	Out int
	Sig string
}

func (in *TxInput) CanUnlock(data string) bool {
	return in.Sig == data
}

func (out *TxOutput) CanBeUnlocked(data string) bool {
	return out.PubKey == data
}
