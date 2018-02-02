package main

import (
  "github.com/VivaLaPanda/fraction-parse/types"
  "github.com/VivaLaPanda/fraction-parse/parse"
  "sync"
)

numWorkers := 5

func main() {
  // Open file


  // Spawn worker to read in file, split on spaces, and push each
  // fraction into a buffered chan
  fracStrings := make(chan string, 100)
  go func(){
    fileToParse, err := ioutil.ReadFile("/tmp/dat")
    tokensToParse = strings.Split(string(fileToParse), " ")
    for _, elem := range tokensToParse {
      fracStrings <- elem
    }

    close(fracStrings)
  }()

  // Spawn workers that take a buffered chan of fraction strings
  // and pushes them as Fraction types unto a results chan
  parsedFractions := make(chan types.Fraction, 100)
  wg := &sync.WaitGroup
  for w := 0; w < numWorkers; w++ {
    wg.Add(1)
    go parse.StartParseWorker(w, fracStrings, parsedFractions, wg)
  }

  // Close the output channel when all workers are done
  go func(){
    wg.Wait()
    close(parsedFractions)
  }()
  // Combine and sort results
  for fraction := range parsedFractions {

  }
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}
