package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/docopt/docopt-go"
)

type (
	Config struct {
		Args      []string
		Delimiter string
		Stdin     bool
	}
)

const (
	version = `v1.0.0
Copyright (C) 2019, jiro4989
Released under the MIT License.
https://github.com/jiro4989/tinyutils`

	doc = `rep is repeating input stream.

Usage:
	rep [options] <args>...
	rep -h | --help
	rep -v | --version

Options:
	-h --help                     Show this screen.
	-v --version                  Show version.
	-d --delimiter=<DELIMITER>    Field delimiter. [default:  ]
	-i --stdin                    Use stdin.`
)

func main() {
	parser := &docopt.Parser{}
	args, _ := parser.ParseArgs(doc, nil, version)
	config := Config{}
	err := args.Bind(&config)
	if err != nil {
		panic(err)
	}

	// 標準入力を使う指定がなければ
	// 引数の最後の文字を繰り返す文字として取得
	if !config.Stdin {
		args := config.Args
		word := args[len(args)-1]
		args = args[:len(args)-1]
		printRepeatedLine(args, word, config)
		return
	}

	// 標準入力を使う指定があれば
	// 標準入力から繰り返す文字を取得
	sc := bufio.NewScanner(os.Stdin)
	var repeat []string
	for sc.Scan() {
		s := sc.Text()
		repeat = append(repeat, s)
	}
	if err := sc.Err(); err != nil {
		panic(err)
	}
	printRepeatedLine(repeat, config.Args[0], config)
}

// printRepeatedLine は繰り返し文字列の回数だけwordを繰り返して出力する。
func printRepeatedLine(repeat []string, word string, config Config) {
	for _, s := range repeat {
		cnt, err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
		var words []string
		for i := 0; i < cnt; i++ {
			words = append(words, word)
		}
		line := strings.Join(words, config.Delimiter)
		fmt.Println(line)
	}
}
