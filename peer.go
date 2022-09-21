package peer

import (
	"peer/account"
)

type Peer struct {
	ledger account.Ledger
	peers  list[Peer]
	ip     int
	addr   string
}

func FloodMessage(msg *account.Ledger) *Peer {
	for i, e := range peers {
		ch := make(chan<- account.Ledger)
		ch <- msg
	}
}

func FloodTransaction(tx *Transaction) *Peer {

}

func Recive(c *chan account.Ledger) *Peer {
	ledger := <-c
}
