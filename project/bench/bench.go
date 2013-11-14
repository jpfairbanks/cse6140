/* bench.go : microbenchmarks to understand performance in golang
Author: James Fairbanks
Date: 2013-09-19
*/

package main

import (
	//"fmt"
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

func add_place(input []int, i int, ch chan int) {
	temp := input[i] + input[len(input)/2+i]
	input[i] = temp
	ch <- temp
}
func ParSum(xarr []int) int {
	length := len(xarr)
	if length == 1 {
		return xarr[0]
	} else {
		ch := make(chan int, length/2)
		for i := 0; i < length/2; i++ {
			go add_place(xarr, i, ch)
		}
		for i := 0; i < length/2; i++ {
			<-ch
		}
		newxarr := xarr[:length/2]
		return ParSum(newxarr)
	}
}

func main() {

}
