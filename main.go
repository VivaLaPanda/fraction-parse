package main

import (
	"flag"
	"fmt"
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
	wg.Add(numWorkers)
	for w := 0; w < numWorkers; w++ {
		go parse.StartParseWorker(w, fracStrings, parsedFractions, wg)
	}

	// Close the output channel when all workers are done
	go func() {
		wg.Wait()
		close(parsedFractions)
	}()

	//results := sortedFractions.Walker()
	sum := types.Fraction{Numerator: 0, Denominator: 1}
	sortedFractions := types.NewTree()
	for fraction := range parsedFractions {
		sum = sum.Add(fraction)
		sortedFractions = sortedFractions.Insert(fraction)
	}

	results := sortedFractions.Walker()

	fmt.Printf("The sum of the fractions is: %s\n", sum)
	for fraction := range results {
		fmt.Printf("%s\n", fraction)
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
