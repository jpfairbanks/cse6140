/* hashes.go: utilities to help with hashing things.
Author: James Fairbanks
Liscence: BSD

we need to find primes and compute a*x+b mod p as the hash,
then we need to test that these hashes are indeed pairwise independent for
distinct values of a and b.
*/

package hashes

import (
	"github.com/cznic/sortutil"
	"math/rand"
)

//MOD: a nice big prime
const MOD int64 = (2 << 31) - 1

//HL: highlevel bits are defined as top 31 bits.
const HL uint64 = 31

//Hash: a struct for hashing things
type Hash struct {
	a, b int64
}

//Apply: compute the hash for a single input you need to take this modulo the width
// in order to make sure it fits in the array that you want to store it in.
func (h *Hash) Apply(x int64) int64 {
	var result int64
	result = (h.a*x + h.b)
	//take the highbits 64bit+ lowbits and make sure what is left is less than 2^31-1
	result = ((result >> HL) + result) & MOD
	return result
}

//ApplyBatch: compute the hash for an array of input, you need to take this modulo the width
// in order to make sure it fits in the array that you want to store it in.
// The array yarr will be modified to house the output.
// If you want to use this function in parallel, then you must provide different arrays to each thread.
// The allocation of yarr must happen before you call this otherwise. We do not check that yarr is the right size.
// Very unsafe programming practice imported from C. However this does shave off function call overhead compared to
// calling Apply in a tight loop.
func (h *Hash) ApplyBatch(xarr []int64, yarr []int64) {
	var result int64
	for i, x := range xarr {
		result = (h.a*x + h.b)
		//take the highbits 64bit+ lowbits and make sure what is left is less than 2^31-1
		result = ((result >> HL) + result) & MOD
		yarr[i] = result
	}
}

//ApplyBatchSort: Applies the hash to a batch and sorts the output in increasing order.
//Takes and input and an output array, and leaves the sorted output values in the output array.
func (h *Hash) ApplyBatchSort(xarr []int64, yarr []int64) {
	var result int64
	for i, x := range xarr {
		result = (h.a*x + h.b)
		//take the highbits 64bit+ lowbits and make sure what is left is less than 2^31-1
		result = ((result >> HL) + result) & MOD
		yarr[i] = result
	}
	sortutil.Int64Slice(yarr).Sort()
}

//New: create a new hash and make sure that it is valid
func New(a int64, b int64) Hash {
	return Hash{a, b}
}

func Rand(r *rand.Rand) Hash {
	a := rand.Int63n(MOD)
	b := rand.Int63n(MOD)
	return New(a, b)
}

//Equal: check that all of the hashes have the same parameters
func (h Hash) Equal(other Hash) bool {
	if h.a != other.a {
		return false
	}
	if h.b != other.b {
		return false
	} else {
		return true
	}
}
