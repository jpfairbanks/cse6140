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

func TestLoad(){
	rounds := 21
	tmr := timing.New(rounds-1)
	var k int
	for k = 1; k < rounds; k++  {
		count := 2 << uint(k)
		ch := make(chan int, count)
		tmr.Tic(k-1)
		for i := 0; i < count; i++ {
			go run(i,ch)
		}
		var sum int
		for i := 0; i < count; i++ {
			sum += <- ch
		}
		tmr.Toc(k-1)
	}
	tmr.Resolve()
	fmt.Println(tmr.TupleString("\n"))
	return
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
	rounds := 10
	numdata :=2*3*5*7*9
	tmr := timing.New(rounds-1)
	rightans := (numdata*(numdata-1))/2
	fmt.Println(rightans)
	for i := 1; i < rounds; i++ {
		tmr.Tic(i-1)
		j := Parfor(i, numdata)
		fmt.Println(j)
		tmr.Toc(i-1)
	}
	tmr.Resolve()
	fmt.Println(tmr.TupleString("\n"))
}