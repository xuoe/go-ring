// Package ring provides a ring buffer implementation.
package ring

// New creates a ring buffer of the given capacity. It panics if capacity <= 0.
func New[T any](capacity int) *Buffer[T] {
	if capacity <= 0 {
		panic("ring: new buffer capacity <= 0")
	}
	return &Buffer[T]{
		vals: make([]T, capacity),
	}
}

// FromSlice creates a ring buffer with the values of s, allocating a new slice
// to hold them. It panics if len(s) == 0 or if s is nil.
func FromSlice[T any](s []T) *Buffer[T] {
	if len(s) == 0 {
		panic("ring: new buffer from empty/nil slice")
	}
	vals := make([]T, len(s))
	copy(vals, s)
	return &Buffer[T]{
		vals: vals,
		size: len(vals),
	}
}

// Buffer is a ring buffer backed by a slice.
type Buffer[T any] struct {
	vals     []T
	size     int
	idx, off int
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

// Push pushes value v onto the buffer. If the buffer is full, it discards and
// returns the value at the tail end.
func (b *Buffer[T]) Push(v T) (T, bool) {
	ov := b.vals[b.idx]
	b.vals[b.idx] = v
	b.idx = (b.idx + 1) % len(b.vals)
	full := b.size == len(b.vals)
	if !full {
		b.size++
	}
	return ov, full
}

// Pop removes the value at the tail end, if any.
func (b *Buffer[T]) Pop() (T, bool) {
	v, ok := b.Tail()
	if ok {
		b.size--
	}
	return v, ok
}

// Get returns the value at the given index. It panics if the buffer is empty
// or if the index falls outside of buffer range.
func (b *Buffer[T]) Get(idx int) T {
	if idx >= b.size || idx < 0 {
		panic("ring: index out of bounds")
	}
	idx = (b.idx - b.size + idx + len(b.vals)) % len(b.vals)
	return b.vals[idx]
}

// Head returns the value at the offset-adjusted buffer head, if any.
func (b *Buffer[T]) Head() (v T, ok bool) {
	if ok = b.size > 0; !ok {
		return
	}
	v = b.Get(b.Len() - 1 + b.off)
	return
}

// Tail returns the value at the tail end of the buffer, if any,
func (b *Buffer[T]) Tail() (v T, ok bool) {
	if ok = b.size > 0; !ok {
		return
	}
	v = b.Get(0)
	return
}

// SetOffset adjusts the buffer head to point to older values (n < 0) or to
// more recent ones (n > 0). If n is 0, the buffer head is reset to point to
// the latest value. It panics if the offset falls outside buffer range.
//
// Adjusting the offset allows to temporarily keep -n head values from being
// pushed over.
func (b *Buffer[T]) SetOffset(n int) {
	b.off += n
	if n == 0 {
		b.off = 0
	}
	if b.off > 0 || b.off <= -b.size {
		panic("ring: buffer head out of bounds")
	}
}

// Offset returns the buffer head offset.
//
// Calling SetOffset with the result of Offset() is equivalent to SetOffset(0).
func (b *Buffer[T]) Offset() int {
	return -b.off
}

// Len returns the buffer length.
func (b *Buffer[T]) Len() int { return b.size }

// Cap returns the buffer capacity.
func (b *Buffer[T]) Cap() int { return len(b.vals) }

// Full reports whether the buffer is full.
func (b *Buffer[T]) Full() bool { return b.size == len(b.vals) }

// Empty reports whether the buffer is empty.
func (b *Buffer[T]) Empty() bool { return b.size == 0 }
