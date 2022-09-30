package connection

import (
	"encoding/gob"
	"net"
)

type Connection struct {
	Connection *net.Conn
	Address    string
	Decoder    *gob.Decoder
	Encoder    *gob.Encoder
}
