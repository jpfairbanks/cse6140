package main

import (
	"flag"
	"fmt"
	"github.com/jpfairbanks/cse6140/project/hashes"
	"log"
	"math/rand"
	"time"
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

//Clone: Return a copy of the skecth with the same parameters and hashes,
//but a fresh set of counters. Suggested use is to compare two different
//manipulations for accuracy.
func (cms *CMSketch) Clone() CMSketch {
	Counter := NewMatrix(cms.Depth, cms.Width)
	out := CMSketch{cms.Depth, cms.Width, cms.Hash, Counter}
	return out
}

//Equal: Tell if two cms instances have the same parameters and data.
func (cms *CMSketch) Equal(other *CMSketch) bool {
	depths := cms.Depth == other.Depth
	widths := cms.Width == other.Width
	var hash bool
	hash = true
	for i, h := range cms.Hash {
		hash = other.Hash[i].Equal(h) && hash
		if !hash {
			break
		}
	}
	counters := cms.Counter.Equal(&other.Counter)
	out := depths && widths && hash && counters
	return out
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
//If counts can be negative, then you must estimate differently
func (cms *CMSketch) UpdateSerial(position int64, count int64) {
	for i, h := range cms.Hash {
		cms.Counter.Add(int64(i), h.Apply(position)%cms.Width, count)
	}
}

//AddSignalSlow: Updates the counter and signals completion on a channel
func (cms *CMSketch) AddSignalSlow(row int64, position int64,
	count int64, ch chan int64) {
	h := cms.Hash[row]
	cms.Counter.Add(row, h.Apply(position)%cms.Width, count)
	ch <- 1
}

//UpdateDepthParallelSlow: Insert a single item into the sketch with a count.
func (cms *CMSketch) UpdateDepthParallelSlow(position int64, count int64) {
	ch := make(chan int64, cms.Depth)
	var i int64
	for i = 0; i < cms.Depth; i++ {
		go cms.AddSignalSlow(i, position, count, ch)
	}
	for i = 0; i < cms.Depth; i++ {
		<-ch
	}
}

//AddSignal: Updates the counter and signals completion on a channel
func (cms *CMSketch) AddSignal(start_row int64, end_row int64, position int64,
	count int64, ch chan int64) {
	var h hashes.Hash
	for row := start_row; row < end_row; row++ {
		h = cms.Hash[row]
		cms.Counter.Add(row, h.Apply(position)%cms.Width, count)
	}
	ch <- 1
}

//UpdateDepthParallel: Insert a single item into the sketch with a count.
func (cms *CMSketch) UpdateDepthParallel(position int64, count int64, numProcs int64) {
	var i int64
	NbyP := cms.Depth / numProcs
	ch := make(chan int64, numProcs)
	for i = 0; i < numProcs; i++ {
		go cms.AddSignal(i*NbyP, (i+1)*NbyP, position, count, ch)
	}
	for i = 0; i < numProcs; i++ {
		<-ch
	}
}

//BatchUpdate: insert a batch of edges all at once.
//You must wait for a signal on the channel in order to ensure correct results
func (cms *CMSketch) BatchUpdate(elements []int64, ch chan int64, numProcs int64) {
	batchSize := len(elements)
	for i, h := range cms.Hash {
		yarr := make([]int64, batchSize)
		h.ApplyBatch(elements, yarr)
		for k, y := range yarr {
			yarr[k] = y % cms.Width
		}
		for _, y := range yarr {
			cms.Counter.Add(int64(i), y, 1)
		}
	}
	ch <- 1
}

//BatchUpdateSort: insert a batch of edges all at once.
//You must wait for a signal on the channel in order to ensure correct results
//This method uses a sort to improve cache locality.
func (cms *CMSketch) BatchUpdateSort(elements []int64, ch chan int64, numProcs int64) {
	//TODO: this method needs to be optimized separately from BatchUpdate
	batchSize := len(elements)
	for i, h := range cms.Hash {
		yarr := make([]int64, batchSize)
		h.ApplyBatchSort(elements, yarr)
		for k, y := range yarr {
			yarr[k] = y % cms.Width
		}
		for _, y := range yarr {
			cms.Counter.Add(int64(i), y, 1)
		}
	}
	ch <- 1
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

var depthPtr, widthPtr, efactorPtr *int64
var numProcs int64

func init() {
	depthPtr = flag.Int64("depth", 50, "sets the number of rows (and hash functions) in the CMSketch")
	widthPtr = flag.Int64("width", 80, "sets the number of columns in the CMSketch")
	efactorPtr = flag.Int64("efactor", 10, "number of elements = depth*width*efactor")
	NumProcsPtr := flag.Int64("procs", 1, "number of concurrent worker processes")
	flag.Parse()
	numProcs = *NumProcsPtr
}

func main() {
	//Handling command line parameters
	log.Printf("starting main\n")
	src := rand.NewSource(0)
	r := rand.New(src)
	var Depth, Width, efactor, numElements int64
	Depth = *depthPtr
	Width = *widthPtr
	efactor = *efactorPtr
	numElements = Depth * Width * efactor
	log.Printf("params:Depth:%d\n", Depth)
	log.Printf("params:Width:%d\n", Width)
	log.Printf("params:efactor:%d\n", efactor)
	log.Printf("params:numElements:%d\n", numElements)

	//Initialize Data Structures
	hslice := RandomHashes(r, Depth)
	cms := NewCMSketch(Depth, Width)
	cms.Hash = hslice
	//Make the zipf distribution of random input
	var j, z int64
	var s, v float64
	var imax uint64
	s = 1.2
	v = 1.0
	imax = 2 << 10
	zipfer := rand.NewZipf(r, s, v, imax)

	//Use set to store the exact answers
	set := make(map[int64]int64)
	log.Printf("Inserting\n")
	ts := time.Now()
	for j = 0; j < numElements; j++ {
		z = int64(zipfer.Uint64())
		//set[z] += 1
		//fmt.Println(z)
		cms.UpdateSerial(z, 1)
	}
	te := time.Now().Sub(ts)
	fmt.Printf("time: %v\n", te)
	fmt.Printf("%s\n", cms.Counter.String())
	var qj int64 //approximate answers
	var totalLoss float64
	loss := func(cj, qj float64) float64 {
		return (cj - qj) * (cj - qj)
	}
	for j, cj := range set {
		qj = cms.PointQuery(j)
		//fmt.Printf("results:%d %d %d %f\n", j, qj, cj, float64(qj)/float64(cj))
		totalLoss += loss(float64(cj), float64(qj))
	}
	fmt.Printf("Total Loss: %f/%d\n", totalLoss, numElements)
}
