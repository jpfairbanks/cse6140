/* bench.go : microbenchmarks to understand performance in golang
Author: James Fairbanks
Date: 2013-09-19
*/

package main

import (
	"fmt"
	"github.com/jpfairbanks/timing"
	"runtime"
	"testing"
)

func tmrKeyPrint(tmr timing.Timing, key string) {
	for index, dura := range tmr.Td {
		fmt.Printf("%s:%d %d\n", key, index, dura.Nanoseconds())
	}
}

func TestLoad(t *testing.T) {
	rounds := 21
	tmr := timing.New(rounds - 1)
	var k int
	for k = 1; k < rounds; k++ {
		count := 2 << uint(k)
		ch := make(chan int, count)
		tmr.Tic(k - 1)
		for i := 0; i < count; i++ {
			go run(i, ch)
		}
		var sum int
		for i := 0; i < count; i++ {
			sum += <-ch
		}
		tmr.Toc(k - 1)
	}
	tmr.Resolve()
	tmrKeyPrint(tmr, "Load")
	return
}

const rounds = 10
const numdata = 2 * 3 * 5 * 7 * 9 * 128 * 4
const reps = 5
const rightans = (numdata * (numdata - 1)) / 2

func TestParforMem(t *testing.T) {
	data := make([]int, numdata)
	fmt.Printf("numdata:%d\n", numdata)
	for i := 0; i < numdata; i++ {
		data[i] = i
	}
	tmr := timing.New(rounds - 1)
	var ans int
	for i := 1; i < rounds; i++ {
		runtime.GOMAXPROCS(i)
		tmr.Tic(i - 1)
		for k := 0; k < reps; k++ {
			ans = Parfor(memsum, i, numdata, data)
		}
		tmr.Toc(i - 1)
		if ans != rightans {
			t.Fail()
		}
	}
	tmr.Resolve()
	tmrKeyPrint(tmr, "ParforMem")
}

func TestParforCPU(t *testing.T) {
	var data []int
	tmr := timing.New(rounds - 1)
	var ans int
	for i := 1; i < rounds; i++ {
		runtime.GOMAXPROCS(i)
		tmr.Tic(i - 1)
		for k := 0; k < reps; k++ {
			ans = Parfor(cpusum, i, numdata, data)
		}
		tmr.Toc(i - 1)
		if ans != rightans {
			t.Fail()
		}
	}
	tmr.Resolve()
	tmrKeyPrint(tmr, "ParforCPU")
}
