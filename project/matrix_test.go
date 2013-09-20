package main

import (
	"fmt"
	"testing"
	"time"
)

func TestRacey(t *testing.T) {
	M := NewMatrix(2, 4)
	var k, rounds int64
	rounds = 900
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
	fmt.Println(M)
	ts = time.Now()
	for k = 0; k < rounds; k++ {
		go raceyUpdater(&M, ch)
	}
	for k = 0; k < rounds; k++ {
		<-ch
	}
	timeRacey := time.Now().Sub(ts)
	fmt.Println(M)
	fmt.Printf("time atomic: %d\n", timeAtomic)
	fmt.Printf("time racey: %d\n", timeRacey)
	fmt.Printf("atomic/racey: %d\n", timeAtomic/timeRacey)
}
