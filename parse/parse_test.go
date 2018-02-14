package parse

import (
	"fmt"
	"sync"
	"testing"

	"github.com/VivaLaPanda/fraction-parse/types"
)

var parseTests = []struct {
	stringToParse string
	expected      string
}{
	{"18", "18_0/1"},
	{"18.5", "18_1/2"},
	{"18_3/5", "18_3/5"},
	{"-18_3/5", "-18_3/5"},
	{"3/5", "0_3/5"},
	{"-0.5", "-0_1/2"},
	{"-123", "-123_0/1"},
	{"+123", "123_0/1"},
}

func TestStartParseWorker(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	parsedFractions := make(chan types.Fraction, 100)
	fracStrings := make(chan string, 100)
	for _, test := range parseTests {
		fracStrings <- test.stringToParse
	}
	go StartParseWorker(0, fracStrings, parsedFractions, wg)
	close(fracStrings)
	go func() {
		wg.Wait()
		close(parsedFractions)
	}()

	for fraction := range parsedFractions {
		fmt.Printf("%s\n", fraction)

		found := false
		for _, test := range parseTests {
			if fraction.String() == test.expected {
				found = true
			}
		}

		if found == false {
			t.Errorf("Error occured while testing ParseWorker: '%v' was not a possible output.", fraction)
		}
	}
}
