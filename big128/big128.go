package big128

import (
	"math/big"
)

// Some constants to help determine how many big.Words are
// needed to store 128 bits
//
// Borrowed from speter.net/go/exp/math/big/pre128
const (
	bits = 128

	// Compute the size _S of a Word in bytes.
	_m    = ^big.Word(0)
	_logS = _m>>8&1 + _m>>16&1 + _m>>32&1
	_S    = 1 << _logS
	_W    = _S << 3 // word size in bits

	pS = bits / _W // preallocated array size in Words
)

// NewInt() returns a new *big.Int, with preallocated
// storage for 128 bits worth of data and set initially
// to the value of x. It is meant to be a drop-in replacement
// for math/big.NewInt(x), which only allocates memory once.
func NewInt(x int64) *big.Int {
	y := struct {
		bigint   big.Int
		prealloc [pS]big.Word
	}{}

	return y.bigint.SetBits(y.prealloc[0:0]).SetInt64(x)
}

// PreallocInts() preallocates the storage of all the
// big.Ints in the given slice to contain 128 bits. This
// operation only allocates memory once.
//
// Any data in those values will be lost, as this function
// is meant to be used on a newly allocated []big.Int
func PreallocInts(ints []big.Int) {
	mem := make([][pS]big.Word, len(ints))
	for i := range ints {
		ints[i].SetBits(mem[i][0:0])
	}
}

// NewInts() generates a slice of big.Int values where each
// has a 128 bit storage preallocated. This operation only
// allocates memory twice.
func NewInts(n int) []big.Int {
	ints := make([]big.Int, n)
	PreallocInts(ints)
	return ints
}

// NewIntPtrs() generates a slice of *big.Int value where each
// is initialized and has a preallocation of 128 bits for its
// data storage. This operation only allocates memory three times.
func NewIntPtrs(n int) []*big.Int {
	ints := NewInts(n)
	ptrs := make([]*big.Int, n)
	for i := range ints {
		ptrs[i] = &ints[i]
	}
	return ptrs
}
