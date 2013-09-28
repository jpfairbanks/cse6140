package main

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"testing"
	"time"
)

func tic() time.Time {
	return time.Now()
}

func toc(ts time.Time) time.Duration {
	return time.Now().Sub(ts)
}

func TestInit(t *testing.T) {
	var d, w int64
	d = 75
	w = 700
	cms := NewCMSketch(d, w)
	if cms.Width != w {
		t.Errorf("Width: %d != %d\n", w, cms.Width)
	}
	if cms.Depth != d {
		t.Errorf("Depth: %d != %d\n", d, cms.Depth)
	}
}

func TestInsert(t *testing.T) {
	cms := NewCMSketch(2, 5)
	r := rand.New(rand.NewSource(0))
	rh := RandomHashes(r, 2)
	cms.Hash = rh
	if rh[0].Apply(1)%5 != 2 {
		t.Errorf("hash[0](1) % 5 != 2\n")
	}
	if rh[1].Apply(1)%5 != 3 {
		t.Errorf("hash[0](1) % 5 != 2\n")
	}
	fmt.Printf("%v\n", cms)
	t.Logf("Before:\n%v\n", cms.Counter.String())
	cms.UpdateSerial(1, 1)
	t.Logf("After:\n%v\n", cms.Counter.String())
	if cms.Counter.Read(0, 2) != 1 {
		t.Errorf("Bad first row\n")
	}
	if cms.Counter.Read(1, 3) != 1 {
		t.Errorf("Bad second row\n")
	}
}

//loss: implements a generic loss function for scoring the implementation.
func loss(cj, qj float64) float64 {
	return (cj - qj) * (cj - qj)
}

//makeCMS: Initialize the CMSKETCH to standard values for test
func makeCMS(r *rand.Rand) *CMSketch {
	var Depth, Width int64
	Depth = 80
	Width = 160
	hslice := RandomHashes(r, Depth)
	cms := NewCMSketch(Depth, Width)
	cms.Hash = hslice
	return &cms
}

//makeZipfer: Initialize the stream of random elements for the tests.
func makeZipfer(r *rand.Rand) *rand.Zipf {
	//Make the zipf distribution of random input
	var s, v float64
	var imax uint64
	s = 1.2
	v = 1.0
	imax = 2 << 10
	zipfer := rand.NewZipf(r, s, v, imax)
	return zipfer
}

//TestSpeed: Test the speed for Serial insertions
func TestSpeed(t *testing.T) {
	src := rand.NewSource(0)
	r := rand.New(src)
	cms := *makeCMS(r)
	var efactor, numElements int64
	efactor = 10
	numElements = cms.Depth * cms.Width * efactor
	zipfer := makeZipfer(r)
	log.Printf("Inserting\n")
	ts := time.Now()
	var j, z int64
	for j = 0; j < numElements; j++ {
		z = int64(zipfer.Uint64())
		cms.UpdateSerial(z, 1)
	}
	te := time.Now().Sub(ts)
	fmt.Printf("time Serial: %v\n", te)
	numProcs := runtime.NumCPU()
	runtime.GOMAXPROCS(numProcs)
	ts = time.Now()
	for j = 0; j < numElements; j++ {
		z = int64(zipfer.Uint64())
		cms.UpdateDepthParallel(z, 1, int64(numProcs))
	}
	te = time.Now().Sub(ts)
	fmt.Printf("time parallel: %v\n", te)
}

//TestAccuracy: Run a test storing the right answers
//compute the loss using the function loss
func TestAccuracy(t *testing.T) {
	src := rand.NewSource(0)
	r := rand.New(src)
	cms := *makeCMS(r)
	var efactor, numElements int64
	efactor = 10
	numElements = cms.Depth * cms.Width * efactor
	zipfer := makeZipfer(r)

	//Use set to store the exact answers
	set := make(map[int64]int64)
	log.Printf("Inserting\n")
	ts := time.Now()
	var j, z int64
	for j = 0; j < numElements; j++ {
		z = int64(zipfer.Uint64())
		set[z] += 1
		//fmt.Println(z)
		cms.UpdateSerial(z, 1)
	}
	te := time.Now().Sub(ts)
	fmt.Printf("time accuracy: %v\n", te)
	//fmt.Printf("%s\n", cms.Counter.String())
	var qj int64 //approximate answers
	var totalLoss float64
	for j, cj := range set {
		qj = cms.PointQuery(j)
		//fmt.Printf("results:%d %d %d %f\n", j, qj, cj, float64(qj)/float64(cj))
		totalLoss += loss(float64(cj), float64(qj))
	}
	fmt.Printf("Total Loss: %f/%d\n", totalLoss, numElements)
}

//TestBatchInsert: Compare the data structures after making batch inserts to one
//where we have made serial insertions to test that they produce the same result
func TestBatchInsert(t *testing.T) {
	src := rand.NewSource(0)
	r := rand.New(src)
	cms := makeCMS(r)
	var numElements, batchsize int64
	batchsize = 10
	zipfer := makeZipfer(r)

	elements := make([]int64, batchsize)
	var i, z int64
	var ts time.Time
	var te time.Duration
	ts = tic()
	for i = 0; i < batchsize; i++ {
		z = int64(zipfer.Uint64())
		elements[i] = z
	}
	te = toc(ts)
	fmt.Printf("time zipfer: %s\n", te)
	ts = tic()
	for i = 0; i < numElements; i++ {
		cms.UpdateSerial(z, 1)
	}
	te = toc(ts)
	fmt.Printf("time single insertions: %v\n", te)
	t.Logf("cms:\n%v\n", cms)
	batchcms := cms.Clone()
	ch := make(chan int64)
	go batchcms.BatchUpdate(elements, ch)
	<-ch
	result := cms.Equal(&batchcms)
	if !result {
		t.Errorf("the sketches did not come up equal\n")
		t.Logf("batchcms:\n%v\n", batchcms)
	}
}

func TestDecrement(t *testing.T) {
	t.Fail()
}
