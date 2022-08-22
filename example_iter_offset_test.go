package ring_test

import (
	"fmt"

	"github.com/xuoe/ring"
)

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
