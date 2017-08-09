package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type info struct {
	names []string
	count int
}

type input struct {
	name string
	data []byte
}

func main() {
	prog := filepath.Base(os.Args[0])
	counts := make(map[string]info)
	inputs := make([]input, 1)

	// Collect all input data
	if len(os.Args[1:]) > 0 {
		for _, fname := range os.Args[1:] {
			data, err := ioutil.ReadFile(fname)
			if err != nil {
				log.Printf("%s: failed to %v", prog, err)
				continue
			}
			inputs = append(inputs, input{filepath.Base(fname), data})
		}
	} else {
		data, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal("%s: failed to %v", prog, err)
		}
		inputs = []input{input{"/dev/stdin", data}}
	}

	// Collect stats
	for _, input := range inputs {
		countLines(input, counts)
	}

	// Display duplicates
	for line, info := range counts {
		if info.count > 1 {
			fmt.Printf("%s(%d): %s\n", strings.Join(info.names, ","), info.count, line)
		}
	}
}

func countLines(in input, counts map[string]info) {
	// input := bufio.NewScanner(r)
	input := bytes.Split(in.data, []byte("\n"))
	for _, text := range input {
		// Get text
		line := string(text)
		// Skip if text is empty
		if len(text) == 0 {
			continue
		}
		// Get info from the map
		info := counts[line]
		// Update info
		info.names = updateNames(info.names, in.name)
		info.count++
		// Update map
		counts[line] = info
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
