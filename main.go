package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/textproto"
	"strconv"
	"strings"

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
	// headerとbodyを分けて処理している関係でバッファがなくなる
	// textprotoを利用してreaderを生成しheaderを読み取る
	scanner := textproto.NewReader(reader)

	var contentLength int

	// 1行ずつ処理する
	for {
		line, err := scanner.ReadLine()
		// headerとbodyの間の空行があるのでheaderだけを読み取ることになる
		if line == "" {
			break
		}
		if err != nil {
			return errors.WithStack(err)
		}

		// ex) line = "Content-Length: 23"
		if strings.HasPrefix(line, "Content-Length") {
			contentLength, err = strconv.Atoi(strings.TrimSpace(strings.Split(line, ":")[1]))
			if err != nil {
				return errors.WithStack(err)
			}
		}

		fmt.Println(line)
	}

	fmt.Println("read body")

	// リクエストボディ
	// Content-Lengthに指定された分のbodyがある
	buf := make([]byte, contentLength)
	_, err = io.ReadFull(reader, buf)
	if err != nil {
		return errors.WithStack(err)
	}

	fmt.Println("BODY:", string(buf))

	fmt.Println("<<< end")

	return nil
}
