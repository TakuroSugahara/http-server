package main

import (
	"io"
	"net"
)

type responser struct {
	conn net.Conn
}

func NewResponser(conn net.Conn) *responser {
	return &responser{
		conn: conn,
	}
}

func (r *responser) WriteOK() {
	io.WriteString(r.conn, "HTTP/1.1 200 OK \r\n")
}
