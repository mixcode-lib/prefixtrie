package prefixtrie

import (
	"errors"

	"github.com/mixcode-lib/arraymap"
)

var (
	ErrInvalidNode = errors.New("invalid node")
)

// Trie, or a prefix tree, is a search tree that each node stores a common substring of keys.
type TrieNode[K comparable, P any] struct {
	Prefix   []K                                    // Common search string to this node
	Children *arraymap.ArrayMap[K, *TrieNode[K, P]] // Children nodes

	Payload    P    // payload for a key that exactly matches to this node.
	PayloadSet bool // true if Payload is valid. i.e. if PayloadSet==false then this node is a non-leaf intermediate node.
}

// Make a trie node with key type []K and payload type P.
// The return value of this function should be used as the root node of a Trie.
func NewTrieRoot[K comparable, P any]() *TrieNode[K, P] {
	return &TrieNode[K, P]{}
}

// Put a payload to the path of []K in a trie.
func (tr *TrieNode[K, P]) Put(path []K, payload P) (err error) {
	// find the common part
	l := 0
	for l < len(tr.Prefix) && l < len(path) && tr.Prefix[l] == path[l] {
		l++
	}
	if l == len(tr.Prefix) { // path has the same prefix with the current node
		path = path[l:]
		if len(path) == 0 { // exact match to the current node
			// Note that the payload is always overwritten
			tr.Payload, tr.PayloadSet = payload, true
			return nil
		}
		if tr.Children == nil {
			tr.Children = arraymap.New[K, *TrieNode[K, P]]()
		}
		child, ok := tr.Children.Get(path[0])
		if ok {
			// traverse down to a child
			return child.Put(path, payload)
		} else {
			// add a new node
			tr.Children.Put(path[0], &TrieNode[K, P]{
				Prefix:     path,
				Children:   nil,
				Payload:    payload,
				PayloadSet: true,
			})
		}
		return nil
	}
	if l < len(tr.Prefix) { // path and the current prefix matched partially

		// a node that inherits the children of the current node
		node1 := &TrieNode[K, P]{
			Prefix:     tr.Prefix[l:],
			Children:   tr.Children,
			Payload:    tr.Payload,
			PayloadSet: tr.PayloadSet,
		}

		s := path[l:]
		if len(s) == 0 {
			// New path is shorter than the current prefix.
			// Replace the current node with the new path.
			tr.Prefix = tr.Prefix[:l]
			tr.Children = arraymap.New[K, *TrieNode[K, P]]()
			tr.Children.Put(node1.Prefix[0], node1)
			tr.Payload, tr.PayloadSet = payload, true
			return nil
		}

		// a node for the new string
		node2 := &TrieNode[K, P]{
			Prefix:     s,
			Children:   nil,
			Payload:    payload,
			PayloadSet: true,
		}

		// update the current node
		tr.Prefix = tr.Prefix[:l]
		tr.Children = arraymap.New[K, *TrieNode[K, P]]()
		tr.Children.Put(node1.Prefix[0], node1)
		tr.Children.Put(s[0], node2)
		var zeroP P
		tr.Payload, tr.PayloadSet = zeroP, false // remove the payload
	}
	return nil
}

// Search for a node.
func (tr *TrieNode[K, P]) Lookup(path []K) *TrieNode[K, P] {
	l := 0
	for l < len(tr.Prefix) && l < len(path) && tr.Prefix[l] == path[l] {
		l++
	}
	if l < len(tr.Prefix) {
		// not match
		return nil
	}
	if len(path) == l {
		// exact match found
		if !tr.PayloadSet {
			// the node is an intermediate node
			return nil
		}
		return tr
	}
	if tr.Children == nil {
		// no more search path to continue
		return nil
	}
	path = path[l:]
	c, ok := tr.Children.Get(path[0])
	if !ok {
		// no more search path to continue
		return nil
	}
	return c.Lookup(path)
}

// Do a depth-first traverse on the Trie node.
// callback() is called for each node, with a slice of key prefixes and a path from the root to the current node.
func (tr *TrieNode[K, P]) Traverse(callback func(prefix [][]K, path []*TrieNode[K, P])) {
	prefix := make([][]K, 0)
	path := make([]*TrieNode[K, P], 0)
	tr.traverse(prefix, path, callback)
}

func (tr *TrieNode[K, P]) traverse(prefix [][]K, path []*TrieNode[K, P], callback func(prefix [][]K, path []*TrieNode[K, P])) {
	path = append(path, tr)
	if tr.PayloadSet {
		callback(prefix, path)
	}
	if tr.Children != nil && tr.Children.Len() != 0 {
		prefix = append(prefix, tr.Prefix)
		for _, c := range tr.Children.Value {
			c.traverse(prefix, path, callback)
		}
	}
}
