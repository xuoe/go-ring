// Package ring provides a generic ring buffer implementation.
//
// The ring buffer is a FIFO queue with limited capacity. Once the capacity is
// reached, values are discarded from the tail end of the buffer to accommodate
// the ones being pushed at the head, the other end.
//
// Values can be pushed, popped and retrieved (from any index) in O(1) time.
package ring

// New creates an empty ring buffer of the given capacity. It panics if
// capacity <= 0.
func New[T any](capacity int) *Buffer[T] {
	if capacity <= 0 {
		panic("ring.New: capacity <= 0")
	}
	return &Buffer[T]{
		vals: make([]T, capacity),
	}
}

// FromSlice creates a ring buffer out of the given value slice, allocating
// a new slice to hold the values. It panics if cap(vals) == 0 (or if vals ==
// nil).
func FromSlice[T any](vals []T) *Buffer[T] {
	if cap(vals) == 0 {
		panic("ring.FromSlice: empty/nil slice")
	}
	ours := make([]T, cap(vals))
	copy(ours, vals)
	return wrapSlice(ours)
}
	}

func wrapSlice[T any](vals []T) *Buffer[T] {
	return &Buffer[T]{
		vals: vals[:cap(vals)],
		size: len(vals),
	}
}

// Buffer is a ring buffer backed by a slice.
type Buffer[T any] struct {
	vals     []T
	size     int
	idx, off int
	zero     T
}

// ToSlice allocates and returns a slice of the buffer values. If the buffer
// length is 0, the returned slice is nil.
func (b *Buffer[T]) ToSlice() []T {
	if b.size == 0 {
		return nil
	}
	res := make([]T, b.size)
	for i := 0; i < b.size; i++ {
		res[i] = b.Get(i)
	}
	return res
}

// Push pushes value val onto the buffer. If full == true, tail holds the value
// at the tail end of the buffer, which has been discarded to make room for
// val; otherwise, if full == false, tail holds the zero value of T.
func (b *Buffer[T]) Push(val T) (tail T, full bool) {
	tail = b.vals[b.idx]
	b.vals[b.idx] = val
	b.idx = (b.idx + 1) % len(b.vals)
	full = b.size == len(b.vals)
	if !full {
		tail = b.zero
		b.size++
	}
	return tail, full
}

// Pop removes the value at the tail end, if any.
func (b *Buffer[T]) Pop() (T, bool) {
	val, ok := b.Tail()
	if ok {
		b.size--
	}
	return val, ok
}

// Get returns the value at the given index idx. It panics if the buffer is
// empty or if the index falls outside buffer range.
func (b *Buffer[T]) Get(idx int) T {
	if idx >= b.size || idx < 0 {
		panic("ring.Get: index out of bounds")
	}
	idx = (b.idx - b.size + idx + len(b.vals)) % len(b.vals)
	return b.vals[idx]
}

// Head returns the value at the offset-adjusted buffer head, if any.
func (b *Buffer[T]) Head() (T, bool) {
	if b.size > 0 {
		return b.Get(b.Len() - 1 + b.off), true
	}
	return b.zero, false
}

// Tail returns the value at the tail end of the buffer, if any.
func (b *Buffer[T]) Tail() (T, bool) {
	if b.size > 0 {
		return b.Get(0), true
	}
	return b.zero, false
}

// SetOffset adjusts the buffer head to point to older values (n < 0),
// effectively ignoring the last n values, or to more recent ones (n > 0), if
// already adjusted. If n == 0, the buffer head is reset to point to the last
// pushed value. It panics if the offset falls outside buffer range.
//
// Note that SetOffset(Offset()) is equivalent to SetOffset(0), regardless of
// the result of Offset().
func (b *Buffer[T]) SetOffset(n int) {
	b.off += n
	if n == 0 {
		b.off = 0
	}
	if b.off > 0 || b.off <= -b.size {
		panic("ring.SetOffset: buffer head out of bounds")
	}
}

// Offset returns the buffer head offset, where 0 indicates no offset.
func (b *Buffer[T]) Offset() int { return -b.off }

// Len returns the buffer length.
func (b *Buffer[T]) Len() int { return b.size }

// Cap returns the buffer capacity, i.e., the maximum number of values the
// buffer can hold without discarding older ones.
func (b *Buffer[T]) Cap() int { return len(b.vals) }

// Full reports whether the buffer is full.
func (b *Buffer[T]) Full() bool { return b.size == len(b.vals) }

// Empty reports whether the buffer is empty.
func (b *Buffer[T]) Empty() bool { return b.size == 0 }
