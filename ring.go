package ring

// New creates a ring buffer of the given size. It panics if size < 0.
func New[T any](size int) *Buffer[T] {
	if size < 0 {
		panic("ring: new buffer size < 0")
	}
	return &Buffer[T]{
		vals: make([]T, size),
	}
}

// FromSlice creates a ring buffer with the values of s, allocating a new slice
// to hold them.
func FromSlice[T any](s []T) *Buffer[T] {
	vals := make([]T, len(s))
	copy(vals, s)
	return &Buffer[T]{
		vals: vals,
		size: len(vals),
	}
}

// Buffer is a ring buffer.
type Buffer[T any] struct {
	vals     []T
	size     int
	idx, off int
}

// Push pushes value v onto the buffer, discarding and returning the value
// at the tail end of the buffer, if any.
func (b *Buffer[T]) Push(v T) (T, bool) {
	ov := b.vals[b.idx]
	b.vals[b.idx] = v
	b.idx = (b.idx + 1) % len(b.vals)
	full := b.size == len(b.vals)
	if b.size < len(b.vals) {
		b.size++
	}
	return ov, full
}

// Head returns the value at the offset-adjusted buffer head, if any.
func (b *Buffer[T]) Head() (T, bool) {
	idx := (b.idx + b.off - 1 + len(b.vals)) % len(b.vals)
	return b.vals[idx], b.size > 0
}

// Tail returns the value at the tail end of the buffer, if any.
func (b *Buffer[T]) Tail() (T, bool) {
	return b.vals[b.idx], b.size == len(b.vals)
}

// SetOffset adjusts the buffer head to point to older values (n < 0) or to
// more recent ones (n > 0). If n is 0, the buffer head is reset to point to
// the latest value. It panics if the offset falls outside buffer range.
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

// Full reports whether the buffer is full.
func (b *Buffer[T]) Full() bool { return b.size == len(b.vals) }

// Empty reports whether the buffer is empty.
func (b *Buffer[T]) Empty() bool { return b.size == 0 }
