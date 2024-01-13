package ebnetworking

import (
	"fmt"
	"net"
)

type IPAddress struct {
	IP   string
	Port int
}

func (i *IPAddress) AsString() string {
	return fmt.Sprintf("%s:%d", i.IP, i.Port)
}

func (i *IPAddress) Resolved() (*net.UDPAddr, error) {
	return net.ResolveUDPAddr("udp", i.AsString())
}
