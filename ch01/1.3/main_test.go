package main

import (
	"strings"
	"testing"
)

func TestEcho(t *testing.T) {
	// t.Fatal("not implemented")
	for _, tc := range testCases {
		if got := tc.proc(tc.in...); got != tc.want {
			t.Fatalf("wanted %q, but got %q", tc.want, got)
		}
	}
}

func BenchmarkSpeed(b *testing.B) {
	for _, bm := range testCases {
		b.Run(bm.name, func(b *testing.B) {
			b.ReportAllocs()
			proc, in := bm.proc, bm.in // capture
			for i := 0; i < b.N; i++ {
				proc(in...)
			}
		})
	}
}

var testCases = []struct {
	name string
	proc func(...string) string
	in   []string
	want string
}{
	{
		"Slow",
		slow,
		[]string{"Hello", "World"},
		"Hello World",
	},
	{
		"Fast",
		fast,
		[]string{"Hello", "World"},
		"Hello World",
	},
}

func slow(args ...string) string {
	res, sep := "", ""
	for _, v := range args {
		res += sep + v
		sep = " "
	}
	return res
}

func fast(args ...string) string {
	return strings.Join(args, " ")
}
