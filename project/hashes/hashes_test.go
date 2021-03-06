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
	"math/rand"
	"sort"
	"testing"
	"time"
)

//all: Compare two slices return true if both have the same lengths and data
//This code was copied from matrix.go
func all(a, b []int64) bool {
	if len(a) != len(b) {
		return false
	}
	for i, _ := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

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

func TestDummy(t *testing.T) {
	fmt.Println("start test")
}

func TestNew(t *testing.T) {
	var aval, bval int64
	aval = 5
	bval = 10
	h := New(aval, bval)
	if h.a != aval {
		t.Errorf("a!=%d", aval)
	}
	if h.b != bval {
		t.Errorf("b!=%d", bval)
	}
}

func TestEqual(t *testing.T) {
	var aval, bval int64
	aval = 5
	bval = 10
	h := New(aval, bval)
	other := New(aval, bval)
	if !h.Equal(other) {
		t.Errorf("False negative\n")
	}
	other = New(aval, bval+1)
	if h.Equal(other) {
		t.Errorf("False positive\n")
	}
}

func TestRand(t *testing.T) {
	h := Rand(rand.New(rand.NewSource(0)))
	var aval, bval int64
	aval = 1432518515
	bval = 3617697886
	fmt.Printf("%v\n", h)
	if h.a != aval {
		t.Errorf("a!=%d", aval)
	}
	if h.b != bval {
		t.Errorf("b!=%d", bval)
	}
}

func initHashes() []Hash {
	var avals, bvals []int64
	var hashes []Hash
	avals = []int64{1, 2, 4, 8, 16, 2394871231}
	bvals = []int64{5, 7, 10, 45, 3, 9283742213}
	numhashes := len(avals)
	hashes = make([]Hash, numhashes)
	for i, a := range avals {
		hashes[i] = New(a, bvals[i])
	}
	return hashes
}
func TestApply(t *testing.T) {
	hashes := initHashes()
	var j int64
	for j = 0; j < 10000; j += 23 {
		fmt.Printf("%d:", j)
		for _, h := range hashes {
			fmt.Printf("\t%d", h.Apply(j))
		}
		fmt.Printf("\n")
	}
}

func TestApplyBatch(t *testing.T) {
	hashes := initHashes()
	var j int64
	var N int64
	N = 10000000 / 23
	t.Logf("N: %d\n", N)
	input := make([]int64, N)
	batchanswer := make([]int64, N)
	for j = 0; j < N; j += 1 {
		input[j] = 23 * j
	}
	rightanswer := make([]int64, N)
	for i, h := range hashes {
		tsserial := time.Now()
		for j, x := range input {
			rightanswer[j] = h.Apply(x)
		}
		teserial := time.Now().Sub(tsserial)
		tsbatch := time.Now()
		h.ApplyBatch(input, batchanswer)
		tebatch := time.Now().Sub(tsbatch)
		fmt.Printf("Apply time: %v\n", teserial)
		fmt.Printf("Apply time batch: %v\n", tebatch)
		if !all(rightanswer, batchanswer) {
			t.Errorf("Failed on hash[%d]\n", i)
			t.Errorf("%v\n", rightanswer)
			t.Errorf("%v\n", batchanswer)
		}
	}
}

func TestApplyBatchSort(t *testing.T) {
	hashes := initHashes()
	var j, N int64
	N = 100000 / 23
	batchanswer := make([]int64, N)
	input := make([]int64, N)
	for j = 0; j < N; j += 1 {
		input[j] = 23 * j
	}
	sortanswer := make([]int64, N)
	for i, h := range hashes {
		tsbatch := time.Now()
		h.ApplyBatch(input, batchanswer)
		tebatch := time.Now().Sub(tsbatch)
		tssort := time.Now()
		h.ApplyBatchSort(input, sortanswer)
		tesort := time.Now().Sub(tssort)
		fmt.Printf("Apply time batch: %v\n", tebatch)
		fmt.Printf("Apply time sort: %v\n", tesort)
		sortedBatch := make([]int, N)
		for i, x := range batchanswer {
			sortedBatch[i] = int(x)
		}
		sort.Sort(sort.IntSlice(sortedBatch))
		for i, x := range sortedBatch {
			batchanswer[i] = int64(x)
		}
		if !all(batchanswer, sortanswer) {
			t.Errorf("Failed on hash[%d]\n", i)
			t.Errorf("%v\n", sortanswer)
			t.Errorf("%v\n", batchanswer)
		}
	}
}
