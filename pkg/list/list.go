// Copyright (c) 2024 Gophers. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

// Package list implements support for a generic ordered List.
// A List is a Collection that wraps an underlying doubly linked list
// and provides convenience methods and synthatic sugar on top of it.
//
// Compared to a Sequence, a List allows for efficient insertion and removal of elements
// at the edges of the list, but slower access to arbitrary elements. This makes the List a good choice
// for implementing queues or stacks.
//
// For a list of comparable types, consider using ComparableList,
// which provides additional methods for comparable types.
package list

import (
	"fmt"
	"iter"
	"math/rand"
	"slices"

	"github.com/charbz/gophers/pkg/collection"
)

type Node[T any] struct {
	value T
	next  *Node[T]
	prev  *Node[T]
}

type List[T any] struct {
	head *Node[T]
	tail *Node[T]
	size int
}

func NewList[T any](s ...[]T) *List[T] {
	list := new(List[T])
	if len(s) == 0 {
		return list
	}
	for _, slice := range s {
		for _, v := range slice {
			list.Add(v)
		}
	}
	return list
}

// The following methods implement
// the Collection interface.

// Add adds a value to the end of the list.
func (l *List[T]) Add(v T) {
	node := &Node[T]{value: v}
	if l.head == nil {
		l.head = node
		l.tail = node
	} else {
		l.tail.next = node
		node.prev = l.tail
		l.tail = node
	}
	l.size++
}

// Length returns the number of nodes in the list.
func (l *List[T]) Length() int {
	return l.size
}

// New returns a new list.
func (l *List[T]) New(s ...[]T) collection.Collection[T] {
	return NewList(s...)
}

// Random returns a random value from the list.
func (l *List[T]) Random() T {
	if l.size == 0 {
		return *new(T)
	}
	return l.At(rand.Intn(l.size))
}

// Values returns an iterator for all values in the list.
func (l *List[T]) Values() iter.Seq[T] {
	return func(yield func(T) bool) {
		for node := l.head; node != nil; node = node.next {
			if !yield(node.value) {
				break
			}
		}
	}
}

// The following methods implement
// the OrderedCollection interface.

// At returns the value of the node at the given index.
func (l *List[T]) At(index int) T {
	if index < 0 || index >= l.size {
		panic(collection.IndexOutOfBoundsError)
	}
	node := l.head
	for i := 0; i < index; i++ {
		node = node.next
	}
	return node.value
}

// All returns an index/value iterator for all nodes in the list.
func (l *List[T]) All() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		i := 0
		for node := l.head; node != nil; node = node.next {
			if !yield(i, node.value) {
				break
			}
			i++
		}
	}
}

// Backward returns an index/value iterator for all nodes in the list in reverse order.
func (l *List[T]) Backward() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		i := l.size - 1
		for node := l.tail; node != nil; node = node.prev {
			if !yield(i, node.value) {
				break
			}
			i--
		}
	}
}

// Slice returns a new list containing only the nodes between the start and end indices.
func (l *List[T]) Slice(start, end int) collection.OrderedCollection[T] {
	if start < 0 || end > l.size || start > end {
		panic(collection.IndexOutOfBoundsError)
	}
	list := &List[T]{}
	for i, v := range l.All() {
		if i < start {
			continue
		}
		if i >= start && i < end {
			list.Add(v)
		}
		if i >= end {
			break
		}
	}
	return list
}

// NewOrdered returns a new ordered collection.
func (l *List[T]) NewOrdered(s ...[]T) collection.OrderedCollection[T] {
	return NewList(s...)
}

// ToSlice returns a slice containing all values in the list.
func (l *List[T]) ToSlice() []T {
	slice := make([]T, 0, l.size)
	for v := range l.Values() {
		slice = append(slice, v)
	}
	return slice
}

// Implement the Stringer interface.
func (l *List[T]) String() string {
	return fmt.Sprintf("List(%T) %v", *new(T), l.ToSlice())
}

// Clone returns a copy of the list. This is a shallow clone.
func (l *List[T]) Clone() *List[T] {
	clone := &List[T]{}
	for v := range l.Values() {
		clone.Add(v)
	}
	return clone
}

// Count is an alias for collection.Count
func (l *List[T]) Count(f func(T) bool) int {
	return collection.Count(l, f)
}

// Concat returns a new list concatenating the passed in lists.
func (l *List[T]) Concat(lists ...*List[T]) *List[T] {
	clone := l.Clone()
	for _, list := range lists {
		for v := range list.Values() {
			clone.Add(v)
		}
	}
	return clone
}

// Contains tests whether a predicate holds for at least one element of this list.
func (l *List[T]) Contains(f func(T) bool) bool {
	i, _ := collection.Find(l, f)
	return i > -1
}

// Corresponds is an alias for collection.Corresponds
func (l *List[T]) Corresponds(s *List[T], f func(T, T) bool) bool {
	return collection.Corresponds(l, s, f)
}

