package peer

import (
	"account"
	"connection"
	"encoding/gob"
	"fmt"
	"net"
	"strconv"
	"strings"
	"utils"
)

type Peer struct {
	Address          string
	Ledger           *account.Ledger
	Connections      []connection.Connection
	AddressesToConns map[string]connection.Connection
}

type String struct {
	Msgfmt string
}

// initialize an empty peer
func (p Peer) Init() {
	l := p.ListenInitial()
	addr, port := utils.GetHostInfo(l)
	p.Address = addr + ":" + port
	p.Connections = make([]connection.Connection, 0)
	p.AddressesToConns = make(map[string]connection.Connection, 0)
	p.Ledger = account.MakeLedger()

	ownConnection := connection.Connection{
		Connection: nil,
		Address:    p.Address,
		Decoder:    nil,
		Encoder:    nil,
	}

	p.Connections = append(p.Connections, ownConnection)
}

func (p Peer) ConductTransaction(tx *account.Transaction) {
	//Need to do some checks, but what?
	ledger := p.Ledger
	ledger.Transact(tx)
}

func (p Peer) FloodTransaction(tx *account.Transaction) {
	for _, c := range p.Connections {
		c.Encoder.Encode("Conduct Transaction")
		c.Encoder.Encode(tx)
	}
}

func (p Peer) JoinNetwork(firstConnection *connection.Connection) {
	err := firstConnection.Encoder.Encode("Ask for connections")

	if err != nil {
		panic(-1)
	}

	var fmt map[string]connection.Connection
	err = firstConnection.Decoder.Decode(&fmt)

	if err != nil {
		panic(-1)
	}

	for k := range fmt {
		conn, err := net.Dial("tcp", k)
		if err != nil {
			panic(-1)
		}
		p.Connections = append(p.Connections, connection.Connection{
			Connection: &conn,
			Address:    k,
			Decoder:    gob.NewDecoder(conn),
			Encoder:    gob.NewEncoder(conn),
		})

	}

}

func (p Peer) SendMessage(msgType string, connection connection.Connection, tx ...account.Transaction) {
	switch msgType {
	case "New peer":
		connection.Encoder.Encode(msgType)
		connection.Encoder.Encode(p.Address)
		return
	case "Ask for connections":
		for _, c := range p.Connections {
			c.Encoder.Encode(msgType)
			c.Encoder.Encode(c.Address)
			// How do i do this? Do i ask for the connecitons for each peer?
			// And then decode it in a loop?
		}
		return
	case "Conduct Transaction":
		for _, transaction := range tx {
			p.ConductTransaction(&transaction)
		}
	}
}

//Initialises the required data and lets the peer listen on a random port
func (p Peer) ListenInitial() net.Listener {
	l, err := net.Listen("tcp", ":")
	if err != nil {
		defer l.Close()
	}
	if err != nil {
		panic(0)
	}
	println("Now listening on " + l.Addr().String())
	go p.listenUtil(l)
	return l
}

//Handles the actual listening
func (p Peer) listenUtil(l net.Listener) {
	for {
		conn, _ := l.Accept()
		println("Got at connection from ", conn.RemoteAddr().String())
		go p.HandleConnection(conn)
	}
}

// If the peer is  unable to connect to the network for any reason
// it simply makes it own
func (p Peer) MakeOwnNetwork() {
	fmt.Println("Listening on port " + strings.Split(p.Address, ":")[0])
}

// Peer tries to connect to the network on addr:port.
// If it fails it initializes its own network
func (p Peer) Connect(addr string, port int) {
	fullAddr := addr + ":" + strconv.Itoa(port)
	c, err := net.Dial("tcp", fullAddr)
	if err != nil {
		fmt.Println("Failed to establish connection, initializing own network")
		p.MakeOwnNetwork()
		return
	}
	firstConnection := &connection.Connection{
		Connection: &c,
		Address:    c.RemoteAddr().String(),
		Decoder:    gob.NewDecoder(c),
		Encoder:    gob.NewEncoder(c),
	}
	p.Connections = append(p.Connections, *firstConnection)
	p.JoinNetwork(firstConnection)
}

func (p Peer) HandleNewConnection(conn net.Conn) {
	defer conn.Close()
	newConnection := &connection.Connection{
		Connection: &conn,
		Address:    conn.LocalAddr().String(),
		Decoder:    gob.NewDecoder(conn),
		Encoder:    gob.NewEncoder(conn),
	}
	p.Connections = append(p.Connections, *newConnection)

}

// Decides what to do with a established connection, depending on the
// type. Type is specified by the first part of the input from the conn
// Object
func (p Peer) HandleConnection(connNew net.Conn) {
	address := connNew.RemoteAddr().String()
	connOld, isPresent := p.AddressesToConns[address]
	if !isPresent {
		//This is a new connection!
		p.HandleNewConnection(connNew)
		return
	}
	inputfmt := &String{}
	err := connOld.Decoder.Decode(inputfmt)
	if err != nil {
		panic(-1)
	}
	msgType := inputfmt.Msgfmt
	switch msgType {
	case "New Peer joined":
		//TODO: Find out how to add a new peers connection
		inputfmt := &Peer{}
		connOld.Decoder.Decode(inputfmt)
		newConn, _ := net.Dial("tcp", inputfmt.Address)
		newConnection := &connection.Connection{
			Connection: &newConn,
			Address:    inputfmt.Address,
			Decoder:    gob.NewDecoder(newConn),
			Encoder:    gob.NewEncoder(newConn),
		}
		p.Connections = append(p.Connections, *newConnection)
		p.AddressesToConns[inputfmt.Address] = *newConnection
		return
	case "Conduct transaction":
		inputfmt := &account.Transaction{}
		err := connOld.Decoder.Decode(inputfmt)
		if err != nil {
			panic(-1)
		}
		transaction := inputfmt
		p.Ledger.Transact(transaction)
		return
	case "Ask for connections":
		var addresses []string
		for k := range p.AddressesToConns {
			addresses = append(addresses, k)
		}
		connOld.Encoder.Encode(addresses)
		return
	}
}
