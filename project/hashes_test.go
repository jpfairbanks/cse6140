/* hashes_test.go: testing the hashes.go file
Author: James Fairbanks
Liscence: BSD

we need to find primes and compute a*x+b mod p as the hash,
then we need to test that these hashes are indeed pairwise independent for
distinct values of a and b.
*/

package hashes

import (
	"fmt"
	"testing"
)

//const MOD int64 = (2 << 31) - 1
//const HL uint64 = 31
//
////Hash: a struct for hashing things
//type Hash struct {
//	a, b int64
//}
//
////Apply: compute the hash for a single input you need to take this modulo the width
//// in order to make sure it fits in the array that you want to store it in.
//func (h Hash) Apply(x int64) int64 {
//	var result int64
//	result = (h.a*x + h.b)
//	result = (result >> HL) % MOD
//	return result
//}
//
////New: create a new hash and make sure that it is valid
//func New(a int64, b int64) Hash {
//	return Hash{a, b}
//}

func TestNew(t *testing.T) {
	var aval, bval int64
	aval = 5
	bval = 10
	h := New(aval, 10)
	if h.a != aval {
		t.Errorf("a!=%d", aval)
	}
	if h.b != bval {
		t.Errorf("b!=%d", bval)
	}
}

func TestApply(t *testing.T) {
	var avals, bvals []int64
	var hashes []Hash
	avals = []int64{1, 2, 4, 8, 16}
	bvals = []int64{5, 7, 10, 45, 3}
	for i, a := range avals {
		hashes[i] = New(a, bvals[i])
	}
	var j int64
	for j = 0; j < 1000; j++ {
		for i, h := range hashes {
			fmt.Printf("%d %d:%d\n", j, i, h.Apply(j))
		}
	}
}
