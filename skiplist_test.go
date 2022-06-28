package skiplist

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func init() {
	seed := time.Now().Unix()
	rand.Seed(seed)
}

// perm returns a random permutation of n Int items in the range [0, n).
func perm(n int) (out []Item) {
	out = make([]Item, 0, n)
	for _, v := range rand.Perm(n) {
		out = append(out, Int(v))
	}
	return
}

// rang returns an ordered list of Int items in the range [0, n).
func rang(n int) (out []Item) {
	for i := 0; i < n; i++ {
		out = append(out, Int(i))
	}
	return
}

func TestSkipList(t *testing.T) {
	sl := New()
	const listSize = 10000
	for i := 0; i < 10; i++ {
		for _, item := range perm(listSize) {
			sl.Insert(item)
		}
		if sl.Len() != listSize {
			t.Fatal("insert failed", listSize, sl.Len())
		}
		for _, item := range perm(listSize) {
			if sl.Search(item) == nil {
				t.Fatal("has did not find item", item)
			}
		}
		for _, item := range perm(listSize) {
			sl.Insert(item)
		}
		it := sl.NewIterator()
		if min, want := it.Value(), Item(Int(0)); min != want {
			t.Fatalf("min: want %+v, got %+v", want, min)
		}

		for _, item := range perm(listSize) {
			if !sl.Delete(item) {
				t.Fatalf("didn't find %v", item)
			}
		}
	}
}

func ExampleSkipList() {
	sl := New()
	for i := Int(0); i < 10; i++ {
		sl.Insert(i)
	}
	fmt.Println("len:       ", sl.Len())
	fmt.Println("search3:   ", sl.Search(Int(3)))
	fmt.Println("search100: ", sl.Search(Int(100)))
	fmt.Println("del4:      ", sl.Delete(Int(4)))
	fmt.Println("del100:    ", sl.Delete(Int(100)))
	sl.Insert(Int(5))
	sl.Insert(Int(100))
	fmt.Println("len:       ", sl.Len())
	fmt.Printf("for:        ")
	for it := sl.NewIterator(); it.Valid(); it.Next() {
		fmt.Print(it.Value().(Int))
		fmt.Print(" ")
	}
	fmt.Println()
	// Output:
	// len:        10
	// search3:    3
	// search100:  <nil>
	// del4:       true
	// del100:     false
	// len:        10
	// for:        0 1 2 3 5 6 7 8 9 100
}

func TestIterator(t *testing.T) {
	sl := New()
	for _, v := range perm(100) {
		sl.Insert(v)
	}

	var got = make([]Item, 0, 100)
	for it := sl.NewIterator(); it.Valid(); it.Next() {
		got = append(got, it.Value())
	}

	if want := rang(100); !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}

	{
		it := sl.NewIterator()
		it.MoveTo(Int(20))
		if !it.Valid() || it.Value() != Int(20) {
			t.Fatal("iterator didn't move to 100")
		}
	}
}

func TestRange(t *testing.T) {
	sl := New()
	for _, v := range rang(10) {
		sl.Insert(v)
	}

	var got = make([]Item, 0, 10)
	for rang := sl.NewRange(Int(1), Int(3)); !rang.End(); rang.Next() {
		got = append(got, rang.Value())
	}
	var want = []Item{Int(1), Int(2), Int(3)}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

const benchmarkListSize = 10000

func BenchmarkInsert(b *testing.B) {
	b.StopTimer()
	insertP := perm(benchmarkListSize)
	b.StartTimer()
	i := 0
	for i < b.N {
		sl := New()
		for _, item := range insertP {
			sl.Insert(item)
			i++
			if i >= b.N {
				return
			}
		}
	}
}

func BenchmarkSearch(b *testing.B) {
	b.StopTimer()
	insertP := perm(benchmarkListSize)
	searchP := perm(benchmarkListSize)
	b.StartTimer()
	i := 0
	for i < b.N {
		b.StopTimer()
		tr := New()
		for _, v := range insertP {
			tr.Insert(v)
		}
		b.StartTimer()
		for _, item := range searchP {
			tr.Search(item)
			i++
			if i >= b.N {
				return
			}
		}
	}
}

func BenchmarkDeleteInsert(b *testing.B) {
	b.StopTimer()
	insertP := perm(benchmarkListSize)
	sl := New()
	for _, item := range insertP {
		sl.Insert(item)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sl.Delete(insertP[i%benchmarkListSize])
		sl.Insert(insertP[i%benchmarkListSize])
	}
}

func BenchmarkDelete(b *testing.B) {
	b.StopTimer()
	insertP := perm(benchmarkListSize)
	removeP := perm(benchmarkListSize)
	b.StartTimer()
	i := 0
	for i < b.N {
		b.StopTimer()
		sl := New()
		for _, v := range insertP {
			sl.Insert(v)
		}
		b.StartTimer()
		for _, item := range removeP {
			sl.Delete(item)
			i++
			if i >= b.N {
				return
			}
		}
		if sl.Len() > 0 {
			panic(sl.Len())
		}
	}
}
