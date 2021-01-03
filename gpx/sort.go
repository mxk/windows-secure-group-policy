package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	uni "unicode"

	"github.com/mxk/go-cli"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

var _ = cli.Main.Add(&cli.Cfg{
	Name:    "sort",
	Usage:   "<lgpo.txt> ...",
	Summary: "Sort LGPO text files",
	MinArgs: 1,
	New:     func() cli.Cmd { return sortCmd{} },
})

type sortCmd struct{}

func (sortCmd) Main(files []string) error {
	for _, name := range files {
		p, err := ParseLGPO(name)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", name, err)
		}
		if sort.SliceIsSorted(p.Entries, p.less) {
			continue
		}
		sort.Slice(p.Entries, p.less)
		fmt.Println("sorted:", name)
		if err = cli.WriteFileAtomic(name, p.Write); err != nil {
			return fmt.Errorf("failed to write %s: %w", name, err)
		}
	}
	return nil
}

type Entry struct{ Comment, Conf, Key, Value, Action, sortKey string }

type LGPO struct {
	Header  string
	Entries []*Entry
	Footer  string
}

const CRLF = "\r\n"

func ParseLGPO(name string) (*LGPO, error) {
	b, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	t := unicode.BOMOverride(charmap.Windows1252.NewDecoder().Transformer)
	if b, _, err = transform.Bytes(t, b); err != nil {
		return nil, err
	}
	s := string(b)
	lines := strings.Split(s, CRLF)
	if len(lines) == 1 {
		lines = strings.Split(s, "\n")
	}
	comment := 0
	entries := make([]*Entry, 0, len(lines)/5)
	isEntry := func(ln string) bool { return len(ln) > 0 && ln[0] != ';' }
	for i := 0; i < len(lines); i++ {
		ln := lines[i]
		if !isEntry(ln) {
			if len(ln) == 0 && comment == i {
				comment++ // Skip leading blank lines
			}
			continue
		}
		if !(strings.EqualFold(ln, "Computer") || strings.EqualFold(ln, "User") ||
			(len(ln) > 5 && strings.EqualFold(ln[:5], "User:"))) {
			return nil, fmt.Errorf("invalid configuration on line %d: %q", i, ln)
		}
		if len(lines)-i < 4 {
			return nil, fmt.Errorf("incomplete entry on line %d", i)
		}
		e := &Entry{
			Comment: strings.Join(lines[comment:i], CRLF),
			Conf:    ln,
			Key:     lines[i+1],
			Value:   lines[i+2],
			Action:  lines[i+3],
		}
		if !(isEntry(e.Key) && isEntry(e.Value) && isEntry(e.Action)) {
			return nil, fmt.Errorf("invalid entry on line %d: %+v", i, e)
		}
		e.sortKey = fmt.Sprintf("%s\x00\x00%s\x00\x00%s",
			strings.ToUpper(e.Conf),
			strings.Map(func(r rune) rune {
				if r == '\\' || r == '/' {
					return 0
				}
				return uni.ToUpper(r)
			}, e.Key),
			strings.ToUpper(e.Value),
		)
		entries = append(entries, e)
		i, comment = i+3, i+4
	}
	p := &LGPO{
		Entries: entries,
		Footer:  strings.Join(lines[comment:], CRLF),
	}
	if len(entries) > 0 {
		p.Header, entries[0].Comment = entries[0].Comment, ""
	} else {
		p.Header, p.Footer = p.Footer, ""
	}
	return p, nil
}

func (p *LGPO) Write(f *os.File) error {
	w := bufio.NewWriter(f)
	writeln := func(s string) { w.WriteString(s); w.WriteString(CRLF) }
	if p.Header != "" {
		w.WriteString(p.Header)
	}
	for _, e := range p.Entries {
		if w.WriteString(CRLF); e.Comment != "" {
			writeln(e.Comment)
		}
		writeln(e.Conf)
		writeln(e.Key)
		writeln(e.Value)
		writeln(e.Action)
	}
	if p.Footer != "" {
		w.WriteString(CRLF)
		w.WriteString(p.Footer)
	}
	return w.Flush()
}

func (p *LGPO) less(i, j int) bool {
	return p.Entries[i].sortKey < p.Entries[j].sortKey
}
