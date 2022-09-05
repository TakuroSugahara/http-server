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

	// 1行ずつ処理する
	var method, path string
	header := make(map[string]string)

	isFirst := true
	for {
		line, err := scanner.ReadLine()
		// headerとbodyの間の空行があるのでheaderだけを読み取ることになる
		if line == "" {
			break
		}
		if err != nil {
			return errors.WithStack(err)
		}

		if isFirst {
			isFirst = false
			// headerの1行目はmethod, pathなどを表している
			// 空白で分けているので分割
			headerLine := strings.Fields(line)
			header["Method"] = headerLine[0]
			header["Path"] = headerLine[1]
			method = headerLine[0]
			path = headerLine[1]
			fmt.Println(method, path)
			continue
		}

		// Header Fields
		headerFields := strings.SplitN(line, ": ", 2)
		fmt.Printf("%s: %s\n", headerFields[0], headerFields[1])
		header[headerFields[0]] = headerFields[1]
	}

	// リクエストボディ
	method, ok := header["Method"]
	if !ok {
		return errors.New("no method found")
	}
	// body読み込む
	if method == "POST" || method == "PUT" {
		len, err := strconv.Atoi(header["Content-Length"])
		if err != nil {
			return errors.WithStack(err)
		}
		// Content-Lengthに指定された分のbodyがある
		buf := make([]byte, len)
		_, err = io.ReadFull(reader, buf)
		if err != nil {
			return errors.WithStack(err)
		}
		fmt.Println("BODY:", string(buf))
	}

	// status line
	io.WriteString(conn, "HTTP/1.1 200 OK \r\n")

	// header
	io.WriteString(conn, "Content-Type: text/html\r\n")
	io.WriteString(conn, "\r\n")
	io.WriteString(conn, "<h1>Hello World!</h1>")

	fmt.Println("<<< end")

	return nil
}
