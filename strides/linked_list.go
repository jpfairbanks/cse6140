/* linked_list.go : Implement linked list with too many memory allocations
Author: James Fairbanks
Date: 2012-09-02
Liscence: BSD
*/

package list

import (
//    "fmt"
//    "github.com/jpfairbanks/timing"
//    "math/rand"
)

//Dtype: the data type of the arrays
type Dtype int64

var NULL *Node
//NULL := 0

//Node: Struct containing a datum and a pointer to the next Node.
type Node struct{
	Datum Dtype
	Next  *Node
}

//List: Struct for containing the head of the list.
type List struct{
	Head *Node
}

//New: Make a new list of Nodes.
func New() List {
	var ell List
	head := &Node{ Datum:0, Next:NULL }
	ell.Head = head
	return ell
}

//Insert: Insert a Node into the list.
func (ell List) Insert(pos int, element Dtype){
	node := ell.Head
	for i:=0; i < pos; i++ {
		node = node.Next
	}
	next := node.Next
	node.Next = &Node{Datum:element, Next:next}
}

//Remove: Remove a Node into the list.
func (ell List) Remove(pos int) Dtype{
	node := ell.Head
	for i:=0; i < pos; i++ {
		node = node.Next
	}
	next := node.Next
	node.Next = node.Next.Next
	return next.Datum
}

//Strider: a slice of type dtype with a fixed stride
type Strider interface{
    Walk() Dtype
}

