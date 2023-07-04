package prefixtrie

import (
	"fmt"
	"testing"

	"github.com/mixcode-lib/arraymap"
)

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
	trie := NewRoot[rune, int]()

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
	trie.Traverse(func(prefix [][]rune, path []*Node[rune, int]) {
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

func TestParse(t *testing.T) {

	var err error

	type testT struct {
		S, Prefix string
	}

	testS := []testT{
		{"abcdef", "abcde"},
		{"abcdeg", "abcde"},
		{"abc", ""},
		{"ijkl", "ijk"},
		{"ijklm", "ijkl"},
		{"ijkxy", "ijk"},
	}
	testM := arraymap.New[string, string]()

	//root := &TrieNode[rune, string]{}
	root := NewRoot[rune, string]()

	// put test
	for _, e := range testS {
		err = root.Put([]rune(e.S), e.S)
		if err != nil {
			t.Fatal(err)
		}
		testM.Put(e.S, e.Prefix)
	}

	// traverse test
	root.Traverse(func(prefix [][]rune, path []*Node[rune, string]) {
		node := path[len(path)-1]
		if node.Payload == "" {
			t.Errorf("payload not found")
		}
		prefS := ""
		for _, s := range prefix {
			prefS += string(s)
		}
		key := prefS + string(node.Prefix)
		if key != node.Payload {
			t.Errorf("path and payload not match: [%v] [%v]", path, node.Payload)
		}
		pref, ok := testM.Get(key)
		if !ok {
			t.Errorf("path not in the test list: %v", prefS)
		}
		if prefS != pref {
			t.Errorf("prefix does not match for path %v, expected %v, actual %v", path, pref, prefS)
		}
	})

	// lookup test
	for _, e := range testS {
		n := root.Lookup([]rune(e.S))
		if n == nil || n.Payload != e.S {
			t.Errorf("lookup failed for key %v", e.S)
		}
		n = root.Lookup([]rune(e.S + "_"))
		if n != nil {
			t.Errorf("lookup must fail")
		}
		if !testM.HasKey(e.Prefix) {
			n = root.Lookup([]rune(e.Prefix))
			if n != nil {
				t.Errorf("lookup must fail")
			}
		}
	}

}
