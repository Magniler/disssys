package utils

import (
	"net"
	"os"
)

func GetHostInfo(listener net.Listener) (string, string) {
	name, _ := os.Hostname()
	//errors.PrintIfError(err, "Error while getting host name from os")

	addrs, _ := net.LookupHost(name)
	//errors.PrintIfError(err, "Error while looking up address of host")

	_, port, _ := net.SplitHostPort(listener.Addr().String())
	//errors.PrintIfError(err, "Error while getting port of tcp listener")

	return addrs[0], port
}
