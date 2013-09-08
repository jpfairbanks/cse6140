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

func TestMultiInsert(t *testing.T) {
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
