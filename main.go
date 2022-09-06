package main

import (
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/TakuroSugahara/http-server/server"
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
	listen, err := net.Listen("tcp", "localhost:12345")
	if err != nil {
		errors.WithStack(err)
	}
	defer listen.Close()

	conn, err := listen.Accept()
	if err != nil {
		return errors.WithStack(err)
	}
	fmt.Println("Start server listen...")

	wg := sync.WaitGroup{}

	wg.Add(1)
	// goroutineを使って複数リクエストをさばけるようにする
	go func(conn net.Conn) {
		defer wg.Done()
		defer conn.Close()

		fmt.Println("accept")
		if err := server.Run(conn); err != nil {
			log.Printf("internal server error: %+v", err)
		}
	}(conn)

	wg.Wait()

	return nil
}
