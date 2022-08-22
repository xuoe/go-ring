package ring_test

import (
	"fmt"

	"github.com/xuoe/ring"
)

func Example_iterate() {
	buf := ring.FromSlice([]int{1, 2, 3, 4})
	for i := 0; i < buf.Len(); i++ {
		fmt.Print(buf.Get(i))
	}
	// Output:
	// 1234
}
