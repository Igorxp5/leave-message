package queue

import (
	"errors"
	"fmt"
	"sync"
)

const Max uint64 = ^uint64(0)

type UniqueQueue struct {
	firstNode *Node
	lastNode  *Node
	max       uint64
	size      uint64
	nodeMap   map[string]*Node
	mu        sync.Mutex
}

type Node struct {
	index    uint64
	value    string
	previous *Node
	next     *Node
}

func New(max uint64) UniqueQueue {
	nodeMap := make(map[string]*Node)
	return UniqueQueue{max: max, nodeMap: nodeMap}
}

func (queue *UniqueQueue) Add(item string) (uint64, error) {
	defer queue.mu.Unlock()
	queue.mu.Lock()
	if queue.isIn(item) {
		return 0, errors.New("cannot add same item")
	}
	if queue.size == queue.max {
		return 0, errors.New("max queue size reached")
	}
	node := &Node{queue.size, item, nil, nil}
	if queue.firstNode == nil {
		queue.firstNode = node
	} else {
		node.previous = queue.lastNode
		queue.lastNode.next = node
	}
	queue.nodeMap[item] = node
	queue.lastNode = node
	queue.size += 1
	return node.index, nil
}

func (queue *UniqueQueue) First() (string, error) {
	defer queue.mu.Unlock()
	queue.mu.Lock()
	if queue.size == 0 {
		return "", errors.New("queue is empty")
	}
	return queue.firstNode.value, nil
}

func (queue *UniqueQueue) Index(item string) (uint64, error) {
	defer queue.mu.Unlock()
	queue.mu.Lock()
	if !queue.isIn(item) {
		return 0, errors.New(fmt.Sprintf("item \"%s\" is not in the queue", item))
	}
	node := queue.nodeMap[item]
	return node.index, nil
}

func (queue *UniqueQueue) Pop() (string, error) {
	defer queue.mu.Unlock()
	queue.mu.Lock()
	if queue.size == 0 {
		return "", errors.New("queue is empty")
	}
	node := queue.firstNode
	queue.firstNode = node.next
	if node == queue.lastNode {
		queue.lastNode = nil
	}
	if queue.firstNode != nil {
		queue.firstNode.previous = nil
	}
	delete(queue.nodeMap, node.value)
	decrementNextNodesIndex(node)
	queue.size -= 1
	return node.value, nil
}

func (queue *UniqueQueue) Remove(item string) error {
	defer queue.mu.Unlock()
	queue.mu.Lock()
	if !queue.isIn(item) {
		return errors.New(fmt.Sprintf("item \"%s\" is not in the queue", item))
	}
	node := queue.nodeMap[item]
	if node.previous != nil {
		node.previous.next = node.next
	}
	if node.next != nil {
		node.next.previous = node.previous
	}
	if node == queue.firstNode {
		queue.firstNode = node.next
	}
	if node == queue.lastNode {
		queue.lastNode = node.previous
	}
	delete(queue.nodeMap, node.value)
	decrementNextNodesIndex(node)
	queue.size -= 1
	return nil
}

func (queue *UniqueQueue) IsIn(item string) bool {
	defer queue.mu.Unlock()
	queue.mu.Lock()
	return queue.isIn(item)
}

func (queue *UniqueQueue) Size() uint64 {
	defer queue.mu.Unlock()
	queue.mu.Lock()
	return queue.size
}

func (queue *UniqueQueue) isIn(item string) bool {
	_, ok := queue.nodeMap[item]
	return ok
}

func decrementNextNodesIndex(node *Node) {
	//FIXME: I'm a queue very inneficient
	for node.next != nil {
		node = node.next
		node.index -= 1
	}
}
