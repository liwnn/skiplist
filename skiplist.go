package skiplist

import (
	"math/rand"
)

const (
	kMaxLevel = 12   // (1/p)^MaxLevel >= maxNode
	kP        = 0.25 // Skiplist P = 1/4

	DefaultFreeListSize = 32
)

var (
	nilNodes = make([]*node, 16)
)

type Item interface {
	Less(than Item) bool
}

// node is an element of a skip list
type node struct {
	item    Item
	forward []*node
}

type FreeList struct {
	freelist []*node
}

func NewFreeList(size int) *FreeList {
	return &FreeList{freelist: make([]*node, 0, size)}
}

func (f *FreeList) newNode(lvl int) (n *node) {
	index := len(f.freelist) - 1
	if index < 0 {
		n = new(node)
		n.forward = make([]*node, lvl)
		return
	}
	n = f.freelist[index]
	f.freelist[index] = nil
	f.freelist = f.freelist[:index]

	if cap(n.forward) < lvl {
		n.forward = make([]*node, lvl)
	} else {
		n.forward = n.forward[:lvl]
	}
	return
}

func (f *FreeList) freeNode(n *node) (out bool) {
	// for gc
	n.item = nil
	toClear := n.forward
	for len(toClear) > 0 {
		toClear = toClear[copy(toClear, nilNodes):]
	}

	if len(f.freelist) < cap(f.freelist) {
		f.freelist = append(f.freelist, n)
		out = true
	}
	return
}

// SkipList implemente "Skip Lists: A Probabilistic Alternative to Balanced Trees"
type SkipList struct {
	header   *node
	level    int // current max level
	freelist *FreeList
	length   int
}

// New creates a skip list
func New() *SkipList {
	sl := &SkipList{
		level:    1,
		freelist: NewFreeList(DefaultFreeListSize),
		header: &node{
			forward: make([]*node, kMaxLevel),
		},
	}
	return sl
}

// Search for an element by traversing forward pointers
func (sl *SkipList) Search(key Item) Item {
	x := sl.header
	// loop : x→key < searchKey <= x→forward[i]→key
	for i := sl.level - 1; i >= 0; i-- {
		for x.forward[i] != nil && x.forward[i].item.Less(key) {
			x = x.forward[i]
		}
	}
	x = x.forward[0]
	if x != nil && !key.Less(x.item) {
		return x.item
	}
	return nil
}

// Insert adds the given item to the skip list.
func (sl *SkipList) Insert(item Item) {
	if item == nil {
		panic("nil item being added to SkipList")
	}
	var prev = make([]*node, kMaxLevel)
	x := sl.header
	for i := sl.level - 1; i >= 0; i-- {
		for x.forward[i] != nil && x.forward[i].item.Less(item) {
			x = x.forward[i]
		}
		prev[i] = x
	}
	x = x.forward[0]
	if x != nil && !item.Less(x.item) {
		x.item = item
	} else {
		lvl := sl.randomLevel()
		if lvl > sl.level {
			for i := sl.level; i < lvl; i++ {
				prev[i] = sl.header
			}
			sl.level = lvl
		}

		x = sl.freelist.newNode(lvl)
		x.item = item
		for i := 0; i < lvl; i++ {
			x.forward[i] = prev[i].forward[i]
			prev[i].forward[i] = x
		}
		sl.length++
	}
}

// Delete remote an item equal to the passed in item. return true if success, else false.
func (sl *SkipList) Delete(item Item) bool {
	var prev = make([]*node, kMaxLevel)
	x := sl.header
	for i := sl.level - 1; i >= 0; i-- {
		for x.forward[i] != nil && x.forward[i].item.Less(item) {
			x = x.forward[i]
		}
		prev[i] = x
	}
	x = x.forward[0]
	if x != nil && !item.Less(x.item) {
		for i := 0; i < sl.level; i++ {
			if prev[i].forward[i] != x {
				break
			}
			prev[i].forward[i] = x.forward[i]
		}
		for sl.level > 1 && sl.header.forward[sl.level-1] == nil {
			sl.level--
		}
		sl.freelist.freeNode(x)
		sl.length--
		return true
	}
	return false
}

func (sl *SkipList) randomLevel() int {
	lvl := 1
	for lvl < kMaxLevel && rand.Float64() < kP {
		lvl++
	}
	return lvl
}

func (sl *SkipList) Len() int {
	return sl.length
}

func (sl *SkipList) NewIterator() *Iterator {
	return &Iterator{sl: sl, x: sl.header.forward[0]}
}

type Iterator struct {
	sl *SkipList
	x  *node
}

func (it *Iterator) Valid() bool {
	return it.x != nil
}

func (it *Iterator) Next() {
	it.x = it.x.forward[0]
}

func (it *Iterator) Value() Item {
	return it.x.item
}

type Int int

// Less returns true if int(a) < int(b).
func (a Int) Less(b Item) bool {
	return a < b.(Int)
}
