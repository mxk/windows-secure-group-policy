//
// Written by Maxim Khitrov (May 2017)
//

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) < 2 || os.Args[1] == "/?" || os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Fprintf(os.Stderr, "Usage: %s <src> [<dst>]\n", os.Args[0])
		os.Exit(2)
	}
	src := os.Args[1]
	dst := src
	if len(os.Args) > 2 {
		dst = os.Args[2]
	}
	pol, err := ioutil.ReadFile(src)
	if err != nil {
		panic(err)
	}
	// Change encoding for "MULTISZ:\0" action lines to match gpedit (2 instead of 6 null bytes for data)
	old := []byte{0x3B, 0x00, 0x07, 0x00, 0x00, 0x00, 0x3B, 0x00, 0x06, 0x00, 0x00, 0x00, 0x3B, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x5D, 0x00}
	new := []byte{0x3B, 0x00, 0x07, 0x00, 0x00, 0x00, 0x3B, 0x00, 0x02, 0x00, 0x00, 0x00, 0x3B, 0x00, 0x00, 0x00, 0x5D, 0x00}
	pol = bytes.Replace(pol, old, new, -1)
	if err = ioutil.WriteFile(dst, pol, 0666); err != nil {
		panic(err)
	}
}
