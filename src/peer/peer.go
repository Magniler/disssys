package peer

import (
	"account"
	"connection"
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"strconv"
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
	case "New peer":
		msg := msgType + ":" + connection.Address
		connection.Encoder.Encode(msg)
		return
	case "Ask for connections":
		for _, c := range p.connections {
			msg := msgType + ":" + c.Address
			err := c.Encoder.Encode(msg)
			if err != nil {
				panic(-1)
			}
		}
		return
	}
}

func (p Peer) Listen() {
	l, _ := net.Listen("tcp", p.address)
	defer l.Close()
	for {
		fmt.Println("Now listening on " + p.address)
		conn, _ := l.Accept()
		go p.HandleConnection(conn)
	}

}

// If the peer is  unable to connect to the network for any reason
// it simply makes it own
func (p Peer) MakeOwnNetwork() {
	name, _ := os.Hostname()
	ip, _ := net.LookupHost(name)
	// We need to join the ip to make into one string
	hostPort := ip[0] + ":" + "1616"
	// :0 os suppossed to be an ephereal port, which means the os
	// finds the next avaivable port
	addr, err := net.ResolveTCPAddr("tcp", hostPort)
	if err != nil {
		panic(-1)
	}
	conn, err := net.Dial("tcp", addr.String())
	if err != nil {
		panic(-1)
	}
	ownConnection := connection.Connection{
		Connection: &conn,
		Address:    addr.String(),
		Decoder:    gob.NewDecoder(conn),
		Encoder:    gob.NewEncoder(conn),
	}
	p.connections = append(p.connections, ownConnection)
	p.Listen()
}

// Peer tries to connect to the network on addr:port.
// If it fails it initializes its own network
func (p Peer) Connect(addr string, port int) {
	fullAddr := addr + ":" + strconv.Itoa(port)
	c, err := net.Dial("tcp", fullAddr)
	if err != nil {
		fmt.Println("Failed to establish connection, initialising own network")
		p.MakeOwnNetwork()
		return
	}
	newConnection := &connection.Connection{
		Connection: &c,
		Address:    fullAddr,
		Decoder:    gob.NewDecoder(c),
		Encoder:    gob.NewEncoder(c),
	}
	// TODO: Right now it does not add its own conneciton to its lis of
	// connections
	p.connections = append(p.connections, *newConnection)

type String struct {
	Msgfmt string
}

// Decides what to do with a newly established connection, depending on the
// type. Type is specified by the first
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
	case "New Peer":
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
}
