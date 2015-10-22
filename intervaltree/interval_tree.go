// Package intervaltree provides a very limited variation of IntervalTree using
// AVL trees.
package intervaltree

import (
	"fmt"
	"sync"
)

// IntervalTree represents an IntervalTree to which intervals can be added
// through Insert(x, y) and membership to an interval in the tree can be checked
// with Contains(x). This implementation does not support overlapping intervals
// nor common IntervalTree operations. The next uint64 not contained in the tree
// can be obtained with Next(x).
type IntervalTree struct {
	root *node
	sync.RWMutex
}

// node holds an interval [I, J] and pointers to nodes holding intervals lesser
// and greater than its own.
type node struct {
	I, J        uint64 // Interval bounds
	Left, Right *node  // Left and right children
	height      uint8  // Nodes on the longest path to a leaf (for AVL retracing)
}

// newNode returns a pointer to a new node to be added as a leaf.
func newNode(x, y uint64) *node {
	ret := &node{
		I:      x,
		J:      y,
		height: 1,
	}
	return ret
}

// insert adds the interval [x, y] to the tree. [x, y] cannot overlap with the
// current tree. If prunning can be done it will be done.
func (n *node) insert(x, y uint64, pRef **node) error {
	if x < n.I && y >= n.I {
		return OverlapError(n.I)
	} else if x >= n.I && x <= n.J {
		return OverlapError(x)
	}

	defer n.rebalance(pRef)

	// New interval is to the left of this nodes interval
	if y < n.I {
		if n.I == y+1 { // Neighbour, expand current interval
			if n.Left == nil {
				n.I = x
				return nil
			}

			// Check if we can join with a child interval
			if n.Left.J == x-1 { // Absorb our child
				n.I = n.Left.I
				n.Left = n.Left.Left
			} else { // Try to take child from our child
				g, err := n.Left.tryJoinGreatestFirst(x, &n.Left)
				if err != nil {
					return err
				}
				n.I = g
			}
			return nil
		}

		// Not neighbouring
		if n.Left == nil { // Create child
			n.Left = newNode(x, y)
			return nil
		}

		// We have a child, let it handle this interval
		return n.Left.insert(x, y, &n.Left)
	}

	// New interval is to the right of this nodes interval
	if n.J == x-1 { // Neighbour, expand current interval
		if n.Right == nil {
			n.J = y
			return nil
		}

		// Check if we can join with a child interval
		if n.Right.I == y+1 { // Absorb our child
			n.J = n.Right.J
			n.Right = n.Right.Right
		} else { // Try to take child from our child
			l, err := n.Right.tryJoinLeastFirst(y, &n.Right)
			if err != nil {
				return err
			}
			n.J = l
		}
		return nil
	}

	// Not neighbouring
	if n.Right == nil { // Create child
		n.Right = newNode(x, y)
		return nil
	}

	// We have a child, let it handle this interval
	return n.Right.insert(x, y, &n.Right)
}

// rebalance fixes AVL invariants violations by applying rotations.
func (n *node) rebalance(nRef **node) {
	bal := n.balanceFactor()
	if bal == 2 {
		if n.Left.balanceFactor() < 0 {
			n.preRotateRight()
		}
		n.rotateLeft(nRef)
	} else if bal == -2 {
		if n.Right.balanceFactor() > 0 {
			n.preRotateLeft()
		}
		n.rotateRight(nRef)
	}

	n.height = max(n.Left.getHeight(), n.Right.getHeight()) + 1
}

// balanceFactor calculates the balance factor for this node.
func (n *node) balanceFactor() int8 {
	return int8(n.Left.getHeight() - n.Right.getHeight())
}

// getHeight returns the number of nodes in the longest path to a leaf
func (n *node) getHeight() uint8 {
	if n == nil {
		return 0
	}
	return n.height
}

// tryJoinGreatestFirst starts a tryJoinGreatest invocation chain. The first
// case is special (nRef is not &p.Right), thats why this function exists.
func (n *node) tryJoinGreatestFirst(x uint64, nRef **node) (uint64, error) {
	if x <= n.J {
		return x, OverlapError(n.J)
	}
	if n.Right == nil {
		return x, nil
	}

	defer n.rebalance(nRef)
	return n.Right.tryJoinGreatest(x, n)
}

// tryJoinLeastFirst starts a tryJoinLeast invocation chain. The first case is
// special (nRef is not &p.Left), thats why this function exists.
func (n *node) tryJoinLeastFirst(y uint64, nRef **node) (uint64, error) {
	if y >= n.I {
		return y, OverlapError(n.I)
	}
	if n.Left == nil {
		return y, nil
	}

	defer n.rebalance(nRef)
	return n.Left.tryJoinLeast(y, n)
}

