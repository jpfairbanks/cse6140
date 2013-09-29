package main

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"testing"
	"time"
)

var width int64 = 100000
var depth int64 = 20

func tic() time.Time {
	return time.Now()
}

func toc(ts time.Time) time.Duration {
	return time.Now().Sub(ts)
}
func makeRectangle(width, depth int64) []int64 {
	size := width * depth
	fmt.Printf("size: %d\n", size)
	arr := make([]int64, size)
	for i, _ := range arr {
		arr[i] = int64(i)
	}
	return arr
}

func TestSerialReads(t *testing.T) {
	log.Printf("starting serial benchmark")
	var ts time.Time
	var te time.Duration
	var sum int64
	arr := makeRectangle(width, depth)
	ts = tic()
	for _, x := range arr {
		sum += x
	}
	te = toc(ts)
	fmt.Printf("sum: %d\n", sum)
	fmt.Printf("time serial: %s\n", te)
}

func rowSum(arr []int64, p int64, ch chan int64) {
	var sum int64
	for i := p * width; i < (p+1)*width; i++ {
		sum += arr[i]
	}
	ch <- sum
}

func rowSumPermute(arr []int64, p int64, ch chan int64, perm []int) {
	var sum int64
	subslice := arr[p*width : (p+1)*width]
	for _, pi := range perm {
		sum += subslice[pi]
	}
	ch <- sum
}

func TestParallelReads(t *testing.T) {
	log.Printf("starting parallel benchmark")
	var ts time.Time
	var te time.Duration
	var sum int64
	var p int64
	arr := makeRectangle(width, depth)
	NumP := int64(runtime.NumCPU())
	ch := make(chan int64, NumP)
	ts = tic()
	for p = 0; p < depth; p++ {
		go rowSum(arr, p, ch)
	}
	var tmp int64
	for p = 0; p < depth; p++ {
		tmp = <-ch
		sum += tmp
	}
	te = toc(ts)
	fmt.Printf("sum: %d\n", sum)
	fmt.Printf("time parallel: %s\n", te)
}

func TestParallelReadsPerm(t *testing.T) {
	log.Printf("starting parallel random benchmark")
	var ts time.Time
	var te time.Duration
	var sum int64
	var p int64
	arr := makeRectangle(width, depth)
	NumP := int64(runtime.NumCPU())
	ch := make(chan int64, NumP)
	var rander []*rand.Rand
	rander = make([]*rand.Rand, depth)
	for p = 0; p < depth; p++ {
		rander[p] = rand.New(rand.NewSource(p))
	}

	perm := make([][]int, depth)
	for p = 0; p < depth; p++ {
		perm[p] = rander[p].Perm(int(width))
	}
	ts = tic()
	for p = 0; p < depth; p++ {
		go rowSumPermute(arr, p, ch, perm[p])
	}
	var tmp int64
	for p = 0; p < depth; p++ {
		tmp = <-ch
		sum += tmp
	}
	te = toc(ts)
	fmt.Printf("sum: %d\n", sum)
	fmt.Printf("time parallel perm: %s\n", te)
	sum = 0
	ts = tic()
	for p = 0; p < depth; p++ {
		rowSumPermute(arr, p, ch, perm[p])
		tmp = <-ch
		sum += tmp
	}
	te = toc(ts)
	fmt.Printf("sum: %d\n", sum)
	fmt.Printf("time serial perm: %s\n", te)
}
