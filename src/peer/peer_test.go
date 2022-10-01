package peer_test

import (
	"peer"
)

func ExamplePeer_Connect() {
	p := peer.Peer{}
	p.Connect("Giberish", 42)
	// Output: Failed to establish connection, initialising own network
	//FIXME: Throws exception, probaply an permissions proble, let andreas see try it.
}
