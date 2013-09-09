/* access.go : Implement strided memory access to see the effects of cache
Author: James Fairbanks
Date: 2012-09-02
Liscence: BSD
*/

package main

import (
	"fmt"
	"github.com/jpfairbanks/timing"
	"math/rand"
	"flag"
)

//Dtype: the data type of the arrays
type Dtype int64

//Strider: a slice of type dtype with a fixed stride
type Strider struct {
	data   []Dtype
	stride int
}

//Fill: fill a strided array so that we can walk it later
func (s Strider) Fill() {
	n := len(s.data)
	iterations := s.stride
	for k := 0; k < iterations; k++ {
		for i := 0; i < len(s.data)-k; i += 1 {
			s.data[(s.stride*i+k) % n] = Dtype((s.stride*(i+1) + k) % n)
		}
	}
	for i:=1; i < s.stride; i++{
		//fmt.Printf("%d ", i)
		s.data[(n-i)-1] = s.data[(n-i)-1] + 1
	}
	s.data[n-1] = 0
}

//Walk: walk a strided array by the stride
func (s Strider) Walk() Dtype {
	var v Dtype
	ptr := s.data[0]
	//fmt.Printf("stride: %d, ", s.stride)
	for ptr != 0 {
		v += ptr
		ptr = s.data[ptr]
	}
	return v
}

var scaleptr *uint
func init(){
	scaleptr = flag.Uint("scale", 10, "set the length of the array to access as a power of 2")
	flag.Parse()
}
//main: time.Time the Walk
func main() {
	var n, s int
	var scale uint
	scale = *scaleptr
	n = 2 << scale
	//fmt.Printf("scale:%d,", scale)
	s = 1
	k := 6
	stepsize := make([]int, k)
	tg := timing.New(k)
	correctsum := Dtype((n * (n - 1)) / 2)
	stepsize[0] = 1
	for i := 1; i < len(stepsize); i++ {
		stepsize[i] = 2 * stepsize[i-1]
	}
	//fmt.Printf("strides:%v,\n",stepsize)
	sdr := Strider{make([]Dtype, n), s}
	var sum Dtype
	for i, k := range stepsize {
		sdr.stride = k
		sdr.Fill()
		tg.Tic(i)
		for iter := 0; iter < 10; iter++ {
			sum = sdr.Walk()
		}
		tg.Toc(i)
		if  correctsum != sum {
			fmt.Println("we did not hit all the elements")
			fmt.Println(correctsum, sum)
		}
	}

	perm := rand.Perm(n)
	ptg := timing.New(1)
	ptg.Tic(0)
	var randsum Dtype
	for run:=0; run < 10; run++{
		randsum = 0
		for _, j := range perm {
			randsum += sdr.data[j]
		}
	}
	ptg.Toc(0)
	if randsum != correctsum {
		fmt.Printf("random ordered produced wrong sum: %v, %v,\n", randsum, correctsum)
	}
	tg.Resolve()
	fmt.Print(tg)
	ptg.Resolve()
	fmt.Printf(" %v\n", ptg)
}
