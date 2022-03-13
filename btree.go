// Package btree implements a binary tree.
package btree

// LessFunc must be supplied by the caller of `Upsert()`. It is responsible for comparing two nodes
// `a` and `b` and must return `true` when `a` is "smaller".
type LessFunc func(a, b *Node) bool

// WalkFunc must be supplied by the caller of traversal functions such as `DepthFirstInOrder()`.
// `btree` will activate this callback for every node in the binary tree.
type WalkFunc func(n *Node)

// Node defines what is stored in a binary tree.
type Node struct {
	// Payload is an amorph placeholder that can be filled in case-by-case by the caller.
	Payload interface{}
	// Left and Right are next `Node`s. The fields are exported so that callers may easily
	// manipulate binary trees themselves.
	Left, Right *Node
}

// BTree holds a binary tree.
type BTree struct {
	// Root is the tree's root.
	Root *Node
	// Less is the `LessFunc` that is caller-supplied. It is repeatedly called when inserting.
	Less LessFunc
}

// New instantiates a new `BTree`.
func New(less LessFunc) *BTree {
	return &BTree{
		Less: less,
	}
}

// Upsert examines the tree and if needed, inserts a new node. The return value `intree` points
// to where the node was inserted (or where a previously inserted node was already found). The
// return value `inserted` is `true` when the node was added to the tree.
func (b *BTree) Upsert(n *Node) (intree *Node, inserted bool) {
	if b.Root == nil {
		b.Root = n
		return b.Root, true
	}
	return b.upsertFrom(b.Root, n)
}

func (b *BTree) upsertFrom(from, n *Node) (intree *Node, inserted bool) {
	switch {
	case b.Less(from, n):
		if from.Left == nil {
			from.Left = n
			return from.Left, true
		}
		return b.upsertFrom(from.Left, n)
	case b.Less(n, from):
		if from.Right == nil {
			from.Right = n
			return from.Right, true
		}
		return b.upsertFrom(from.Right, n)
	default:
		return from, false
	}
}

// DepthFirstInOrder "walks" along the tree and calls the `WalkFunc` for each node. Nodes are
// visited depth first, in order.
func (b *BTree) DepthFirstInOrder(walk WalkFunc) {
	if b.Root == nil {
		return
	}
	b.depthFirstInOrderFrom(b.Root, walk)
}

func (b *BTree) depthFirstInOrderFrom(n *Node, walk WalkFunc) {
	if n.Left != nil {
		b.depthFirstInOrderFrom(n.Left, walk)
	}
	walk(n)
	if n.Right != nil {
		b.depthFirstInOrderFrom(n.Right, walk)
	}
}

// DepthFirstReverse "walks" along the tree and calls the `WalkFunc` for each node. Nodes are
// visited depth first, reverse order.
func (b *BTree) DepthFirstReverse(walk WalkFunc) {
	if b.Root == nil {
		return
	}
	b.depthFirstReverseFrom(b.Root, walk)
}

func (b *BTree) depthFirstReverseFrom(n *Node, walk WalkFunc) {
	if n.Right != nil {
		b.depthFirstInOrderFrom(n.Right, walk)
	}
	walk(n)
	if n.Left != nil {
		b.depthFirstInOrderFrom(n.Left, walk)
	}
}
