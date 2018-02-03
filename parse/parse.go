package parse

import (
	"fmt"
	"github.com/VivaLaPanda/fraction-parse/types"
	"math"
	"strconv"
	"sync"
	"unicode"
)

var dfaTable = [][]int{
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 1, 2, 2, 0, 0, 0, 3, 0},
	{0, 0, 0, 0, 0, 0, 0, 3, 0},
	{0, strWholeRet, 0, 0, strNumGo8, strWholeGo6, strWholeGo4, 3, 0},
	{0, 0, 0, 0, 0, 0, 0, 5, 0},
	{strDecRet, strDecRet, strDecRet, strDecRet, strDecRet, strDecRet, strDecRet, 5, strDecRet},
	{0, 0, 0, 0, 0, 0, 0, 7, 0},
	{0, 0, 0, 0, strNumGo8, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 9, 0},
	{strDenomRet, strDenomRet, strDenomRet, strDenomRet, strDenomRet, strDenomRet, strDenomRet, strDenomRet, strDenomRet},
}

const (
	errState    = 0
	strWholeRet = -1
	strNumGo8   = -2
	strWholeGo6 = -3
	strWholeGo4 = -4
	strDecRet   = -5
	strDenomRet = -6
)

func StartParseWorker(id int, input <-chan string, output chan<- types.Fraction, wg *sync.WaitGroup) {
	defer wg.Done()
	for token := range input {
		state := 1
		buffer := ""
		fraction := types.Fraction{0, 1}
		for _, r := range token + " " {
			state = queryTable(state, r)
			currentBufferAsInt := 0
			currentBufferAsInt, _ = strconv.Atoi(buffer)
			// Crazy switch incoming!!!
			// This switch deals with what we determined to be "special states"
			// If no special state is hit, then we continue adding to the buffer
			// like normal
			switch state {
			case errState:
				fmt.Printf("FATAL ERROR: %v is not a valid fraction.\n", token)
				break
			case strWholeRet: // We have a whole number. Store and emit
				fraction.Numerator += fraction.Denominator * currentBufferAsInt
				buffer = ""
				output <- fraction
				break
			case strNumGo8: // Store just the numerator and then go to state 8
				fraction.Numerator = currentBufferAsInt
				buffer = ""
				state = 8
			case strWholeGo4: // Store the whole number part and then go to state 4
				fraction.Numerator += fraction.Denominator * currentBufferAsInt
				buffer = ""
				state = 4
			case strWholeGo6: // Store the whole number part and then go to state 6
				fraction.Numerator += fraction.Denominator * currentBufferAsInt
				buffer = ""
				state = 6
			case strDecRet: // Store the decimal part, convert to fraction, and then emit
				tempFrac := types.Fraction{1, 0}
				tempFrac.Numerator = currentBufferAsInt
				tempFrac.Denominator = int(math.Pow(float64(10), float64(len(buffer))))
				fraction = fraction.Add(tempFrac)
				buffer = ""
				output <- fraction
				break
			case strDenomRet: // Store the denominator and then emit
				fraction.Denominator = currentBufferAsInt
				buffer = ""
				output <- fraction
				break
			default: // We are parsing a number right now it looks like
				buffer += string(r) // Append current rune to buffer
			}
		}
	}

	return
}

func queryTable(state int, inputRune rune) int {
	return dfaTable[state][symbolToIndex(inputRune)]
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
