package main

import "github.com/mxk/go-cli"

//go:generate go build -trimpath "-ldflags=-s -w" -o ../Tools

func main() {
	cli.Main.Summary = "LGPO extensions"
	cli.Main.Run()
}
