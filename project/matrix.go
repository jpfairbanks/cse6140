/* matrix.go: atomicly updated 2d array
Author: James Fairbanks
Date: 2013-09-20
Licence: BSD
*/

package main

import (
	"fmt"
	"strings"
	"sync/atomic"
)

//Matrix: a logically 2d array of data stored in row major order
//We are using row major because when we parallelize over the hash functions in the count min sketch,
//each hash function will make random access to a single row. When we query a point, then we will access a random element from each row.
type Matrix struct {
	Nrows int64
	Ncols int64
	Data  []int64
}

//all: Compare two slices return true if both have the same lengths and data
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

func NewMatrix(Nrows, Ncols int64) Matrix {
	data := make([]int64, Nrows*Ncols)
	return Matrix{Nrows, Ncols, data}
}

//Read: Access an element from a Row Major order Matrix
//TODO: test atomicity of reads, we might need to use atomic.load(ptr)
func (M *Matrix) Read(row, col int64) int64 {
	return M.Data[(row*M.Ncols)+col]
}

//AtomicAdd: Increment an element Matrix
//Returns old+inc and guarantees that no data was lost
func (M *Matrix) AtomicAdd(row, col int64, inc int64) int64 {
	addr := &M.Data[(row*M.Ncols)+col]
	return atomic.AddInt64(addr, inc)
}

//Add: add to an element from a Matrix without synchronizing
//Returns old+inc and not concurrent safe
func (M *Matrix) Add(row, col int64, inc int64) int64 {
	index := (row * M.Ncols) + col
	M.Data[index] += inc
	return M.Data[index]
}

//String: Print a Matrix like in numpy
func (M *Matrix) String() string {
	strarr := make([]string, M.Nrows)
	var str string
	var i int64
	for i = 0; i < M.Nrows; i++ {
		str = fmt.Sprintf("%v", M.Data[i*M.Ncols:(i+1)*M.Ncols])
		strarr[i] = str
	}
	return strings.Join(strarr, "\n")
}

//Equal: Test for deep equality of matrices
func (M *Matrix) Equal(other *Matrix) bool {
	if M.Nrows != other.Nrows {
		return false
	}
	if M.Ncols != other.Ncols {
		return false
	}
	if !all(M.Data, other.Data) {
		return false
	}
	return true
}
