//
// Written by Maxim Khitrov (May 2017)
//

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Change encoding for "MULTISZ:\0" action lines to match gpedit (2 instead of 6
// null bytes for data)
var (
	old = []byte{0x3B, 0x00, 0x07, 0x00, 0x00, 0x00, 0x3B, 0x00, 0x06, 0x00, 0x00, 0x00, 0x3B, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x5D, 0x00}
	new = []byte{0x3B, 0x00, 0x07, 0x00, 0x00, 0x00, 0x3B, 0x00, 0x02, 0x00, 0x00, 0x00, 0x3B, 0x00, 0x00, 0x00, 0x5D, 0x00}
)

func main() {
	if len(os.Args) < 2 || os.Args[1] == "/?" || os.Args[1] == "-h" ||
		os.Args[1] == "--help" {
		fail(2, "Usage: %s <registry.pol ...>\n", os.Args[0])
	}
	for _, name := range os.Args[1:] {
		b, err := ioutil.ReadFile(name)
		if err != nil {
			fail(1, "error reading %s: %v\n", name, err)
		}
		n := len(b)
		if b = bytes.Replace(b, old, new, -1); len(b) != n {
			fmt.Println("file fixed:", name)
			if err = write(name, b); err != nil {
				fail(1, "error writing %s: %v\n", name, err)
			}
		}
	}
}

// fail writes an error message to stderr and terminates the process with the
// specified status code.
func fail(code int, format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(code)
}

// write writes b to the specified file.
func write(name string, b []byte) (err error) {
	tmp, err := ioutil.TempFile(filepath.Dir(name), filepath.Base(name)+".")
	if err != nil {
		return
	}
	panicked := true
	defer func() {
		if panicked || err != nil {
			tmp.Close()
			os.Remove(tmp.Name())
		}
	}()
	if _, err = tmp.Write(b); err == nil {
		if err = tmp.Close(); err == nil {
			err = os.Rename(tmp.Name(), name)
		}
	}
	panicked = false
	return
}
