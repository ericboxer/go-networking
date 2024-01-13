package ebnetworking

import (
	"fmt"
	"net"
)

type UDPDataHandler func([]byte, *net.UDPAddr)

type UDPComms struct {
	buffer       []byte
	conn         *net.UDPConn
	localAddress IPAddress
	sendPort     int
	handler      UDPDataHandler
	closeChan    chan struct{}
}

func (u *UDPComms) Init(localAddress IPAddress, sendPort int) error {
	u.localAddress = localAddress
	u.sendPort = sendPort
	u.handler = u.handler
	u.closeChan = make(chan struct{})

	addr, err := u.localAddress.Resolved()
	if err != nil {
		return fmt.Errorf("error resolving local address: %w", err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return fmt.Errorf("error listening on port %d: %w", u.localAddress.Port, err)
	}

	u.conn = conn
	return nil
}

func (u *UDPComms) SetHandler(handler UDPDataHandler) {
	u.handler = handler
}

func (u *UDPComms) Start() {
	go func() {
		for {
			select {
			case <-u.closeChan:
				return
			default:
				buffer := make([]byte, 1024) // Consider making this configurable
				n, addr, err := u.conn.ReadFromUDP(buffer)
				if err != nil {
					fmt.Printf("Error reading from UDP: %s\n", err)
					continue
				}
				if u.handler != nil {
					u.handler(buffer[:n], addr)
				}
			}
		}
	}()
}

func (u *UDPComms) Send(data []byte, host IPAddress) (int, error) {
	addr, err := host.Resolved()
	if err != nil {
		return 0, fmt.Errorf("error resolving remote address: %w", err)
	}

	bytesWritten, err := u.conn.WriteToUDP(data, addr)
	if err != nil {
		return 0, fmt.Errorf("error sending data: %w", err)
	}

	return bytesWritten, nil
}

func (u *UDPComms) Close() {
	close(u.closeChan)
	u.conn.Close()
}


func (u *UDPComms) defaultHandler(data []byte,addr *net.UDPAddr) {
	fmt.Println("%s:%d->%v", addr.IP, addr.Port, data)
}