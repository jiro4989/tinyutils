package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/docopt/docopt-go"
)

type (
	Config struct {
		Args        []string `docopt:"<args>"`
		ColumnCount int      `docopt:"--column-count"`
		Delimiter   string   `docopt:"--delimiter"`
	}
)

const (
	version = `v1.0.0
Copyright (C) 2019, jiro4989
Released under the MIT License.
https://github.com/jiro4989/tinyutils`

	doc = `flat is flatting input stream.
Usage:
	flat [options]
	flat [options] <args>...
	flat -h | --help
	flat -v | --version

Options:
	-h --help                     Show this screen.
	-v --version                  Show version.
	-n --column-count=<NUM>       Column count.
	-d --delimiter=<DELIMITER>    Field delimiter. [default:  ]`
)

func main() {
	parser := &docopt.Parser{}
	args, _ := parser.ParseArgs(doc, nil, version)
	config := Config{}
	err := args.Bind(&config)
	if err != nil {
		panic(err)
	}

	// 引数がある場合はそれをファイルとして処理
	if 0 < len(config.Args) {
		for _, file := range config.Args {
			func() {
				f, err := os.Open(file)
				if err != nil {
					panic(err)
				}
				defer f.Close()
				if err := writeFlat(os.Stdout, f, config); err != nil {
					panic(err)
				}
			}()
		}

		return
	}
	// 引数がない場合は標準入力を受け取る
	if err := writeFlat(os.Stdout, os.Stdin, config); err != nil {
		panic(err)
	}
}

func writeFlat(dst io.Writer, src io.Reader, config Config) error {
	var onelineDatas []string
	var i int
	// １行ずつ取得し１行分のデータに追加
	// カラム数を超過したら次の行のデータ用に切り替え
	sc := bufio.NewScanner(src)
	for sc.Scan() {
		line := sc.Text()
		onelineDatas = append(onelineDatas, line)
		i++
		if 0 < config.ColumnCount && config.ColumnCount <= i {
			s := strings.Join(onelineDatas, config.Delimiter)
			fmt.Fprintln(dst, s)
			onelineDatas = []string{}
			i = 0
		}
	}
	if err := sc.Err(); err != nil {
		return err
	}
	// 最後に残ったデータがあれば追加
	if 0 < len(onelineDatas) {
		s := strings.Join(onelineDatas, config.Delimiter)
		fmt.Fprintln(dst, s)
	}
	return nil
}
