package transport

import "net"

type Transporter interface {
	Send(conn net.Conn, req interface{}) error
	Recv(conn net.Conn, resp interface{}) error
}
