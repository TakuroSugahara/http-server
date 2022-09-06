package main

import (
	"bufio"
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

	listen, err := net.Listen("tcp", "localhost:12345")
	if err != nil {
		errors.WithStack(err)
	}
	defer listen.Close()

	conn, err := listen.Accept()
	if err != nil {
		return errors.WithStack(err)
	}
	defer conn.Close()

	fmt.Println(">>> start")

	reader := bufio.NewReader(conn)

	req := NewRequester(reader)

	// TODO: goroutineを使って複数リクエストに対応できるようにする

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
