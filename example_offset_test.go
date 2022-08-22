package ring_test

import (
	"fmt"

	"github.com/xuoe/go-ring"
)

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
