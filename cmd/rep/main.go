package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

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
	Use:     "rep",
	Short:   "rep is repeating input stream",
	Version: version,
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		flags := cmd.Flags()

		delimiter, err := flags.GetString("delimiter")
		if err != nil {
			panic(err)
		}

		useStdin, err := flags.GetBool("stdin")
		if err != nil {
			panic(err)
		}

		opts := options{delimiter: delimiter, useStdin: useStdin}

		if opts.useStdin {
			// 標準出力から繰り返す文字を取得
			sc := bufio.NewScanner(os.Stdin)
			sc.Scan()
			if err := sc.Err(); err != nil {
				panic(err)
			}
			word := sc.Text()

			printRepeatedLine(args, word, opts)
			return
		}

		// 引数の最後の文字を繰り返す文字として取得
		word := args[len(args)-1]
		args = args[:len(args)-1]
		printRepeatedLine(args, word, opts)
	},
}

func printRepeatedLine(repeatStrs []string, word string, opts options) {
	for _, s := range repeatStrs {
		cnt, err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
		var words []string
		for i := 0; i < cnt; i++ {
			words = append(words, word)
		}
		line := strings.Join(words, opts.delimiter)
		fmt.Println(line)
	}
}
