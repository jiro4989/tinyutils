package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/docopt/docopt-go"
)

type (
	Config struct {
		Files []string
	}
)

const (
	version = `v1.0.0
Copyright (C) 2019, jiro4989
Released under the MIT License.
https://github.com/jiro4989/tinyutils`

	doc = `codepoint prints codepoint of text.

Usage:
	codepoint [options]
	codepoint [options] <files>...
	codepoint -h | --help
	codepoint -v | --version

Options:
	-h --help                     Show this screen.
	-v --version                  Show version.`
)

func main() {
	parser := &docopt.Parser{}
	args, _ := parser.ParseArgs(doc, nil, version)
	config := Config{}
	err := args.Bind(&config)
	if err != nil {
		panic(err)
	}

	fmt.Println("char code_point code_point(hex)")

	// 引数が一つでもあれば引数のファイルを処理
	if 0 < len(config.Files) {
		for _, file := range config.Files {
			func() {
				f, err := os.Open(file)
				if err != nil {
					panic(err)
				}
				defer f.Close()
				if err := printCodePoints(f); err != nil {
					panic(err)
				}
			}()
		}

		return
	}
	// 引数が一つもなければ標準入力を処理
	if err := printCodePoints(os.Stdin); err != nil {
		panic(err)
	}
}

func printCodePoints(f *os.File) error {
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		for _, oct := range []rune(line) {
			ch := string(oct)
			hex := strconv.FormatInt(int64(oct), 16)
			line := fmt.Sprintf("%s %d \\U%s", ch, oct, hex)

			fmt.Println(line)
		}
	}
	if err := sc.Err(); err != nil {
		return err
	}
	return nil
}
