package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	prog := filepath.Base(os.Args[0]) + ":"
	fmt.Println(prog, strings.Join(os.Args[1:], " "))
}
