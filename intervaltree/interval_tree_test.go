package intervaltree

import (
	"fmt"
	"testing"
)

func (n *node) isAVL() error {
	// Empty tree is always AVL
	if n == nil {
		return nil
	}

	// check heights consistency
	if n.getHeight() != max(n.Left.getHeight(), n.Right.getHeight())+1 {
		return fmt.Errorf("Height is wrong. Got '%d', expected '%d'", n.getHeight(), max(n.Left.getHeight(), n.Right.getHeight()))
	}

	bal := n.balanceFactor()
	if bal > 1 || bal < -1 {
		return fmt.Errorf("Tree is unbalanced. Balance factor = '%d'", bal)
	}

	if err := n.Left.isAVL(); err != nil {
		return err
	}

	return n.Right.isAVL()
}

func Test1(t *testing.T) {
	it := New()
	if err := it.root.isAVL(); err != nil {
		t.Fatalf("Empty tree is not AVL: %v", err)
	}
}

func Test2(t *testing.T) {
	it := New()
	it.Insert(1, 10)
	if err := it.root.isAVL(); err != nil {
		t.Fatalf("Tree with one node is not AVL: %v", err)
	}

	if it.Contains(0) {
		t.Fatal("Tree with one node contains non-inserted value 0")
	}

	if it.Contains(11) {
		t.Fatal("Tree with one node contains non-inserted value 11")
	}

	for i := 1; i < 11; i++ {
		if !it.Contains(uint64(i)) {
			t.Fatalf("Tree with one node does not contain inserted value %d", i)
		}
	}
}

func Test3(t *testing.T) {
	it := New()
	for i := 0; i < 10; i++ {
		if it.Contains(uint64(i)) {
			t.Fatalf("Empty IntervalTree contains %d", i)
		}
	}

	it.Insert(5, 15)
	for i := 5; i < 16; i++ {
		if !it.Contains(uint64(i)) {
			t.Fatalf("IntervalTree does not contain %d which was added", i)
		}
	}

	it.Insert(0, 1)
	if !it.Contains(0) {
		t.Fatalf("IntervalTree does not contain %d which was added", 0)
	}
	if !it.Contains(1) {
		t.Fatalf("IntervalTree does not contain %d which was added", 1)
	}

	it.Insert(2, 3)
	if !it.Contains(2) {
		t.Fatalf("IntervalTree does not contain %d which was added", 2)
	}
	if !it.Contains(3) {
		t.Fatalf("IntervalTree does not contain %d which was added", 3)
	}

	if it.Contains(4) {
		t.Fatalf("IntervalTree contains %d which was not added", 4)
	}

	it.Insert(40, 50)
	it.Insert(17, 19)
	it.Insert(21, 23)
	it.Insert(25, 29)
	if !it.Contains(18) {
		t.Fatalf("IntervalTree does not contain %d which was added", 18)
	}
	if !it.Contains(26) {
		t.Fatalf("IntervalTree does not contain %d which was added", 26)
	}
	if it.Contains(24) {
		t.Fatalf("IntervalTree contains %d which was not added", 24)
	}

	it.Insert(30, 39)
	for i := 25; i < 51; i++ {
		if !it.Contains(uint64(i)) {
			t.Fatalf("IntervalTree does not contain %d which was added", i)
		}
	}

}
