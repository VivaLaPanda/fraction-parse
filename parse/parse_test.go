package parse

import (
	"github.com/VivaLaPanda/fraction-parse/types"
	"sync"
	"testing"
)

var parseTests = []struct {
	stringToParse string
}{
	{"18"},
	{"18.5"},
	{"18_3/5"},
	{"3/5"},
	{"-0.5"},
}

func TestStartParseWorker(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	parsedFractions := make(chan types.Fraction, 100)
	fracStrings := make(chan string, 100)
	for _, test := range parseTests {
		fracStrings <- test.stringToParse
	}
	StartParseWorker(0, fracStrings, parsedFractions, wg)
	close(fracStrings)
}
