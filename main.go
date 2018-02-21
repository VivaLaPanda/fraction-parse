package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/VivaLaPanda/fraction-parse/parse"
	"github.com/VivaLaPanda/fraction-parse/types"
)

const numWorkers = 20

func main() {
	filepath := os.Args[1]

	// Open file
	var data []byte
	var err error
	if len(os.Args[1:]) == 0 {
		data, err = ioutil.ReadAll(os.Stdin)
		check(err)
	} else {
		data, err = ioutil.ReadFile(filepath)
		check(err)
	}

	cleanString := strings.TrimSuffix(string(data), "\n")

	// Spawn worker to read in file, split on spaces, and push each
	// fraction into a buffered chan
	fracStrings := make(chan string, 100)
	go func() {
		tokensToParse := strings.Fields(cleanString)
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
