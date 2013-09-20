/* hashes.go: utilities to help with hashing things.
Author: James Fairbanks
Liscence: BSD

we need to find primes and compute a*x+b mod p as the hash,
then we need to test that these hashes are indeed pairwise independent for
distinct values of a and b.
*/

package hashes

import ()

const MOD int64 = (2 << 31) - 1
const HL uint64 = 31

//Hash: a struct for hashing things
type Hash struct {
	a, b int64
}

//Apply: compute the hash for a single input you need to take this modulo the width
// in order to make sure it fits in the array that you want to store it in.
func (h Hash) Apply(x int64) int64 {
	var result int64
	result = (h.a*x + h.b)
	result = (result >> HL) % MOD
	return result
}

//New: create a new hash and make sure that it is valid
func New(a int64, b int64) Hash {
	return Hash{a, b}
}
