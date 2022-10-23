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
	Address     string
	Ledger      *account.Ledger
	Connections []connection.Connection
}

func (p Peer) UpdateLedger(tx *account.Transaction) {
	//Need to do some checks, but what?
	ledger := p.Ledger
	ledger.Transact(tx)
}

func (p Peer) FloodTransaction(tx *account.Transaction) {
	for _, c := range p.Connections {
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
		for _, c := range p.Connections {
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
	l, err := net.Listen("tcp", p.Address)
	if err != nil {
		defer l.Close()
	}
	if err != nil {
		panic(0)
	}
	println("Now listening on " + l.Addr().String())
	for {
		conn, _ := l.Accept()
		println("Got at connection from ", conn.RemoteAddr().String())
		go p.HandleConnection(conn)
	}
}

// If the peer is  unable to connect to the network for any reason
// it simply makes it own
func (p Peer) MakeOwnNetwork() {
	name, _ := os.Hostname()
	ip, _ := net.LookupHost(name)
	// We need to join the ip to make into one string
	hostPort := ip[0] + ":" + "16160"
	// :16160 is chosen at random
	// TODO: Make it soo the port isn't hardcoded
	addr, err := net.ResolveTCPAddr("tcp", hostPort)
	if err != nil {
		fmt.Println(err)
		panic(-1)
	}
	ownConnection := connection.Connection{
		Connection: nil,
		Address:    addr.String(),
		Decoder:    nil,
		Encoder:    nil,
	}
	p.Connections = append(p.Connections, ownConnection)
	p.Listen()
}

// Peer tries to connect to the network on addr:port.
// If it fails it initializes its own network
func (p Peer) Connect(addr string, port int) {
	p.Connections = make([]connection.Connection, 0)
	fullAddr := addr + ":" + strconv.Itoa(port)
	c, err := net.Dial("tcp", fullAddr)
	if err != nil {
		fmt.Println("Failed to establish connection, initializing own network")
		p.MakeOwnNetwork()
		return
	}
	newConnection := &connection.Connection{
		Connection: &c,
		Address:    c.LocalAddr().String(),
		Decoder:    gob.NewDecoder(c),
		Encoder:    gob.NewEncoder(c),
	}
	println(c.LocalAddr().String())
	// TODO: Right now it does not add its own conneciton to its lis of
	// connections
	p.Connections = append(p.Connections, *newConnection)
}

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
		p.Connections = append(p.Connections, *newConnection)
		return
	case "Transaction":
		return
	case "Ask for connections":
		return
	}
}
