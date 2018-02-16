package types

// A Tree is a binary tree with Fraction values.
type Tree struct {
	Left  *Tree
	Value Fraction
	Right *Tree
}

// Walk traverses a tree depth-first,
// sending each Value on a channel.
func walk(t *Tree, ch chan Fraction) {
	if t == nil {
		return
	}
	walk(t.Left, ch)
	ch <- t.Value
	walk(t.Right, ch)
}

// Walker launches Walk in a new goroutine,
// and returns a read-only channel of values.
func (t *Tree) Walker() <-chan Fraction {
	ch := make(chan Fraction)
	go func() {
		walk(t, ch)
		close(ch)
	}()
	return ch
}

// Compare reads values from two Walkers
// that run simultaneously, and returns true
// if t1 and t2 have the same contents.
func compare(t1, t2 *Tree) bool {
	c1, c2 := t1.Walker(), t2.Walker()
	for {
		v1, ok1 := <-c1
		v2, ok2 := <-c2
		if !ok1 || !ok2 {
			return ok1 == ok2
		}
		if v1 != v2 {
			break
		}
	}
	return false
}

// New returns a new, random binary tree
// holding the values 1k, 2k, ..., nk.
func NewTree() *Tree {
	var t *Tree
	return t
}

// Insert returns the result of inseting the given fraction into the given
// tree
func (t *Tree) Insert(v Fraction) *Tree {
	if t == nil {
		return &Tree{nil, v, nil}
	}
	if v.LessThan(t.Value) {
		t.Left = t.Left.Insert(v)
		return t
	}
	t.Right = t.Right.Insert(v)
	return t
}
