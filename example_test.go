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

func Example_iterate() {
	buf := ring.FromSlice([]int{1, 2, 3, 4})
	for i := 0; i < buf.Len(); i++ {
		fmt.Print(buf.Get(i))
	}
	// Output:
	// 1234
}

func Example_offset() {
	buf := ring.FromSlice([]int{10, 20, 30, 40, 50})

	// By default, the offset is 0.
	fmt.Println(buf.Offset())

	// Which means that Head() retrieves the last added item.
	fmt.Println(buf.Head())

	// But we can change the head offset to point to an older item:
	buf.SetOffset(-2)
	fmt.Println(buf.Offset())

	// Now the read head points two slots away from the write head.
	fmt.Println(buf.Head())

	// If we push a new value, it gets pushed at the write head (the actual end
	// of the buffer), pushing over values to the read head.
	fmt.Println(buf.Push(60))
	fmt.Println(buf.Head())

	// Clear the offset:
	buf.SetOffset(0) // or buf.SetOffset(buf.Offset())

	// Output:
	// 0
	// 50 true
	// 2
	// 30 true
	// 10 true
	// 40 true
}

func Example_offsetIterate() {
	buf := ring.FromSlice([]int{1, 2, 3, 4})
	buf.SetOffset(-1) // disregard the current head (4).

	// Len() does not account for offsets.
	for i := 0; i < buf.Len(); i++ {
		fmt.Print(buf.Get(i))
	}
	fmt.Println()

	// So we subtract it.
	for i := 0; i < buf.Len()-buf.Offset(); i++ {
		fmt.Print(buf.Get(i))
	}
	// Output:
	// 1234
	// 123
}
