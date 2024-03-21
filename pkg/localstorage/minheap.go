package localstorage

import (
	"math"
)

type MinHeap struct {
	heap  []*node //节点存储
	count int     //当前节点个数
}

func NewMinHeap(initCapacity int) *MinHeap {
	return &MinHeap{
		heap: make([]*node, 0, initCapacity),
	}
}

//AddNode 添加节点
func (m *MinHeap) AddNode(n *node) {
	var currentIndex = 0
	if m.count >= len(m.heap) {
		m.heap = append(m.heap, n)
		currentIndex = len(m.heap) - 1
	} else {
		m.heap[m.count] = n
		currentIndex = m.count
	}
	m.count++

	n.index = currentIndex
	//节点上浮
	m.up(currentIndex)
}

//SetNode 修改节点(通过n.index修改),返回修改成功或失败
//	如果n.index > 索引范围 或者 n=nil则返回失败
func (m *MinHeap) SetNode(n *node) bool {
	if n == nil || n.index >= m.count {
		return false
	}
	beforeNode := m.heap[n.index]
	m.heap[n.index] = n
	if n.expire.After(beforeNode.expire) {
		// fmt.Println("------down----", n.index)
		m.down(n.index)
	} else {
		// fmt.Println("------up----", n.index)
		m.up(n.index)
	}
	return true
}

//DeleteNode 删除节点
func (m *MinHeap) DeleteNode(index int) node {
	var ret node
	if index >= m.count {
		return ret
	}

	n := m.heap[index]
	last := m.heap[m.count-1]
	m.swapNode(last, n)
	m.heap[n.index] = nil
	m.count--
	m.down(last.index)
	return *n
}

//Peek 查看最小节点,不移除节点
func (m *MinHeap) Peek() node {
	var n node
	if m.Empty() {
		return n
	}

	n = *m.heap[0]
	return n
}

//Pop 查看并移除最小节点
func (m *MinHeap) Pop() node {
	var n node
	if m.Empty() {
		return n
	}
	firstNode := m.heap[0]
	lastNode := m.heap[m.count-1]
	// fmt.Println("first:", firstNode.index, " last:", lastNode.index)
	m.swapNode(firstNode, lastNode)
	m.count--
	m.heap[m.count] = nil
	m.down(0)
	return *firstNode
}

//Empty 是否为空
func (m *MinHeap) Empty() bool {
	return m.count == 0
}

//Count 返回元素个数
func (m *MinHeap) Count() int {
	return m.count
}

//up 节点上浮操作
func (m *MinHeap) up(nodeIndex int) {
	if nodeIndex == 0 {
		return
	}
	currentIndex := nodeIndex
	for {
		if currentIndex <= 0 { //已经是根节点
			break
		}
		cNode := m.heap[currentIndex]
		pIndex := getParent(currentIndex)
		pNode := m.heap[pIndex]
		if pNode.expire.Before(cNode.expire) { //父节点已经比当前节点小了
			break
		}
		m.swapNode(cNode, pNode)
		currentIndex = pIndex
	}
}

//swapNode 交换node,返回值是否交换成功
//	当n, n1中任意一个为nil时  或者  index超过count时  或者 n=n1时返回false
func (m *MinHeap) swapNode(n, n1 *node) bool {
	if n == nil || n1 == nil ||
		n == n1 ||
		n.index > m.count || n1.index > m.count {
		return false
	}
	m.heap[n.index] = n1
	m.heap[n1.index] = n
	n.index, n1.index = n1.index, n.index //交换索引
	return true
}

func (m *MinHeap) down(nodeIndex int) {
	if nodeIndex >= m.count {
		return
	}

	currentIndex := nodeIndex

	for {
		cNode := m.heap[currentIndex]
		minNode := cNode
		minIndex := currentIndex

		leftIndex := getLeftChild(currentIndex)
		if leftIndex < m.count {
			lNode := m.heap[leftIndex]
			if lNode.expire.Before(minNode.expire) {
				minNode = lNode
				minIndex = leftIndex
			}
		}

		rightIndex := getRightChild(currentIndex)
		if rightIndex < m.count {
			rNode := m.heap[rightIndex]
			if rNode.expire.Before(minNode.expire) {
				minNode = rNode
				minIndex = rightIndex
			}
		}

		if cNode == minNode { //如果最小节点是当前节点则不需要更换
			break
		}
		currentIndex = minIndex

		m.swapNode(cNode, minNode)
	}
}

//getLeftChild 通过父节点index获取左孩子节点index
func getLeftChild(pIndex int) int {
	return pIndex*2 + 1
}

//getRightChild 通过父节点获取右孩子节点
func getRightChild(pIndex int) int {
	return pIndex*2 + 2
}

//getParent 通过子节点或者父节点index
func getParent(cIndex int) int {
	var index = (float64(cIndex) - 1) / 2
	return int(math.Floor(index))
}
