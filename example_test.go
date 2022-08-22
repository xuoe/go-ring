package ring_test

import (
	"fmt"

	"github.com/xuoe/go-ring"
)

func Example() {
	buf := ring.FromSlice([]string{"foo", "bar", "baz"})

	// Get the head.
	fmt.Println(buf.Head())

	// Pushing another value drops the first one and returns it.
	fmt.Println(buf.Push("qux"))
	fmt.Println(buf.Head())
	fmt.Println(buf.Tail())

	// Popping the current tail will have the buffer tail point to the next
	// element.
	fmt.Println(buf.Pop())
	fmt.Println(buf.Tail())

	// Check len/cap.
	fmt.Println(buf.Len(), buf.Cap())

	// Output
	// "baz" true
	// "foo" true
	// "qux" true
	// "bar" true
	// "baz" true
	// 2 3
}
