package localstorage

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestGetParent(t *testing.T) {
	tab := [][2]int{
		{3, 1}, {4, 1}, {11, 5}, {13, 6}, {14, 6}, {12, 5}, {1, 0}, {2, 0},
	}
	for _, v := range tab {
		r := getParent(v[0])
		if v[1] != r {
			t.Error(v, r)
		}
	}

}

func TestGetLeftChild(t *testing.T) {
	tab := [][2]int{
		{1, 3}, {3, 7}, {0, 1}, {4, 9}, {5, 11}, {6, 13}, {2, 5},
	}
	for _, v := range tab {
		r := getLeftChild(v[0])
		if v[1] != r {
			t.Error(v, r)
		}
	}
}

func TestGetRightChild(t *testing.T) {
	tab := [][2]int{
		{1, 4}, {3, 8}, {0, 2}, {4, 10}, {5, 12}, {6, 14}, {2, 6},
	}
	for _, v := range tab {
		r := getRightChild(v[0])
		if v[1] != r {
			t.Error(v, r)
		}
	}
}

func TestAddNode(t *testing.T) {
	h := initHeap()
	t.Log("init----------")
	printHeap(h)
	if !checkHeap(t, h, 0) {
		printHeap(h)
	}
	// t.Log("Pop---------")
	// t.Log(h.Pop())
	// if !checkHeap(t, h, 0) {
	// 	printHeap(h)
	// }

	t.Log("SetNode---------")
	n := h.Peek()
	t.Log(n)
	n.expire = n.expire.Add(20 * time.Second)
	n.value = n.value.(int) + 20
	h.SetNode(&n)
	if !checkHeap(t, h, 0) {
		printHeap(h)
	}
	// printHeap(h)
}

func TestPop(t *testing.T) {
	for i := 0; i < 10; i++ {

		h := initHeap()
		h.AddNode(&node{
			expire: time.Now(),
			value:  10,
		})

		if !checkHeap(t, h, 0) {
			printHeap(h)
		}
		t.Log("----------------pop")

		min := h.Pop()
		if !checkHeap(t, h, 0) {
			printHeap(h)
		}
		t.Log("-----------", min.expire.Unix())
		// times := 0
		for !h.Empty() {
			c := h.Pop()
			if !checkHeap(t, h, 0) {
				printHeap(h)
			}
			if c.expire.Before(min.expire) {
				t.Error("pop err", min.expire.Unix(), c.expire.Unix())
			}
			min = c
		}

		for i := 0; i < 15; i++ {
			h.AddNode(&node{
				expire: time.Now(),
				value:  i,
				name:   "b",
			})
			if !checkHeap(t, h, 0) {
				printHeap(h)
			}
		}
		h.DeleteNode(5)
		if !checkHeap(t, h, 0) {
			printHeap(h)
		}

		n := h.Peek()
		n.expire = n.expire.Add(10 * time.Second)
		h.SetNode(&n)
		if !checkHeap(t, h, 0) {
			printHeap(h)
		}

		h.DeleteNode(7)
		if !checkHeap(t, h, 0) {
			printHeap(h)
		}
		time.Sleep(time.Second)
	}

}

func TestSetNode(t *testing.T) {
	h := initHeap()
	t.Log("count", h.Count())
	h.DeleteNode(4)

	t.Log("count", h.Count())
	n := &node{
		expire: time.Now(),
		value:  10,
		name:   "b",
	}
	h.AddNode(n)
	t.Log("count", h.Count())
	t.Log(n)

	n.expire = n.expire.Add(11 * time.Second)
	h.SetNode(n)
	t.Log(n)
	checkHeap(t, h, 0)

}

func initHeap() *MinHeap {
	h := NewMinHeap(1)
	now := time.Now()
	rand.Seed(time.Now().UnixNano())
	tab := []int{}
	for i := 0; i < 10; i++ {
		tab = append(tab, rand.Int()%100)
	}
	for _, tb := range tab {
		h.AddNode(&node{
			expire: now.Add(time.Duration(tb) * time.Second),
			value:  tb,
			name:   "a",
		})
	}
	return h
}

func checkHeap(t *testing.T, h *MinHeap, index int) bool {
	t.Helper()
	if index >= h.count {
		return true
	}
	cNode := h.heap[index]
	lIndex := getLeftChild(index)
	if lIndex >= h.count {
		return true
	}
	lNode := h.heap[lIndex]
	if cNode.expire.After(lNode.expire) {
		t.Error(cNode, lNode)
		return false
	}
	if !checkHeap(t, h, lIndex) {
		return false
	}

	rIndex := getRightChild(index)
	if rIndex >= h.count {
		return true
	}
	rNode := h.heap[rIndex]
	if cNode.expire.After(rNode.expire) {
		t.Error(cNode, rNode)
		return false
	}
	if !checkHeap(t, h, rIndex) {
		return false
	}
	return true
}

func printHeap(h *MinHeap) {
	fmt.Println("------------start---------------")
	for _, v := range h.heap {
		fmt.Println(v)
	}
	fmt.Println("-------------end--------------")
}

func BenchmarkAdd(b *testing.B) {
	h := initHeap()
	now := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := rand.Int() % 1000
		h.AddNode(&node{
			expire: now.Add(time.Duration(r) * time.Second),
			value:  nil,
			index:  i,
			name:   "a",
		})
	}
}
func BenchmarkSet(b *testing.B) {
	h := initHeap()
	now := time.Now()
	for i := 0; i < 100000; i++ {
		r := rand.Int() % 1000
		h.AddNode(&node{
			expire: now.Add(time.Duration(r) * time.Second),
			value:  nil,
			index:  i,
			name:   "a",
		})
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := rand.Int() % 1000
		n := h.Pop()
		n.expire.Add(time.Duration(r) * time.Second)
		h.SetNode(&n)
	}
}
