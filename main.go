package main

import (
	"flag"
	"io/ioutil"
	"log"
	"strings"
	"sync"

	"github.com/VivaLaPanda/fraction-parse/parse"
	"github.com/VivaLaPanda/fraction-parse/types"
)

const numWorkers = 1

var filepath = flag.String("filepath", "", "File to pull fractions from")

func main() {
	// Open file
	flag.Parse()

	if *filepath == "" {
		log.Fatalf("You must specify a file to parse")
	}

	// Spawn worker to read in file, split on spaces, and push each
	// fraction into a buffered chan
	fracStrings := make(chan string, 100)
	go func() {
		fileToParse, err := ioutil.ReadFile(*filepath)
		check(err)
		tokensToParse := strings.Split(string(fileToParse), " ")
		for _, elem := range tokensToParse {
			fracStrings <- elem
		}

		close(fracStrings)
	}()

	// Spawn workers that take a buffered chan of fraction strings
	// and pushes them as Fraction types unto a results chan
	parsedFractions := make(chan types.Fraction, 100)
	wg := &sync.WaitGroup{}
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go parse.StartParseWorker(w, fracStrings, parsedFractions, wg)
	}

	// Close the output channel when all workers are done
	go func() {
		wg.Wait()
		close(parsedFractions)
	}()
	// Combine and sort results
	sortedFractions := types.NewTree()
	var sum types.Fraction
	for fraction := range parsedFractions {
		sum = sum.Add(fraction)
		sortedFractions = sortedFractions.Insert(fraction)
	}

	results := sortedFractions.Walker()
	for fraction := range results {
		log.Printf("%s\n", fraction)
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
