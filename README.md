[![Go Reference](https://pkg.go.dev/badge/github.com/xuoe/ring.svg)](https://pkg.go.dev/github.com/xuoe/ring)

This library provides a generic ring buffer implementation.

A ring buffer offers a window over a stream of values. New values are pushed at 
one end (the head), while older ones are discarded from the other end (the 
tail) to accommodate the new ones.

#### API Overview

```go
func New[T any](size int) *Buffer[T]
func FromSlice[T any]([]T) *Buffer[T]
func (*Buffer[T]) Head() (T, bool)
func (*Buffer[T]) Tail() (T, bool)
func (*Buffer[T]) Push(T) (T, bool)
func (*Buffer[T]) Offset() int
func (*Buffer[T]) SetOffset(int)
func (*Buffer[T]) Empty() bool
func (*Buffer[T]) Full() bool
func (*Buffer[T]) Len() int
```
