/* linked_list.go : Implement linked list with too many memory allocations
Author: James Fairbanks
Date: 2012-09-02
Liscence: BSD
*/

package list

import (
	"fmt"
	"strings"

//    "github.com/jpfairbanks/timing"
//    "math/rand"
)

//Dtype: the data type of the arrays
type Dtype int64

var NULL *Node

//NULL := 0

//Node: Struct containing a datum and a pointer to the next Node.
type Node struct {
	Datum Dtype
	Next  *Node
}

//List: Struct for containing the head of the list.
type List struct {
	Head *Node
}

//New: Make a new list of Nodes.
func New() List {
	var ell List
	head := &Node{Datum: 0, Next: NULL}
	ell.Head = head
	return ell
}

//Insert: Insert a Node into the list.
func (ell *List) Insert(pos int, element Dtype) {
	var prev *Node
	node := ell.Head
	for i := 0; i < pos; i++ {
		//fmt.Println("getting the next")
		if i == pos-1 {
			prev = node
		}
		node = node.Next
	}
	//fmt.Printf("found pos node: %v\n", node)
	//fmt.Printf("found prev: %v\n", prev)
	newNode := &Node{Datum: element, Next: node}
	//fmt.Printf("made newNode: %v\n", newNode)
	if pos > 0 {
		prev.Next = newNode
	} else {
		ell.Head = newNode
	}
	//fmt.Println("did previous")
	//fmt.Println(node.Next)
}

//Remove: Remove a Node into the list.
func (ell *List) Remove(pos int) Dtype {
	node := ell.Head
	for i := 0; i < pos; i++ {
		node = node.Next
	}
	next := node.Next
	node.Next = node.Next.Next
	return next.Datum
}

//String: print  list by printing each element on a line
func (ell *List) String() string {
	var curr *Node
	curr = ell.Head
	var repr map[int]string
	repr = make(map[int]string)
	var sterm string
	var i int
	for curr != nil {
		sterm = fmt.Sprintf("%d %v %v",
			i, curr.Datum, curr.Next)
		//fmt.Println(sterm)
		repr[i] = sterm
		i++
		curr = curr.Next
	}
	ordered := make([]string, len(repr))
	for j := 0; j < i; j++ {
		ordered[j] = repr[j]
	}
	var joined string
	joined = strings.Join(ordered, "\n")
	//fmt.Println(joined)
	return joined
}

//Strider: a slice of type dtype with a fixed stride
type Strider interface {
	Walk() Dtype
}
