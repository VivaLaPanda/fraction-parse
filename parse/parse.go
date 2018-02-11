package parse

import (
	"fmt"
	"math"
	"strconv"
	"sync"
	"unicode"

	"github.com/VivaLaPanda/fraction-parse/types"
)

var stateTable = [][]int{
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 1, 2, 2, 0, 0, 4, 3, 0},
	{0, 0, 0, 0, 0, 0, 4, 3, 0},
	{0, emitState, 0, 0, 8, 6, 4, 3, 0},
	{0, 0, 0, 0, 0, 0, 0, 5, 0},
	{emitState, emitState, emitState, emitState, emitState, emitState, emitState, 5, emitState},
	{0, 0, 0, 0, 0, 0, 0, 7, 0},
	{0, 0, 0, 0, 8, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 9, 0},
	{emitState, emitState, emitState, emitState, emitState, emitState, emitState, emitState, emitState},
}

var actionTable = [][]int{
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, strSign, strSign, 0, 0, strWhole, 0, 0},
	{0, 0, 0, 0, 0, 0, strWhole, 0, 0},
	{0, strWhole, 0, 0, strNum, strWhole, strWhole, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{strDec, strDec, strDec, strDec, strDec, strDec, strDec, 0, strDec},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, strNum, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{strDenom, strDenom, strDenom, strDenom, strDenom, strDenom, strDenom, strDenom, strDenom},
}

const (
	errState  = 0
	emitState = -1
	strWhole  = -2
	strNum    = -3
	strDec    = -4
	strDenom  = -5
	strSign   = -6
)

func StartParseWorker(id int, input <-chan string, output chan<- types.Fraction, wg *sync.WaitGroup) {
	defer wg.Done()

	// Here we are iterating over our channel of fraction strings
	for token := range input {
		state := 1
		action := 0
		buffer := ""
		negativeFlag := false
		fraction := types.Fraction{Numerator: 0, Denominator: 1}
		wholeNumber := 0

		// Mini helper to do last minute changes before emitting
		var emit = func(out types.Fraction) {
			// Deal with the whole number component
			out = out.Add(types.Fraction{Numerator: wholeNumber, Denominator: 1})

			if negativeFlag {
				out.Numerator *= -1
			}
			output <- out
		}

		// This loops over the runes in a single fraction string
		for _, r := range token + " " {
			action = queryTable(actionTable, state, r)
			state = queryTable(stateTable, state, r)

			currentBufferAsInt, _ := strconv.Atoi(buffer)

			switch action {
			case strWhole: // We have a whole number. Store
				wholeNumber = currentBufferAsInt
				buffer = ""
			case strNum: // Store just the numerator
				fraction.Numerator = currentBufferAsInt
				buffer = ""
			case strDenom: // Store the denominator
				fraction.Denominator = currentBufferAsInt
				buffer = ""
			case strDec: // Store the decimal part, convert to fraction
				tempFrac := types.Fraction{Numerator: 0, Denominator: 1}
				tempFrac.Numerator = currentBufferAsInt
				tempFrac.Denominator = int(math.Pow(float64(10), float64(len(buffer))))
				fraction = fraction.Add(tempFrac)
				buffer = ""
			case strSign:
				negativeFlag = true
				buffer = ""
			}

			switch state {
			case errState:
				fmt.Printf("FATAL ERROR: %v is not a valid fraction.\n", token)
				break
			case emitState:
				emit(fraction)
				break
			default:
				// If we did an action we want to ignore the last rune
				if action == 0 {
					buffer += string(r) // Append current rune to buffer
				}
			}
		}
	}

	return
}

func queryTable(table [][]int, state int, inputRune rune) int {
	return table[state][symbolToIndex(inputRune)]
}

func symbolToIndex(r rune) int {
	if unicode.IsSpace(r) {
		return 1
	} else if unicode.IsDigit(r) {
		return 7
	}

	switch r {
	case '-':
		return 2
	case '+':
		return 3
	case '/':
		return 4
	case '_':
		return 5
	case '.':
		return 6
	default:
		return 0
	}
}
