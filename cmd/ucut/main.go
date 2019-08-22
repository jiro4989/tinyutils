package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/docopt/docopt-go"
)

type (
	Config struct {
		Files           []string
		Delimiter       string
		OutputDelimiter string `docopt:"--output-delimiter"`
		Fields          string
		OutFile         string `docopt:"--out-file"`
	}
)

const (
	version = `v1.0.0
Copyright (C) 2019, jiro4989
Released under the MIT License.
https://github.com/jiro4989/tinyutils`

	doc = `ucut cuts word with unicode character.

Usage:
	ucut [options]
	ucut [options] <files>...
	ucut -h | --help
	ucut -v | --version

Options:
	-h --help                            Show this screen.
	-v --version                         Show version.
	-d --delimiter=<DELIMITER>           Field delimiter. [default:  ]
	-D --output-delimiter=<DELIMITER>    Output field delimiter. [default:  ]
	-f --fields=<FIELDS>                 Fields. [default: -]
	-o --out-file=<OUTFILE>              Out file.`
)

func main() {
	parser := &docopt.Parser{}
	args, _ := parser.ParseArgs(doc, nil, version)
	config := Config{}
	err := args.Bind(&config)
	if err != nil {
		panic(err)
	}

	// 出力先ファイルの指定があればファイルに書き込む
	// なければ標準出力
	var outfile *os.File
	if config.OutFile == "" {
		outfile = os.Stdout
	} else {
		outfile, err = os.OpenFile(config.OutFile, os.O_RDWR, os.ModePerm)
		if err != nil {
			panic(err)
		}

		defer outfile.Close()
	}

	// 引数がある場合はそれをファイルとして処理
	if 0 < len(config.Files) {
		for _, file := range config.Files {
			func() {
				f, err := os.Open(file)
				if err != nil {
					panic(err)
				}
				defer f.Close()
				if err := ucutIO(outfile, f, config); err != nil {
					panic(err)
				}
			}()
		}

		return
	}
	// 引数がない場合は標準入力を受け取る
	if err := ucutIO(outfile, os.Stdin, config); err != nil {
		panic(err)
	}
}

func ucutIO(dst io.Writer, src io.Reader, config Config) error {
	sc := bufio.NewScanner(src)
	for sc.Scan() {
		line := sc.Text()
		fields, err := ucut(line, config.Delimiter, config.Fields)
		if err != nil {
			return err
		}
		joined := strings.Join(fields, config.OutputDelimiter)
		fmt.Fprintln(dst, joined)
	}
	if err := sc.Err(); err != nil {
		return err
	}
	return nil
}

// ucut はlineをdelimiterで区切り、fieldsの指定のフィールドを指定の順序で取得し
// て返す
func ucut(line, delimiter, fields string) ([]string, error) {
	lineFields := strings.Split(line, delimiter)
	fieldFields := strings.Split(fields, ",")

	// どのフィールドを返すか？を決める
	var fs []int
	for _, v := range fieldFields {
		se := strings.Split(v, "-")
		// "-" で分割指定をしているときは範囲でフィールド位置を取得する
		if 1 < len(se) {
			var (
				startNum, endNum int
				err              error
			)
			// 開始のフィールド位置を取得
			// 未指定の場合は先頭から
			if se[0] == "" {
				startNum = 1
			} else {
				startNum, err = strconv.Atoi(se[0])
				if err != nil {
					return nil, err
				}
			}

			// 終了のフィールド位置を取得
			// 未指定の場合は最後まで
			if se[1] == "" {
				endNum = len(lineFields)
			} else {
				endNum, err = strconv.Atoi(se[1])
				if err != nil {
					return nil, err
				}
			}
			for i := startNum - 1; i < endNum; i++ {
				fs = append(fs, i)
			}
			continue
		}
		n, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		// フィールドの指定はcutに習って1から始まるので1減らす
		if n < 1 {
			msg := fmt.Sprintf("fields index must be over 1. index = %d", n)
			return nil, errors.New(msg)
		}
		fs = append(fs, n-1)
	}

	var ret []string
	for _, v := range fs {
		ret = append(ret, lineFields[v])
	}
	return ret, nil
}
