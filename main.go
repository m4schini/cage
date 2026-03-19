/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"cage/cmd"
	_ "embed"
	"fmt"
	"os"
)

//go:embed Containerfile
var Containerfile string

func main() {
	if len(os.Args) >= 2 && os.Args[1] == "Containerfile" {
		fmt.Println(Containerfile)
		return
	}
	cmd.Execute()
}
