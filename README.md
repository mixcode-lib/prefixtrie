# prefixtrie

A Prefix Trie in Go.


## Examples
```

func Example() {

    // key-value pairs of a string key an a int value
	values := []struct {
		S string
		N int
	}{
		{"abcdef", 2},
		{"abcdeg", 3},
		{"abc", 1},
		{"ijkl", 4},
		{"ijkxy", 6},
		{"ijklm", 5},
	}

    // make a trie of []rune key and int value
	trie := NewTrieRoot[rune, int]()

	// Put values to the trie
	for _, e := range values {
		err := trie.Put([]rune(e.S), e.N)
		if err != nil {
			panic(err)
		}
	}

	// find an entry
	l1 := trie.Lookup([]rune("abcdef"))
	fmt.Printf("%d\n", l1.Payload)

	l2 := trie.Lookup([]rune("abcdez"))
	if l2 != nil {
		panic("must be nil")
	}

	// Traverse the whole trie
	trie.Traverse(func(prefix [][]rune, path []*TrieNode[rune, int]) {
		for i := range prefix {
			if i != 0 {
				fmt.Print("|")
			}
			fmt.Printf("%s", string(prefix[i]))
		}
		node := path[len(path)-1]
		fmt.Printf("|%s : %d\n", string(node.Prefix), node.Payload)
	})

	// Output:
	// 2
	// |abc : 1
	// |abc|de|f : 2
	// |abc|de|g : 3
	// |ijk|l : 4
	// |ijk|l|m : 5
	// |ijk|xy : 6

}


```