// Dequeue removes and returns the first element of the list.
func (l *List[T]) Dequeue() (T, error) {
	if l.size == 0 {
		return *new(T), collection.EmptyCollectionError
	}
	element := l.head.value
	l.head = l.head.next
	l.size--
	return element, nil
}

// Distinct takes an "equality" function as an argument
// and returns a new sequence containing all the unique elements
// If you prefer not to pass an equality function use a ComparableList.
func (l *List[T]) Distinct(f func(T, T) bool) *List[T] {
	return NewList(slices.CompactFunc(l.ToSlice(), f))
}

// Drop is an alias for collection.Drop
func (l *List[T]) Drop(n int) *List[T] {
	return collection.Drop(l, n).(*List[T])
}

// DropWhile is an alias for collection.DropWhile
func (l *List[T]) DropWhile(f func(T) bool) *List[T] {
	return collection.DropWhile(l, f).(*List[T])
}

// DropRight is an alias for collection.DropRight
func (l *List[T]) DropRight(n int) *List[T] {
	return collection.DropRight(l, n).(*List[T])
}

// Enqueue appends an element to the list.
func (l *List[T]) Enqueue(v T) {
	l.Add(v)
}

// Equals takes a list and an equality function as arguments
// and returns true if the two sequences are equal.
// If you prefer not to pass an equality function use a ComparableList.
func (l *List[T]) Equals(s *List[T], f func(T, T) bool) bool {
	if l.size != s.size {
		return false
	}
	n1 := l.head
	n2 := s.head
	for n1 != nil && n2 != nil {
		if !f(n1.value, n2.value) {
			return false
		}
		n1 = n1.next
		n2 = n2.next
	}
	return true
}

// Exists is an alias for Contains
func (l *List[T]) Exists(f func(T) bool) bool {
	return l.Contains(f)
}

// Filter is an alias for collection.Filter
func (l *List[T]) Filter(f func(T) bool) *List[T] {
	return collection.Filter(l, f).(*List[T])
}

// FilterNot is an alias for collection.FilterNot
func (l *List[T]) FilterNot(f func(T) bool) *List[T] {
	return collection.FilterNot(l, f).(*List[T])
}

// Find is an alias for collection.Find
func (l *List[T]) Find(f func(T) bool) (int, T) {
	return collection.Find(l, f)
}

// FindLast is an alias for collection.FindLast
func (l *List[T]) FindLast(f func(T) bool) (int, T) {
	return collection.FindLast(l, f)
}

// ForEach is an alias for collection.ForEach
func (l *List[T]) ForEach(f func(T)) *List[T] {
	return collection.ForEach(l, f).(*List[T])
}

// ForAll is an alias for collection.ForAll
func (l *List[T]) ForAll(f func(T) bool) bool {
	return collection.ForAll(l, f)
}

// Head is an alias for collection.Head
func (l *List[T]) Head() (T, error) {
	return collection.Head(l)
}

// Init is an alias for collection.Init
func (l *List[T]) Init() *List[T] {
	return collection.Init(l).(*List[T])
}

// IsEmpty returns true if the list is empty.
func (l *List[T]) IsEmpty() bool {
	return l.size == 0
}

// Last is an alias for collection.Last
func (l *List[T]) Last() (T, error) {
	return collection.Last(l)
}

// NonEmpty returns true if the list is not empty.
func (l *List[T]) NonEmpty() bool {
	return l.size > 0
}

// Pop removes and returns the last element of the list.
func (l *List[T]) Pop() (T, error) {
	if l.size == 0 {
		return *new(T), collection.EmptyCollectionError
	}
	element := l.tail.value
	l.tail = l.tail.prev
	l.size--
	return element, nil
}

// Push appends an element to the list.
func (l *List[T]) Push(v T) {
	l.Add(v)
}

// Partition is an alias for collection.Partition
func (l *List[T]) Partition(f func(T) bool) (*List[T], *List[T]) {
	left, right := collection.Partition(l, f)
	return left.(*List[T]), right.(*List[T])
}

// Reverse is an alias for collection.Reverse
func (l *List[T]) Reverse() *List[T] {
	return collection.Reverse(l).(*List[T])
}

// SplitAt is an alias for collection.SplitAt
func (l *List[T]) SplitAt(n int) (*List[T], *List[T]) {
	left, right := collection.SplitAt(l, n)
	return left.(*List[T]), right.(*List[T])
}

// Take is an alias for collection.Take
func (l *List[T]) Take(n int) *List[T] {
	return collection.Take(l, n).(*List[T])
}

// TakeRight is an alias for collection.TakeRight
func (l *List[T]) TakeRight(n int) *List[T] {
	return collection.TakeRight(l, n).(*List[T])
}

// Tail is an alias for collection.Tail
func (l *List[T]) Tail() *List[T] {
	return collection.Tail(l).(*List[T])
}