/* bench.go : microbenchmarks to understand performance in golang
Author: James Fairbanks
Date: 2013-09-19
*/

package main

import (
	"time"
)

func run(i int, ch chan int) {
	time.Sleep(10 * time.Millisecond)
	ch <- i
}

func cpusum(xarr []int, start int, stop int, ch chan int) {
	var sum int
	for i := start; i < stop; i++ {
		sum += i
	}
	ch <- sum
}

func memsum(xarr []int, start int, stop int, ch chan int) {
	var sum int
	for _, xi := range xarr[start:stop] {
		sum += xi
	}
	ch <- sum
	return
}

func Parfor(f func([]int, int, int, chan int), numproc int, numdata int, data []int) int {
	ch := make(chan int, numproc)
	var start int
	var stop int
	for i := 0; i < numproc; i++ {
		start = i * (numdata / numproc)
		stop = (i + 1) * (numdata / numproc)
		go f(data, start, stop, ch)
	}
	var finalsum int
	for i := 0; i < numproc; i++ {
		finalsum += <-ch
	}
	return finalsum
}

func main() {

}
