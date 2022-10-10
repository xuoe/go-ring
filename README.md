[![Go Reference](https://pkg.go.dev/badge/github.com/xuoe/go-ring.svg)](https://pkg.go.dev/github.com/xuoe/go-ring)

This library provides a generic [ring buffer](https://en.wikipedia.org/wiki/Circular_buffer) implementation.

### API

```go
func New(int) *Buffer[T]
func FromSlice([]T) *Buffer[T]
func WrapSlice([]T) *Buffer[T]

func (*Buffer[T]) Push(T) (T, bool)
func (*Buffer[T]) Pop() (T, bool)
func (*Buffer[T]) Head() (T, bool)
func (*Buffer[T]) Tail() (T, bool)
func (*Buffer[T]) Get(int) T
func (*Buffer[T]) Cap() int
func (*Buffer[T]) Len() int
func (*Buffer[T]) Empty() bool
func (*Buffer[T]) Full() bool
func (*Buffer[T]) Offset() int
func (*Buffer[T]) SetOffset(int)
func (*Buffer[T]) ToSlice() []T
```

See [the docs](https://pkg.go.dev/github.com/xuoe/go-ring) for more.
