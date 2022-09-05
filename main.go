package main

import (
	"fmt"
	"net"

	"github.com/pkg/errors"
)

// http serverやること
// - クライアントからの接続を待ち受ける
// - クライアントから送信された HTTP リクエストをパースする
// - HTTP リクエストに基づいて HTTP レスポンスを生成/返却する

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
	}
}

func run() error {
	fmt.Println("Start tcp listen...")

	fmt.Println("create port")
	listen, err := net.Listen("tcp", "localhost:12345")
	if err != nil {
		errors.WithStack(err)
	}
	defer listen.Close()

	fmt.Println("create connection")
	conn, err := listen.Accept()
	if err != nil {
		return errors.WithStack(err)
	}
	defer conn.Close()

	fmt.Println(">>> start")

	// 受け取り用バッファ
	buf := make([]byte, 1024)

	for {
		n, err := conn.Read(buf)
		if n == 0 {
			break
		}
		if err != nil {
			return errors.WithStack(err)
		}
		fmt.Println(string(buf[:n]))
	}

	fmt.Println("<<< end")

	return nil
}