// tryJoinGreatest returns the lower endpoint of the greatest interval in the
// children of n if its upper endpoint is a neighbour of x and also removes this
// interval. Otherwise it returns x.
func (n *node) tryJoinGreatest(x uint64, p *node) (uint64, error) {
	defer n.rebalance(&p.Right)
	if n.Right == nil { // n is the greatest interval
		if x <= n.J {
			return x, OverlapError(n.J)
		}
		if n.J == x-1 { // n neighbours
			p.Right = n.Left
			return n.I, nil
		}
		return x, nil
	}
	return n.Right.tryJoinGreatest(x, n)
}

// tryJoinLeast returns the upper endpoint of the least interval in the children
// of n if its lower endpoint is a neighbour of x and also removes this
// interval. Otherwise it returns x.
func (n *node) tryJoinLeast(y uint64, p *node) (uint64, error) {
	defer n.rebalance(&p.Left)
	if n.Left == nil { // n is the least interval
		if y >= n.I {
			return y, OverlapError(n.I)
		}
		if n.I == y+1 { // n neighbours
			p.Left = n.Right
			return n.J, nil
		}
		return y, nil
	}
	return n.Left.tryJoinLeast(y, n)
}

// rotateLeft performs a left tree rotation.
func (n *node) rotateLeft(nRef **node) {
	pivot := n.Left
	n.Left = n.Left.Right
	pivot.Right = n
	*nRef = pivot
}

// rotateRight performs a right tree rotation.
func (n *node) rotateRight(nRef **node) {
	pivot := n.Right
	n.Right = n.Right.Left
	pivot.Left = n
	*nRef = pivot
}

// preRotateRight performs the first rotation in a LeftRight case
func (n *node) preRotateRight() {
	pivot := n.Left
	n.Left = pivot.Right
	pivot.Right = n.Left.Left
	n.Left.Left = pivot
}

// preRotateLeft performs the first rotation in a RightLeft case
func (n *node) preRotateLeft() {
	pivot := n.Right
	n.Right = pivot.Left
	pivot.Left = n.Right.Right
	n.Right.Right = pivot
}

// contains checks recursively if x is contained in this node or its children.
func (n *node) contains(x uint64) bool {
	if n == nil {
		return false
	}

	if n.I <= x && x <= n.J {
		return true
	} else if x < n.I {
		return n.Left.contains(x)
	} else {
		return n.Right.contains(x)
	}
}

// containingNode checks recursively for the node holding the interval that
// contains x and returns this node. If x is not contained it returns nil.
func (n *node) containingNode(x uint64) *node {
	if n == nil {
		return nil
	}

	if n.I <= x && x <= n.J {
		return n
	} else if x < n.I {
		return n.Left.containingNode(x)
	} else {
		return n.Right.containingNode(x)
	}
}

// max returns the greatest of two uint8
func max(a, b uint8) uint8 {
	if a > b {
		return a
	}
	return b
}

// print SPrints recursively the intervals contained in this tree
func (n *node) print() string {
	if n == nil {
		return ""
	}

	return n.Left.print() + fmt.Sprintf("[%d -- %d]", n.I, n.J) + n.Right.print()
}

// ToString returns a string representing all the intervals contained in the tree
func (t *IntervalTree) ToString() string {
	return t.root.print()
}

// Contains checks recursively if x is contained in this node or its children.
func (t *IntervalTree) Contains(x uint64) bool {
	t.RLock()
	defer t.RUnlock()
	return t.root.contains(x)
}

// Next returns the minimum value not contained in the tree that is greater or
// equal to x.
func (t *IntervalTree) Next(x uint64) uint64 {
	t.RLock()
	defer t.RUnlock()
	c := t.root.containingNode(x)
	if c == nil {
		return x
	}

	return c.J + 1
}

// Insert adds an interval to the tree. The interval cannot overlap with the
// tree. If prunning is possible it will be done.
func (t *IntervalTree) Insert(x, y uint64) error {
	if x > y {
		return InvalidIntervalError{x, y}
	}

	t.Lock()
	defer t.Unlock()
	if t.root == nil { // First interval
		t.root = newNode(x, y)
		return nil
	}

	return t.root.insert(x, y, &t.root)
}

// New returns a pointer to an empty IntervalTree.
func New() *IntervalTree {
	return &IntervalTree{}
}
