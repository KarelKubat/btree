package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/KarelKubat/btree"
)

// The payload of a node: a string, and how many times it was seen.
type stringcount struct {
	str   string
	count int64
}

// Visitor function: prints the number of times a word was seen, and the word itself.
func nodeWalk(n *btree.Node) {
	fmt.Println(n.Payload.(*stringcount).count, n.Payload.(*stringcount).str)
}

// Node comparison :`a` is "less" if its string is alphabetically less.
func lessFunc(a, b *btree.Node) bool {
	return a.Payload.(*stringcount).str < b.Payload.(*stringcount).str
}

func main() {
	// Check cmdline, open input file or set input to stdin
	if len(os.Args) != 1 {
		log.Fatalln("Usage: wordcount (reads from stdin, shows words and their frequencies)")
	}

	// Instantiate a binary tree.
	bt := btree.New(lessFunc)

	// Start a scanner that splits by spaces.
	sc := bufio.NewScanner(os.Stdin)
	sc.Split(bufio.ScanWords)
	for sc.Scan() {
		// Insert or find node having a `stringcount` payload with the word. If the node is inserted
		// as fresh, then its count will be zero. If the node was found already in the tree, then its
		// count will be something else. In any case we increment the count.
		// The second return value from `bt.Upsert()` is a boolean indicating whether the node was
		// added to the tree. In this situation we don't care.
		intree, _ := bt.Upsert(&btree.Node{Payload: &stringcount{str: sc.Text()}})
		intree.Payload.(*stringcount).count++

		// Alternatively, one might:
		// intree, inserted := bt.Upsert(&btree.Node{Payload: &stringcount{str: sc.Text(), count: 1}})
		// if !inserted {
		//	 intree.Payload.(*stringcount).count++
		//}
	}
	bt.DepthFirstInOrder(nodeWalk)
	// In reverse order you might use: bt.DepthFirstReverse(nodeWalk)
}
