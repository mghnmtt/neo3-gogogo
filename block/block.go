package block

import (
	//"github.com/mghnmtt/neo-gogogo/rpc/models"
	"github.com/mghnmtt/neo3-gogogo/tx"
)

type Block struct {
	Header
	Transactions []tx.Transaction
}

func (b *Block) GetSize() int {
	sz := 0
	for _, tx := range b.Transactions {
		sz += tx.GetSize()
	}
	return b.Header.GetSize() + sz
}
