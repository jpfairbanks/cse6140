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
)

//Dtype: the data type of the arrays
type Dtype int64

//Strider: a slice of type dtype with a fixed stride
type Strider struct {
	data   []Dtype
	stride int
}

//Walk: walk a strided array by the stride
func (s Strider) Walk() Dtype {
	var v Dtype
	iterations := s.stride
	for k := 0; k < iterations; k++ {
		for i := 0; i < len(s.data)-k; i += s.stride {
			v += s.data[i+k]
		}
	}
	return v
}

//main: time.Time the Walk
func main() {
	var n, s int
	n = 2 << 20
	fmt.Println(n)
	s = 1
	k := 18
	size := make([]int, k)
	tg := timing.New(k)

	size[0] = 1
	for i := 1; i < len(size); i++ {
		size[i] = 2 * size[i-1]
	}
	fmt.Println(size)
	sdr := Strider{make([]Dtype, n), s}
	for i := 0; i < len(sdr.data); i++ {
		sdr.data[i] = Dtype(i)
	}

	for i, k := range size {
		sdr.stride = k
		tg.Tic(i)
		var sum Dtype
		for iter := 0; iter < 1; iter++ {
			sum = sdr.Walk()
		}
		tg.Toc(i)
		if correctsum := (n * (n - 1)) / 2; correctsum != int(sum) {
			fmt.Println("we did not hit all the elements")
			fmt.Println(correctsum, sum)
		}
	}

	perm := rand.Perm(n)
	ptg := timing.New(1)
	ptg.Tic(0)
	var v Dtype
	for _, j := range perm {
		v = sdr.data[j]
	}

	ptg.Toc(0)
	fmt.Println(v)
	tg.Resolve()
	fmt.Println(tg.TupleString("\n"))
	ptg.Resolve()
	fmt.Println("Random Access: ", ptg)
}
