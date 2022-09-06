package main

import (
	"bufio"
	"fmt"
	"io"
	"net/textproto"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

const (
	Method = "Method"
	Path   = "Path"
)

type requester struct {
	scanner *textproto.Reader
	header  map[string]string
}

func NewRequester(r *bufio.Reader) *requester {
	return &requester{
		// headerとbodyを分けて処理している関係でバッファがなくなる
		// textprotoを利用してreaderを生成しheaderを読み取る
		scanner: textproto.NewReader(r),
		header:  make(map[string]string),
	}
}

// 1行ずつ処理する
func (req *requester) ReadHeader() error {
	isFirst := true
	for {
		line, err := req.scanner.ReadLine()
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
			req.header[Method] = headerLine[0]
			req.header[Path] = headerLine[1]
			continue
		}

		// Header Fields
		headerFields := strings.SplitN(line, ": ", 2)
		fmt.Printf("%s: %s\n", headerFields[0], headerFields[1])
		req.header[headerFields[0]] = headerFields[1]
	}
	return nil
}

func (r *requester) GetMethod() (string, error) {
	method, ok := r.header[Method]
	if !ok {
		return "", errors.New("no method found")
	}
	return method, nil
}

func (r *requester) GetPath() string {
	path, ok := r.header[Path]
	if !ok {
		return "/"
	}
	return path
}

func (r *requester) getContentLength() (int, error) {
	len, err := strconv.Atoi(r.header["Content-Length"])
	if err != nil {
		return 0, errors.WithStack(err)
	}
	return len, nil
}

func (r *requester) ReadBody(reader *bufio.Reader) error {
	len, err := r.getContentLength()
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
	return nil
}
