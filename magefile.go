// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
)

var (
	defaultBuildFlags = []string{`-ldflags`, `-s -w -extldflags "-static"`}
	xbuildTarget      = []string{"windows", "linux", "darwin", "386", "amd64"}
	commands          = []string{"flat", "rep", "ucut", "codepoint"}
)

// Default target to run when none is specified
// If not set, running mage will list available targets
// var Default = Build

// A build step that requires additional params, or platform specific steps for example
func Build() error {
	fmt.Println("Building...")
	os.Setenv("GO111MODULE", "on")
	for _, cmd := range commands {
		fmt.Printf("Build %s...\n", cmd)
		args := []string{"build"}
		args = append(args, defaultBuildFlags...)
		args = append(args, "-o", "bin/"+cmd, "./cmd/"+cmd)
		cmd := exec.Command("go", args...)
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}

// A build step that requires additional params, or platform specific steps for example
func Xbuild() error {
	fmt.Println("Building...")
	cmd := exec.Command("go", "build", "-o", "MyApp", ".")
	return cmd.Run()
}

// A custom install step if you need your bin someplace other than go/bin
func Install() error {
	mg.Deps(Build)
	fmt.Println("Installing...")
	return os.Rename("./MyApp", "/usr/bin/MyApp")
}

// Clean up after yourself
func Clean() {
	fmt.Println("Cleaning...")
	os.RemoveAll("MyApp")
}

// Clean up after yourself
func Test() {
	fmt.Println("Cleaning...")
	os.RemoveAll("MyApp")
}

// Clean up after yourself
func Bootstrap() {
	fmt.Println("Cleaning...")
	os.RemoveAll("MyApp")
}
