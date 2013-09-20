/* hashes.go: utilities to help with hashing things.
Author: James Fairbanks
Liscence: BSD

we need to find primes and compute a*x+b mod p as the hash,
then we need to test that these hashes are indeed pairwise independent for
distinct values of a and b.
*/

package hashes

import(
	"fmt"
	"crypto/rand"
)

//Hash: a struct for hashing things
type struct Hash {
	var a, b, p int
}

/*
func (h Hash) Write() int, error {
	buff := make([]int, buffsize)
	fmt.Println("write is fake implemented")
	return 1, error
}
*/

//Return: the hash for a single input
func (h Hash) hash(x int) int {
	return (a * x + b) % p
}

//New: create a new hash and make sure that it is valid
func New(a,b,p int) Hash {
	
}

//New: create a new hash and make sure that it is valid
func New(a,b int) func(int) int {
	p := rand.Prime(crypto.Reader, 64)
	return func(x int)int {return (a * x + b) % p}
}
