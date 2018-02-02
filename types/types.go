package types

import "fmt"

// Fraction is a type which represents any
type Fraction struct {
	Numerator, Denominator int
}

// Add is a method which returns the sum of the given fraction and the Fraction
// passed as a parameter
func (f Fraction) Add(other Fraction) Fraction {
	f.Numerator = (f.Numerator * other.Denominator) + (other.Numerator * f.Denominator)
	f.Denominator = f.Denominator * other.Denominator

	f.Reduce()

	return f
}

// LessThan is a method which checks to see whether the given fraction is
// less than the fraction passed to the function
func (f Fraction) LessThan(other Fraction) bool {
	f = f.Reduce()
	other = other.Reduce()

	if f.Numerator*other.Denominator < f.Denominator*other.Numerator {
		// Cross multiply comparison
		return true
	}

	return false
}

// Reduce makes given fraction reduced.
func (f Fraction) Reduce() Fraction {
	gcd := gcdEuclidean(f.Numerator, f.Denominator)
	f.Numerator /= gcd
	f.Denominator /= gcd

	return f
}

func (f Fraction) String() string {
	wholeNum, regFraction := f.makeRegular()

	return fmt.Sprintf("%d_%d/%d", wholeNum, regFraction.Numerator, regFraction.Denominator)
}

// MakeRegular is a method that Given an irregular fraction will move the
// extra fraction component into the whole number component
func (f Fraction) makeRegular() (int, Fraction) {
	var wholeNumber int
	if f.Numerator > f.Denominator {
		f.Numerator = f.Numerator - f.Denominator
		wholeNumber, f = f.makeRegular()
		wholeNumber++
	}

	return wholeNumber, f
}

// GCDEuclidean calculates GCD by Euclidian algorithm.
func gcdEuclidean(a, b int) int {
	for a != b {
		if a > b {
			a -= b
		} else {
			b -= a
		}
	}

	return a
}
