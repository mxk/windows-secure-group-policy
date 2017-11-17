//
// Written by Maxim Khitrov (November 2017)
//

package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"
)

const NL = "\r\n"

func main() {
	if len(os.Args) < 2 || os.Args[1] == "/?" || os.Args[1] == "-h" ||
		os.Args[1] == "--help" {
		fail(2, "Usage: %s <lgpo.txt ...>\n", os.Args[0])
	}
	cache := Cache{make(map[string]string), make(map[string]string)}
	for _, name := range os.Args[1:] {
		entries, header, footer, err := Parse(name, &cache)
		if err != nil {
			fail(1, "error reading %s: %v\n", name, err)
		}
		if es := (EntrySorter{entries, &cache}); !sort.IsSorted(&es) {
			sort.Sort(&es) // Unstable ok, entries can't be equal
			fmt.Println("file sorted:", name)
			if err = Write(name, entries, header, footer); err != nil {
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

// Parse extracts the entries and header/footer comments from an LGPO text file.
func Parse(name string, c *Cache) (e []*Entry, header, footer string, err error) {
	f, err := os.Open(name)
	if err != nil {
		return
	}
	defer f.Close()

	// Detect Unicode BOM
	r := bufio.NewReader(f)
	if b, _ := r.Peek(2); len(b) == 2 && b[0]&b[1] == 0xFE && b[0] != b[1] {
		err = fmt.Errorf("unicode files not supported")
		return
	}

	// Parse entries
	confs := []string{
		"Computer",
		"User",
		"User:Administrators",
		"User:Non-Administrators",
	}
	conf := confs[0]
	scan := bufio.NewScanner(r)
	line := 0
	var comment []byte
	for scan.Scan() {
		ln := scan.Bytes()
		line++

		// Comment (leading blank lines are ignored)
		if len(ln) == 0 || ln[0] == ';' {
			if comment = append(comment, ln...); len(comment) > 0 {
				comment = append(comment, NL...)
			}
			continue
		} else if len(e) == 0 {
			header, comment = string(comment), comment[:0]
		}

		// Configuration type
		if string(ln) != conf {
			for _, conf = range confs {
				if string(ln) == conf {
					goto kva
				}
			}
			if len(ln) <= 5 || string(ln[:5]) != "User:" {
				err = fmt.Errorf("invalid entry on line %d: %q", line, ln)
				return
			}
			conf = string(ln)
			confs = append(confs, conf)
		}

	kva:
		// Key, value, action
		var kva [3]string
		for i := range kva {
			if scan.Scan() {
				ln = scan.Bytes()
				if line++; len(ln) > 0 && ln[0] != ';' {
					if i != 1 {
						kva[i] = c.String(ln)
					} else {
						kva[i] = string(ln) // Bypass cache for values
					}
					continue
				}
				err = fmt.Errorf("unexpected blank line or comment")
			} else if err = scan.Err(); err == nil {
				err = io.EOF
			}
			err = fmt.Errorf("incomplete entry on line %d: %v", line, err)
			return
		}
		e = append(e, &Entry{
			Comment: string(comment),
			Conf:    conf,
			Key:     kva[0],
			Value:   kva[1],
			Action:  kva[2],
		})
		comment = comment[:0]
	}
	if err = scan.Err(); err != nil {
		err = fmt.Errorf("error on line %d: %v", line, err)
	} else {
		footer = string(comment)
	}
	return
}

// Write saves entries e to the named file.
func Write(name string, e []*Entry, header, footer string) (err error) {
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
	w := bufio.NewWriter(tmp)
	w.WriteString(header)
	for i, n := range e {
		if n.Comment != "" {
			w.WriteString(n.Comment)
		}
		w.WriteString(n.Conf)
		w.WriteString(NL)
		w.WriteString(n.Key)
		w.WriteString(NL)
		w.WriteString(n.Value)
		w.WriteString(NL)
		w.WriteString(n.Action)
		w.WriteString(NL)
		if i+1 < len(e) || footer != "" {
			if _, err = w.WriteString(NL); err != nil {
				return
			}
		}
	}
	w.WriteString(footer)
	if err = w.Flush(); err == nil {
		if err = tmp.Close(); err == nil {
			err = os.Rename(tmp.Name(), name)
		}
	}
	panicked = false
	return
}

// Entry represents a single four-line LGPO text entry with an optional comment.
type Entry struct{ Comment, Conf, Key, Value, Action string }

// Less returns true if e should come before x when sorted.
func (e *Entry) Less(x *Entry, c *Cache) bool {
	if e.Conf != x.Conf {
		return e.Conf < x.Conf
	}
	if ek, xk := c.Norm(e.Key), c.Norm(x.Key); ek != xk {
		return ek < xk
	}
	if ev, xv := c.Norm(e.Value), c.Norm(x.Value); ev != xv {
		return ev < xv
	}
	panic("duplicate entry: " + e.Key + "!" + e.Value)
}

// Cache allows string reuse to reduce memory usage and speed up comparisons.
type Cache struct{ str, norm map[string]string }

// String returns string(b).
func (c *Cache) String(b []byte) string {
	s, ok := c.str[string(b)]
	if !ok {
		s = string(b)
		c.str[s] = s
	}
	return s
}

// Norm converts k to a (somewhat) normalized form for comparison.
func (c *Cache) Norm(k string) string {
	s, ok := c.norm[k]
	if !ok {
		s = strings.Map(func(r rune) rune {
			if r == '\\' {
				return 0
			}
			return unicode.ToUpper(r)
		}, k)
		c.norm[k] = s
	}
	return s
}

// EntrySorter implements sort.Interface for case-insensitive sorting of entries
// in an LGPO text file.
type EntrySorter struct {
	e []*Entry
	c *Cache
}

func (s *EntrySorter) Len() int           { return len(s.e) }
func (s *EntrySorter) Swap(i, j int)      { s.e[i], s.e[j] = s.e[j], s.e[i] }
func (s *EntrySorter) Less(i, j int) bool { return s.e[i].Less(s.e[j], s.c) }
