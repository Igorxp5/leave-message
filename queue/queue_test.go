package queue

import (
	"fmt"
	"testing"
)

func TestZeroQueue(t *testing.T) {
	queue := New(0)
	_, err := queue.Add("something")
	if err == nil {
		t.Fatalf("should not allow add items in zero-sized queue")
	}
}

func TestAddIndex(t *testing.T) {
	var index uint64
	var err error

	queue := New(2)
	index, err = queue.Add("something")
	if err != nil {
		t.Fatalf("%v", err)
	}
	if index != 0 {
		t.Fatalf("the first element added should have index 0")
	}
	index, err = queue.Add("other")
	if err != nil {
		t.Fatalf("%v", err)
	}
	if index != 1 {
		t.Fatalf("the second element added should have index 1")
	}
}

func TestMaxReached(t *testing.T) {
	queue := New(3)
	for i := 0; i < 3; i++ {
		_, err := queue.Add(fmt.Sprintf("%d", i))
		if err != nil {
			t.Fatalf("should allow add items while the queue is not full")
		}
	}
	_, err := queue.Add("something")
	if err == nil {
		t.Fatalf("should not allow add items in a full queue")
	}
}

func TestGetItems(t *testing.T) {
	var value interface{}
	var err error

	queue := New(3)
	queue.Add("Get")

	if queue.Size() != 1 {
		t.Fatalf("queue size should be 1")
	}

	queue.Add("Post")

	if queue.Size() != 2 {
		t.Fatalf("queue size should be 2")
	}

	queue.Add("Update")

	if queue.Size() != 3 {
		t.Fatalf("queue size should be 3")
	}

	value, err = queue.Pop()
	if err != nil {
		t.Fatalf("1: not able to pop item from the queue")
	}
	if value != "Get" {
		t.Fatalf("1: queue shuold be FIFO")
	}

	value, err = queue.Pop()
	if err != nil {
		t.Fatalf("2: not able to pop item from the queue")
	}
	if value != "Post" {
		t.Fatalf("2: queue shuold be FIFO")
	}

	value, err = queue.Pop()
	if err != nil {
		t.Fatalf("3: not able to pop item from the queue")
	}
	if value != "Update" {
		t.Fatalf("3: queue shuold be FIFO")
	}
	value, err = queue.Pop()
	if err == nil {
		t.Fatalf("should return error when popping a empty queue")
	}

	if queue.Size() != 0 {
		t.Fatalf("queue size should be 0")
	}
}

func TestAddSameItem(t *testing.T) {
	var err error
	var value string
	queue := New(2)
	_, err = queue.Add("unique")
	if err != nil {
		t.Fatalf("%v", err)
	}
	_, err = queue.Add("unique")
	if err == nil {
		t.Fatalf("no error should happen when adding a item")
	}
	value, err = queue.Pop()
	if err != nil {
		t.Fatalf("%v", err)
	}
	if value != "unique" {
		t.Fatalf("expecting \"unique\" as value from popping the queue")
	}
	value, err = queue.Pop()
	if err == nil {
		t.Fatalf("should return error when popping a empty queue")
	}
}

func TestRemoveItem(t *testing.T) {
	var err error
	var value string

	queue := New(3)
	queue.Add("one")
	if queue.Size() != 1 {
		t.Fatalf("queue size should be 1")
	}
	queue.Add("two")
	if queue.Size() != 2 {
		t.Fatalf("queue size should be 2")
	}
	err = queue.Remove("three")
	if err == nil {
		t.Fatalf("queue shuold return an error removing a non-added item")
	}
	err = queue.Remove("one")
	if err != nil {
		t.Fatalf("%v", err)
	}
	value, err = queue.Pop()
	if err != nil {
		t.Fatalf("%v", err)
	}
	if value != "two" {
		t.Fatalf("expecting \"two\" as value from popping the queue")
	}
	value, err = queue.Pop()
	if err == nil {
		t.Fatalf("should return error when popping a empty queue")
	}
}

func TestIndex(t *testing.T) {
	queue := New(3)
	queue.Add("foo")
	queue.Add("bar")
	queue.Add("path")

	var index uint64
	var err error

	index, err = queue.Index("bar")
	if err != nil {
		t.Fatalf("%v", err)
	}
	if index != 1 {
		t.Fatalf("\"bar\" index should be 1")
	}

	_, err = queue.Pop()
	if err != nil {
		t.Fatalf("%v", err)
	}

	index, err = queue.Index("bar")
	if err != nil {
		t.Fatalf("%v", err)
	}
	if index != 0 {
		t.Fatalf("\"bar\" index should be 1")
	}

	queue.Add("other")
	index, err = queue.Index("other")
	if err != nil {
		t.Fatalf("%v", err)
	}
	if index != 2 {
		t.Fatalf("\"other\" index should be 2")
	}

	err = queue.Remove("bar")
	if err != nil {
		t.Fatalf("%v", err)
	}

	index, err = queue.Index("other")
	if err != nil {
		t.Fatalf("%v", err)
	}
	if index != 1 {
		t.Fatalf("\"other\" index should be 1")
	}
}
