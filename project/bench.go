/* bench.go : microbenchmarks to understand performance in golang
Author: James Fairbanks
Date: 2013-09-19
*/

package main

import (
	"fmt"
	"time"
	"github.com/jpfairbanks/timing"
)

func run(i int, ch chan int) {
	time.Sleep(10*time.Millisecond)
	ch <- i
}

func Parfor(numproc int, numdata int) int{
	ch := make(chan int, numproc)
	data := make([]int, numdata)
	for i,_ := range data{
		data[i] = i
	}
	var start int
	var stop int
	run := func(xarr []int, ch chan int) {
		var sum int; for _, xi := range xarr{
		 sum += xi}
		  ch <- sum 
		  return
		}
	for i := 0; i < numproc; i++ {
		start = i*(numdata/numproc)
		stop = (i+1)*(numdata/numproc)
		go run(data[start:stop], ch)
	}
	var j int
	for i := 0; i < numproc; i++ {
		j += <-ch
	}
	return j
}

func main(){

}