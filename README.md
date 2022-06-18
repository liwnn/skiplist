# Skip list implementation for Go

This Go package provides an implementation of skip list: [Skip Lists: A Probabilistic Alternative to Balanced Trees](https://www.cl.cam.ac.uk/teaching/2005/Algorithms/skiplists.pdf).


## Usage
All you have to do is to implement a comparison `function Less() bool` for your Item which will be store in the skip list, here are some examples.

A case for `int` items.
``` go
package main

import (
    "github.com/liwnn/skiplist"
)

func main() {
	sl := skiplist.New()

	// Insert some values
	for i := 0; i < 10; i++ {
		sl.Insert(skiplist.Int(i))
	}

	// Get the value of the key
	item := sl.Search(skiplist.Int(5))
	if item != nil {
		fmt.Println(item)
	}

	// Delete the key
	if sl.Delete(skiplist.Int(4)) {
		fmt.Println("Deleted", 4)
	}

	// Traverse the list
	for it := sl.NewIterator(); it.Valid(); it.Next() {
		fmt.Println(it.Value().(skiplist.Int))
	}
}
```

A case for `struct` items:
``` go
package main

import (
    "github.com/liwnn/skiplist"
)

type KV struct {
	Key   int
	Value int
}

func (kv KV) Less(than skiplist.Item) bool {
	return kv.Key < than.(KV).Key
}

func main() {
	sl := skiplist.New()

	// Insert some values
	for i := 0; i < 10; i++ {
		sl.Insert(KV{Key: i, Value: 100 + i})
	}

	// Get the value of the key
	item := sl.Search(KV{Key: 1})
	if item != nil {
		fmt.Println(item.(KV))
	}

	// Delete the key
	if sl.Delete(KV{Key: 4}) {
		fmt.Println("Deleted", 4)
	}

	// Traverse the list
	for it := sl.NewIterator(); it.Valid(); it.Next() {
		fmt.Println(it.Value().(KV))
	}
}
```
