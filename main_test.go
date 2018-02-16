package main

import (
	"os"
	"testing"
)

func TestMain(t *testing.T) {
	os.Args = make([]string, 2)
	os.Args[1] = "-filepath=test-data/fractions.txt"
	main()
}

func BenchMain(b *testing.B) {
	os.Args = make([]string, 2)
	os.Args[1] = "-filepath=test-data/fractions.txt"
	for n := 0; n < b.N; n++ {
		main()
	}
}
