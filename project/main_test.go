package main

import (
	"fmt"
	"math/rand"
	"testing"
)

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
