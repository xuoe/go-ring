package ring_test

import (
	"fmt"

	"github.com/xuoe/go-ring"
)

func Example_reverseIterate() {
	buf := ring.FromSlice([]int{1, 2, 3, 4})
	for i := buf.Len() - 1; i >= 0; i-- {
		fmt.Print(buf.Get(i))
	}
	// Output:
	// 4321
}
