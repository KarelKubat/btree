# btree

Package `btree` implements an in-memory binary tree, where nodes are stored in-order.

## API

### What's in the node?

The payload of a node is user-supplied. It should be a `struct` with all necessary fields.
An example of a payload to store a person's name and how often it occurs (say in some directory), may be:

```go
type person struct {
    name    string  // person's name
    counter int     // how many times this name is seen
}
 ```
Nodes in the tree will have this data as an amorph `interface{}` named `Payload`. Each node has furthermore the pointers `Left` and `Right` to sub-nodes, which may be examined. 

### Instantiating a binary tree

`btree.New()` returns an empty binary tree. The argument to `New()` is a comparison function that `btree` uses to find the right place for nodes. For example, if the stored persons are ordered by the  name:

```go
// lessFunc must conform to the type btree.LessFunc
func lessFunc(a, b *btree.Node) bool {
    return a.Payload.(*person).name < b.Payload.(*person).name
}
...
bt := btree.New(lessFunc)
```

The field `Root` of the structure (in this example `bt`) is the top node. This field is `nil` until the first node is added.

### Adding nodes to the tree

Nodes are added using `btree.Upsert()`.

- The argument is a pointer to a `btree.Node`, with the field `Payload` set to the data to store.
- There are two return values:
  - The spot in the tree where the node resides (actually a pointer, so a `*btree.Node`)
  - A `bool` which is `true` when the node was added, or `false` when the node was already present.
- The returned `bool` can be used if, e.g., already-present nodes needs updating. In that case the first return value can be used to get at the `Payload`, which is then typecast.

Given the above `person` example this may be:

```go
node := &btree.Node{
    Payload: &person{
        name:    "John Smith", // name of the person in my address book
        counter: 1,            // assume it's the first time it's entered into the tree
    },
}
storageNode, inserted := bt.Upsert(node)
if !inserted {
    // this name was seen before
    storageNode.Payload.(*person).counter++
}

/* This would of course be more efficient without the `if`. We can just as well leave
   the counter to its default (zero), and always increment: first time from zero to
   one, or if the name occurs more times: from whatever it was to +1.
storageNode, _ := bt.Upsert(&btree.Node{
    Payload: &person {
        name: "Sponge Bob",
    },
})
storageNode.Payload.(*person).counter++ 
*/
```

### Examining the tree

Method `btree.DepthFirstInOrder()` "walks" the tree and activates a supplied callback:

```go
// printNode must conform to the type btree.WalkFunc
func printNode(btree.Node *n) {
    thisPerson := node.Payload.(*person)
    fmt.Printf("name %q was seen %v times\n", thisPerson.name, thisPerson.count)
}
...
bt.DepthFirstInOrder(printPerson)
```

Method `btree.DepthFirstReverse()` traverses the tree in reverse order.

## Full example (see `main/wordcount.go`)

```go
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
```