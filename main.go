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

	scanner := bufio.NewScanner(conn)

	// 1行ずつ処理する
	for scanner.Scan() {
		// headerとbodyの間の空行があるのでheaderだけを読み取ることになる
		if scanner.Text() == "" {
			break
		}
		fmt.Println(scanner.Text())
	}

	if scanner.Err() != nil {
		return scanner.Err()
	}

	fmt.Println("<<< end")

	return nil
}
