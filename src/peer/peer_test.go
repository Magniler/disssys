package peer_test

import (
	"net"
	"os"
	"peer"
	"testing"
	"time"
)

func ExamplePeer_MakeOwnNetwork() {
	p := peer.Peer{}
	p.Connect("Giberish", 42)
	//Output: Failed to establish connection, initializing own network
}

func TestPeer_Connect(t *testing.T) {
	p1 := peer.Peer{}
	go p1.Connect("Giberish", 42)
	// This peer then initializes its own network on 192.168.1.152:16160
	time.Sleep(5 * time.Second)
	p2 := peer.Peer{}
	osname, _ := os.Hostname()
	ip, _ := net.LookupHost(osname)
	go p2.Connect(ip[0], 16160)
	println(p2.Connections)
}
