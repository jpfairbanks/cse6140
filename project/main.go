package main

import (
	"fmt"
	"github.com/jpfairbanks/cse6140/project/hashes"
	"math/rand"
)

//CMSketch: A streaming data structure for tracking the elements of a long vector in a space efficient manner.
//Provides update and query operations
type CMSketch struct {
	Width   int64
	Depth   int64
	Hash    []hashes.Hash
	Counter Matrix
}

func NewCMSketch(Width, Depth int64) CMSketch {
	Hash := make([]hashes.Hash, Depth)
	Counter := NewMatrix(Width, Depth)
	cms := CMSketch{Width, Depth, Hash, Counter}
	return cms
}

func RandomHashes(r *rand.Rand, Depth int64) []hashes.Hash {

	hslice := make([]hashes.Hash, Depth)
	for i, _ := range hslice {
		hslice[i] = hashes.Rand(r)
	}
	return hslice
}
func (cms *CMSketch) UpdateSerial(position int64, count int64) {
	for i, h := range cms.Hash {
		cms.Counter.Add(int64(i), h.Apply(position)%cms.Width, count)
	}
}

func main() {
	fmt.Println("starting main")
	src := rand.NewSource(0)
	r := rand.New(src)
	var Depth int64
	Depth = 10
	hslice := RandomHashes(r, Depth)
	fmt.Println(hslice)
}
