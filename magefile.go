// +build mage

package main

import (
	"fmt"
	"os"

	"github.com/jiro4989/magetask/v1"
	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
)

var (
	commands = []string{"flat", "rep", "ucut", "codepoint"}
	config   = magetask.BuildConfig{}
)

// Build はパッケージをビルドして所定のディレクトリ配下に配置する。
func Build() error {
	for _, cmd := range commands {
		fmt.Printf("Building %s ... ", cmd)
		if err := magetask.Build("bin/"+cmd, "./cmd/"+cmd, config); err != nil {
			return err
		}
		fmt.Println("OK")
	}
	return nil
}

// Xbuild はパッケージをクロスコンパイルして所定のディレクトリ配下に配置する。
func Xbuild() error {
	for _, cmd := range commands {
		fmt.Printf("Xbuilding %s ... ", cmd)
		if err := magetask.Xbuild(cmd, "./cmd/"+cmd, config); err != nil {
			return err
		}
		fmt.Println("OK")
	}
	return nil
}

// A custom install step if you need your bin someplace other than go/bin
func Install() error {
	mg.Deps(Build)
	fmt.Println("Installing...")
	return os.Rename("bin/", "/usr/bin/MyApp")
}

// Clean up after yourself
func Clean() {
	fmt.Print("Cleaning ... ")
	os.RemoveAll("bin")
	os.RemoveAll("dist")
	fmt.Println("OK")
}

// Test はテストを実行する。
func Test() error {
	return magetask.Test()
}
