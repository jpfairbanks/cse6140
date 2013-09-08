package list

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	ell := New()
	n := ell.Head
	if n.Datum != 0 {
		t.Error("bad initial data")
	}
	if n.Next != NULL {
		t.Error("bad next pointer")
	}
}

func TestInsert(t *testing.T) {
	ell := New()
	ptr := ell.Head
	t.Logf("got ptr:%v", ptr)
	ell.Insert(0, 5)
	t.Logf("ell: %v\n", ell)
	t.Logf("made insert")
	head := ell.Head
	t.Logf("head: %v\n", head)
	if head.Datum != 5 {
		t.Errorf("insert on 0 failed on datum: %v", head.Datum)
	}
	if head.Next != ptr {
		t.Errorf("insert on 0 failed on pointer: %v", ptr)
	}
}

func TestSecondInsert(t *testing.T) {
	ell := New()
	ptr := ell.Head
	t.Logf("got ptr:%v", ptr)
	ell.Insert(1, 5)
	t.Logf("ell: %v\n", ell)
	t.Logf("made insert")
	node := ell.Head
	t.Logf("head: %v\n", node)
	if node.Datum != 0 {
		t.Errorf("insert on 1 failed on datum: %v", node.Datum)
	}
	if node.Next == nil {
		t.Errorf("insert on 1 failed on pointer: %v", ptr)
	}
	if node.Next.Datum != 5 {
		t.Errorf("insert failed to adjust prev.Next: %v", node.Next)
	}
	fmt.Println("about to print")
	fmt.Println(ell.String())
}

func TestBigList(t *testing.T) {
	ell := New()
	k := uint(15)
	size := 2 << k
	//Fill the list full of data
	for i := 0; i < size; i++ {
		ell.Insert(i, Dtype(i))
	}

	//Extract all of the data
	var node *Node
	var sum Dtype
	node = ell.Head
	for node != nil {
		sum += node.Datum
		node = node.Next
	}
	//Checking that we made the traversal correctly
	correctsum := (size * (size - 1)) / 2
	fmt.Printf("complete traversal: %v %v", sum, correctsum)
	if Dtype(correctsum) != sum {
		t.Errorf("incomplete traversal: %v %v", sum, correctsum)
	}
}
