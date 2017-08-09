package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type info struct {
	names []string
	count int
}

func main() {
	prog := filepath.Base(os.Args[0])
	counts := make(map[string]info)

	// Collect stats
	if len(os.Args[1:]) > 0 {
		for _, fname := range os.Args[1:] {
			f, err := os.Open(fname)
			if err != nil {
				log.Printf("%s: failed to %v", prog, err)
				continue
			}
			defer f.Close()
			countLines(f, fname, counts)
		}
	} else {
		countLines(os.Stdin, "/dev/stdin", counts)
	}

	// Display duplicates
	for line, info := range counts {
		if info.count > 1 {
			fmt.Printf("%s(%d): %s\n", strings.Join(info.names, ","), info.count, line)
		}
	}
}

func countLines(r io.Reader, name string, counts map[string]info) {
	input := bufio.NewScanner(r)
	for input.Scan() {
		if err := input.Err(); err != nil && err != io.EOF {
			log.Fatal(err)
		}

		// Get text
		text := input.Text()
		// Skip if text is empty
		if text == "" {
			continue
		}
		// Get info from the map
		info := counts[text]
		// Update info
		info.names = updateNames(info.names, filepath.Base(name))
		info.count++
		// Update map
		counts[text] = info
	}
}

func updateNames(names []string, name string) []string {
	for _, n := range names {
		if n == name {
			return names // already in the list
		}
	}
	return append(names, name)
}
