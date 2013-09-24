package main

import (
	"fmt"
	"github.com/jpfairbanks/cse6140/project/hashes"
	"math/rand"
)

//CMSketch: A streaming data structure for tracking the elements of
//a long vector in a space efficient manner.
//Provides update and query operations
type CMSketch struct {
	Depth   int64
	Width   int64
	Hash    []hashes.Hash
	Counter Matrix
}

//NewCMSketch: Allocates a new sketch
//Does not create the hashes. We leave that up to you so that
//the consumer can pick hashes according to different schema.
func NewCMSketch(Depth, Width int64) CMSketch {
	Hash := make([]hashes.Hash, Depth)
	Counter := NewMatrix(Depth, Width)
	cms := CMSketch{Depth, Width, Hash, Counter}
	return cms
}

//RandomHashes: Make uniformly random hashes from a user provided
//stream of random numbers.
func RandomHashes(r *rand.Rand, Depth int64) []hashes.Hash {
	hslice := make([]hashes.Hash, Depth)
	for i, _ := range hslice {
		hslice[i] = hashes.Rand(r)
	}
	return hslice
}

//UpdateSerial: Insert a single item into the sketch with a count.
//If counts can be negative, then you mist estimate differently
func (cms *CMSketch) UpdateSerial(position int64, count int64) {
	for i, h := range cms.Hash {
		cms.Counter.Add(int64(i), h.Apply(position)%cms.Width, count)
	}
}

//AddSignal: Updates the counter and signals completion on a channel
func (cms *CMSketch) AddSignal(row int64, position int64,
	count int64, ch chan int64) {
	h := cms.Hash[row]
	cms.Counter.Add(row, h.Apply(position)%cms.Width, count)
	ch <- 1
}

//UpdateDepthParallel: Insert a single item into the sketch with a count.
//If counts can be negative, then you mist estimate differently
func (cms *CMSketch) UpdateDepthParallel(position int64, count int64) {
	ch := make(chan int64, cms.Depth)
	var i int64
	for i = 0; i < cms.Depth; i++ {
		go cms.AddSignal(i, position, count, ch)
	}
	for i = 0; i < cms.Depth; i++ {
		<-ch
	}
}

//PointQuery: query the value at position
func (cms *CMSketch) PointQuery(position int64) int64 {
	temp := make([]int64, cms.Depth)
	var h hashes.Hash
	var i int64
	for i = 0; i < cms.Depth; i++ {
		h = cms.Hash[i]
		temp[i] = cms.Counter.Read(i, h.Apply(position)%cms.Width)
	}
	mini := temp[0]
	for i = 0; i < cms.Depth; i++ {
		if mini > temp[i] {
			mini = temp[i]
		}
	}
	return mini
}

func main() {
	fmt.Println("starting main")
	src := rand.NewSource(0)
	r := rand.New(src)
	var Depth, Width int64
	Depth = 160
	Width = 480
	hslice := RandomHashes(r, Depth)
	cms := NewCMSketch(Depth, Width)
	cms.Hash = hslice
	fmt.Printf("%s\n", cms.Counter.String())
	fmt.Printf("Inserting\n")
	cms.UpdateSerial(1, 1)
	fmt.Printf("%s\n", cms.Counter.String())
	fmt.Printf("Inserting\n")
	var j, z int64
	var s, v float64
	var imax uint64
	s = 1.2
	v = 1.0
	imax = 2 << 10
	zipfer := rand.NewZipf(r, s, v, imax)
	set := make(map[int64]int64)
	for j = 0; j < Depth*Width*10; j++ {
		z = int64(zipfer.Uint64())
		set[z] += 1
		fmt.Println(z)
		cms.UpdateSerial(z, 1)
	}
	fmt.Printf("%s\n", cms.Counter.String())
	var qj int64
	for j, cj := range set {
		qj = cms.PointQuery(j)
		fmt.Printf("results:%d %d %d %f\n", j, qj, cj, float64(qj)/float64(cj))
	}
}
