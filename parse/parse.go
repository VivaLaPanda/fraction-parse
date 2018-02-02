package parse

import (
	"github.com/VivaLaPanda/fraction-parse/types"
	"sync"
)

func StartParseWorker(id int, input <-chan string, output chan<- types.Fraction, wg sync.WaitGroup) {
	defer wg.Done()
	for token := range input {
		_ = token
		output <- types.Fraction{}
	}

	return
}
