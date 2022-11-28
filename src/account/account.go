package account

import (
	"sync"
)

type Transaction struct {
	From   string
	To     string
	Amount int
}

type Ledger struct {
	Accounts map[string]int
	lock     sync.Mutex
}

func MakeLedger() *Ledger {
	ledger := new(Ledger)
	ledger.Accounts = make(map[string]int)
	return ledger
}

func (l *Ledger) Transact(t *Transaction) {
	l.lock.Lock()
	defer l.lock.Unlock()
	_, toFound := l.Accounts[t.To]
	_, fromFound := l.Accounts[t.From]
	if toFound && fromFound {
		l.Accounts[t.From] -= t.Amount
		l.Accounts[t.To] += t.Amount
	}
}
