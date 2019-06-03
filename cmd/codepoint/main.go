package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

// options ... main で受け取られる引数、オプション
type options struct {
	delimiter string
	useStdin  bool
}

const (
	version = `v1.0.0
Copyright (C) 2019, jiro4989
Released under the MIT License.
https://github.com/jiro4989/tinyutils`
)

func init() {
	cobra.OnInitialize()
	command.Flags().StringP("delimiter", "d", "", "delimiter")
	command.Flags().BoolP("stdin", "i", false, "use stdin")
}

func main() {
	if err := command.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var command = &cobra.Command{
	Use:     "codepoint",
	Short:   "codepoint prints codepoint of text",
	Version: version,
	Run: func(cmd *cobra.Command, args []string) {
		// flags := cmd.Flags()
		//
		// delimiter, err := flags.GetString("delimiter")
		// if err != nil {
		// 	panic(err)
		// }
		//
		// useStdin, err := flags.GetBool("stdin")
		// if err != nil {
		// 	panic(err)
		// }
		//
		// opts := options{delimiter: delimiter, useStdin: useStdin}

		fmt.Println("char code_point code_point(hex)")

		// 引数が一つでもあれば引数のファイルを処理
		if 0 < len(args) {
			for _, file := range args {
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
	},
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
