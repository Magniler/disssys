package peer

import (
	"account"
	"connection"
	"encoding/gob"
	"net"
	"strings"
)

type Peer struct {
	address     string
	ledger      *account.Ledger
	connections []connection.Connection
}

func (p Peer) UpdateLedger(tx *account.Transaction) {
	//Need to do some checks, but what?
	ledger := p.ledger
	ledger.Transact(tx)
}

func (p Peer) FloodTransaction(tx *account.Transaction) {
	for _, c := range p.connections {
		conn := c.Connection
		encoder := gob.NewEncoder(*conn)
		encoder.Encode(tx)
	}
}

func (p Peer) SendMessage(msgType string, connection connection.Connection) {
	switch msgType {
	case "New Connection":
		msg := msgType + ":" + connection.Address
		connection.Encoder.Encode(msg)
		return
	case "Transaction":
		return
	case "Ask for connections":
		for _, c := range p.connections {
			msg := msgType + ":" + c.Address
			c.Encoder.Encode(msg)
		}
		return
	}
}

func (p Peer) Listen() {
	l, _ := net.Listen("tcp", p.address)
	defer l.Close()
	for {
		conn, _ := l.Accept()
		go p.HandleConnection(conn)
	}

}

type String struct {
	Msgfmt string
}

func (p Peer) HandleConnection(conn net.Conn) {
	defer conn.Close()
	decoder := gob.NewDecoder(conn)
	inputfmt := &String{}
	err := decoder.Decode(inputfmt)
	if err != nil {
		panic(-1)
	}
	input := strings.Split(inputfmt.Msgfmt, ":")
	msgType := input[0]
	switch msgType {
	case "New Connection":
		newConnection := &connection.Connection{
			Connection: &conn,
			Address:    input[1],
			Decoder:    decoder,
			Encoder:    gob.NewEncoder(conn),
		}
		p.connections = append(p.connections, *newConnection)
		return
	case "Transaction":
		return
	case "Ask for connections":
		return
	}
}

func (p Peer) Init() {
	p.Listen()
}
