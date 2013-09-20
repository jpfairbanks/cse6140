/* matrix.go: atomicly updated 2d array
Author: James Fairbanks
Date: 2013-09-20
Licence: BSD
*/

package main

import (
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

func NewMatrix(Nrows, Ncols int64) Matrix {
	data := make([]int64, Nrows*Ncols)
	return Matrix{Nrows, Ncols, data}
}

//Read: Access an element from a Row Major order Matrix
//TODO: test atomicity of reads, we might need to use atomic.load(ptr)
func (M *Matrix) Read(row, col int64) int64 {
	return M.Data[(row*M.Ncols)+col]
}

//AtomicAdd: Increment an element from a Row Major order Matrix
//Returns old+inc and guarantees that no data was lost
func (M *Matrix) AtomicAdd(row, col int64, inc int64) int64 {
	addr := &M.Data[(row*M.Ncols)+col]
	return atomic.AddInt64(addr, inc)
}

//Add: add to an element from a Row Major order Matrix without synchronizing
//Returns old+inc and guarantees that no data was lost
func (M *Matrix) Add(row, col int64, inc int64) {
	M.Data[(row*M.Ncols)+col] += inc
}

func updater(M *Matrix, ch chan int64) {
	var i, j int64
	for i = 0; i < M.Ncols; i++ {
		for j = 0; j < 100; j++ {
			M.AtomicAdd(1, i, 1)
		}
	}
	ch <- 1
}

func raceyUpdater(M *Matrix, ch chan int64) {
	var i, j int64
	for i = 0; i < M.Ncols; i++ {
		for j = 0; j < 100; j++ {
			M.Add(0, i, 1)
		}
	}
	ch <- 1
}
