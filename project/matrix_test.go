package main

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

//updater: Perform atomic updates to row 1 of the matrix
func updater(M *Matrix, ch chan int64) {
	var i, j int64
	for i = 0; i < M.Ncols; i++ {
		for j = 0; j < 100; j++ {
			M.AtomicAdd(1, i, 1)
		}
	}
	ch <- 1
}

//raceyUpdater: Perform unsynchronized updates to row 0 of the matrix
func raceyUpdater(M *Matrix, ch chan int64) {
	var i, j int64
	for i = 0; i < M.Ncols; i++ {
		for j = 0; j < 100; j++ {
			M.Add(0, i, 1)
		}
	}
	ch <- 1
}
func TestRacey(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	M := NewMatrix(2, 4)
	var k, rounds int64
	rounds = 500
	fmt.Println(M)
	ch := make(chan int64, rounds)
	ts := time.Now()
	for k = 0; k < rounds; k++ {
		go updater(&M, ch)
	}
	for k = 0; k < rounds; k++ {
		<-ch
	}
	timeAtomic := time.Now().Sub(ts)
	var i int64
	for i = 0; i < M.Ncols; i++ {
		if M.Read(1, i) != rounds*100 {
			t.Errorf("Atomics did not work: 1,%d\n", i)
		}
	}
	t.Logf("Post Atomics:\n%v\n", M)
	ts = time.Now()
	for k = 0; k < rounds; k++ {
		go raceyUpdater(&M, ch)
	}
	for k = 0; k < rounds; k++ {
		<-ch
	}
	for i = 0; i < M.Ncols; i++ {
		if M.Read(0, i) == rounds*100 {
			t.Errorf("Race condition did not appear 0,%d\n", i)
		}
	}
	t.Logf("Post Racey:\n%v\n", M)
	timeRacey := time.Now().Sub(ts)
	fmt.Printf("time atomic: %d\n", timeAtomic)
	fmt.Printf("time racey: %d\n", timeRacey)
	fmt.Printf("atomic/racey: %d\n", timeAtomic/timeRacey)
}

func TestString(t *testing.T) {
	rightS := "[0 0 0 0]\n[0 0 0 0]"
	M := NewMatrix(2, 4)
	S := M.String()
	t.Logf("%s\n", S)
	if S != rightS {
		t.Errorf(S)
	}
}

//TestEqual: make sure that we can test for deep equality
func TestEqual(t *testing.T) {
	M := NewMatrix(2, 4)
	other := NewMatrix(2, 5)
	if M.Equal(&other) {
		t.Errorf("Failed on different sized matrices\n")
	}
	other = NewMatrix(2, 4)
	if !M.Equal(&other) {
		t.Errorf("False negative on zeros")
	}
	other.Data[0] = 1
	if M.Equal(&other) {
		t.Errorf("False positive")
	}
}
