package server

import (
	"bufio"
	"fmt"
	"net"
)

func Run(conn net.Conn) error {
	fmt.Println(">>> start")

	reader := bufio.NewReader(conn)
	req := NewRequester(reader)

	if err := req.ReadHeader(); err != nil {
		return err
	}

	// リクエストボディ
	method, err := req.GetMethod()
	if err != nil {
		return err
	}
	// body読み込む
	if method == "POST" || method == "PUT" {
		req.ReadBody(reader)
	}

	// status line
	res := NewResponser(conn)
	res.WriteOK()

	fmt.Println("<<< end")

	return nil
}
