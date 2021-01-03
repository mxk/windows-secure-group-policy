package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mxk/go-cli"
)

var _ = cli.Main.Add(&cli.Cfg{
	Name:    "polfix",
	Usage:   "<registry.pol> ...",
	Summary: `Change pol file encoding for "MULTISZ:\0" to match gpedit`,
	MinArgs: 1,
	New:     func() cli.Cmd { return polfixCmd{} },
})

type polfixCmd struct{}

// 2 instead of 6 null bytes for data
var (
	multiszLGPO   = []byte{0x3B, 0x00, 0x07, 0x00, 0x00, 0x00, 0x3B, 0x00, 0x06, 0x00, 0x00, 0x00, 0x3B, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x5D, 0x00}
	multiszGPEdit = []byte{0x3B, 0x00, 0x07, 0x00, 0x00, 0x00, 0x3B, 0x00, 0x02, 0x00, 0x00, 0x00, 0x3B, 0x00, 0x00, 0x00, 0x5D, 0x00}
)

func (polfixCmd) Main(files []string) error {
	for _, name := range files {
		b, err := ioutil.ReadFile(name)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", name, err)
		}
		n := len(b)
		if b = bytes.Replace(b, multiszLGPO, multiszGPEdit, -1); len(b) == n {
			continue
		}
		fmt.Println("fixed:", name)
		err = cli.WriteFileAtomic(name, func(f *os.File) error {
			_, err := f.Write(b)
			return err
		})
		if err != nil {
			return fmt.Errorf("failed to write %s: %w", name, err)
		}
	}
	return nil
}
