package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// Config ... main で受け取られる引数、オプション
type Config struct {
	Args        []string
	ColumnCount int
	Delimiter   string
}

const (
	version = `v1.0.0
Copyright (C) 2019, jiro4989
Released under the MIT License.
https://github.com/jiro4989/tinyutils`
)

func init() {
	cobra.OnInitialize()
	RootCommand.Flags().IntP("columncount", "n", 0, "column count")
	RootCommand.Flags().StringP("delimiter", "d", " ", "delimiter")
}

var RootCommand = &cobra.Command{
	Use:     "flat",
	Short:   "flat is flatting input stream",
	Version: version,
	Run: func(cmd *cobra.Command, args []string) {
		flags := cmd.Flags()

		columnCount, err := flags.GetInt("columncount")
		if err != nil {
			panic(err)
		}

		delimiter, err := flags.GetString("delimiter")
		if err != nil {
			panic(err)
		}

		config := Config{Args: args, ColumnCount: columnCount, Delimiter: delimiter}

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
	},
}

func main() {
	if err := RootCommand.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
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
