package peer

import (
	"account"
	"encoding/gob"
	"net"
)

type Peer struct {
	ledger      account.Ledger
	connections []Connection
}

func (p Peer) UpdateLedger(tx *account.Transaction) {
	//Need to do some checks, but what?
	ledger := *p.ledger
	ledger.Transact(tx)
}

func (p Peer) FloodMessage(msg *nt.Ledger) {
	encodedmsg := encoder
	for i, conn := range p.connections {
		conn, err := net.Dial("tcp", conn.Address.String())
	}
}

func (p Peer) FloodTransaction(tx *Transaction) {
	for i, c := range p.connections {
		conn, err := net.Dial(tcp, c.address)
		encoder := gob.NewEncoder(conn)
		encoder.Encode(msg)
	}
}

func (p Peer) Listen(c *chan account.Ledger) {
	l, _ = net.Listen("tcp", Peer.address)
	for {
		conn, eer := l.Accept()
		if err != nil {
			panic(err)
		}
		msg := conn.
	}
}

func (p Peer) Init() {
	p.TryToConnect()
	p.Listen()
}


func (p Peer) TryToConnect() {
	for i, conn := range p.connections {
		go HandleConnection(conn.Conneciton)
	}
}

type Connection struct {
	Connection net.conn
	Address    string
	Decoder    gob.decoder
	Encoder    gob.encoder
}
