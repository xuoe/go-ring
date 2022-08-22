package ring_test

import (
	"strconv"
	"testing"

	"github.com/xuoe/go-ring"
)

func TestPush(t *testing.T) {
	for i, test := range []struct {
		size int
		newv []int
		oldv []int
		ok   []bool
	}{
		{
			5,
			[]int{3, 2, 1},
			[]int{0, 0, 0},
			[]bool{false, false, false},
		},
		{
			5,
			[]int{3, 2, 1, 7, 2, 1, 10},
			[]int{0, 0, 0, 0, 0, 3, 2},
			[]bool{false, false, false, false, false, true, true},
		},
		{
			2,
			[]int{3, 2, 1, 7, 2},
			[]int{0, 0, 3, 2, 1},
			[]bool{false, false, true, true, true},
		},
	} {
		t.Run(strconv.Itoa(i+1), func(t *testing.T) {
			b := ring.New[int](test.size)
			if exp, got := 0, b.Len(); exp != got {
				t.Fatalf("new Len(): %d -%d", got, exp)
			}
			for i := range test.newv {
				t.Run(strconv.Itoa(i+1), func(t *testing.T) {
					old, ok := b.Push(test.newv[i])
					if exp, got := test.oldv[i], old; exp != got {
						t.Errorf("old value: %d -%d", got, exp)
					}
					if exp, got := test.ok[i], ok; exp != got {
						t.Errorf("had old value: %t -%t", got, exp)
					}
				})
			}
		})
	}
}

func TestPop(t *testing.T) {
	for i, test := range []struct {
		size int
		push []int
		pop  []int
	}{
		{
			5,
			[]int{3, 2, 1, 10, 2},
			[]int{3, 2, 1},
		},
		{
			3,
			[]int{3, 2, 1},
			[]int{3, 2, 1, 0, 0, 0},
		},
		{
			5,
			[]int{3, 2, 1},
			[]int{3, 2, 1, 0, 0, 0},
		},
	} {
		t.Run(strconv.Itoa(i+1), func(t *testing.T) {
			b := ring.New[int](test.size)
			for _, v := range test.push {
				b.Push(v)
			}
			for i := 0; i < len(test.pop); i++ {
				v, ok := b.Pop()
				if exp, got := test.pop[i], v; exp != got {
					t.Errorf("%d: Pop() value: -%d +%d", i, exp, got)
				}
				if i >= test.size {
					if exp, got := false, ok; exp != got {
						t.Errorf("%d: Pop() ok: -%t +%t", i, exp, got)
					}
				}
			}
		})
	}
}

func TestHeadTail(t *testing.T) {
	b := ring.New[int](3)
	for i, test := range []struct {
		do   func()
		head int
		tail int
		size int
	}{
		{
			do:   func() { b.Push(1) },
			head: 1,
			tail: 1,
			size: 1,
		},
		{
			do:   func() { b.Push(2) },
			head: 2,
			tail: 1,
			size: 2,
		},
		{
			do:   func() { b.Push(3) },
			head: 3,
			tail: 1,
			size: 3,
		},
		{
			do:   func() { b.Push(4) },
			head: 4,
			tail: 2,
			size: 3,
		},
		{
			do:   func() { b.Push(5) },
			head: 5,
			tail: 3,
			size: 3,
		},
		{
			do:   func() { b.Pop() },
			head: 5,
			tail: 4,
			size: 2,
		},
		{
			do:   func() { b.Pop() },
			head: 5,
			tail: 5,
			size: 1,
		},
		{
			do:   func() { b.Pop() },
			head: 0,
			tail: 0,
			size: 0,
		},
	} {
		t.Run(strconv.Itoa(i+1), func(t *testing.T) {
			test.do()
			{
				exp := test.head
				got, _ := b.Head()
				if exp != got {
					t.Errorf("Head(): -%d +%d", exp, got)
				}
			}
			{
				exp := test.tail
				got, _ := b.Tail()
				if exp != got {
					t.Errorf("Tail(): -%d +%d", exp, got)
				}
			}
			{
				exp := test.size
				got := b.Len()
				if exp != got {
					t.Errorf("Len(): -%d +%d", exp, got)
				}
			}
		})
	}
}

func TestOffset(t *testing.T) {
	type value struct {
		int
		bool
	}
	b := ring.New[int](4)
	for i, test := range []struct {
		op         func()
		head, tail value
		off        int
	}{
		{
			func() {},
			value{0, false},
			value{0, false},
			0,
		},
		{
			func() { b.Push(1) },
			value{1, true},
			value{1, true},
			0,
		},
		{
			func() { b.Push(2) },
			value{2, true},
			value{1, true},
			0,
		},
		{
			func() { b.Push(3) },
			value{3, true},
			value{1, true},
			0,
		},
		{
			func() { b.Push(4) },
			value{4, true},
			value{1, true},
			0,
		},
		{
			func() { b.SetOffset(-1) },
			value{3, true},
			value{1, true},
			1,
		},
		{
			func() { b.SetOffset(-1) },
			value{2, true},
			value{1, true},
			2,
		},
		{
			func() { b.SetOffset(0) },
			value{4, true},
			value{1, true},
			0,
		},
		{
			func() { b.SetOffset(-1); b.Pop() },
			value{3, true},
			value{2, true},
			1,
		},
	} {
		t.Run(strconv.Itoa(i+1), func(t *testing.T) {
			test.op()
			{
				exp := test.head
				v, ok := b.Head()
				got := value{v, ok}
				if exp != got {
					t.Errorf("head: %v -%v", got, exp)
				}
			}
			{
				exp := test.tail
				v, ok := b.Tail()
				got := value{v, ok}
				if exp != got {
					t.Errorf("tail: %v -%v", got, exp)
				}
			}
			{
				exp := test.off
				got := b.Offset()
				if exp != got {
					t.Errorf("offset: %v -%v", got, exp)
				}
			}
		})
	}
}

func TestFullEmptyLen(t *testing.T) {
	b := ring.New[int](5)
	for i, test := range []struct {
		op    func()
		len   int
		full  bool
		empty bool
	}{
		{
			func() {},
			0, false, true,
		},
		{
			func() { b.Push(10) },
			1, false, false,
		},
		{
			func() { b.Push(11) },
			2, false, false,
		},
		{
			func() {
				b.Push(11)
				b.Push(11)
				b.Push(11)
			},
			5, true, false,
		},
		{
			func() {
				b.Push(11)
				b.Push(11)
				b.Push(11)
			},
			5, true, false,
		},
	} {
		t.Run(strconv.Itoa(i+1), func(t *testing.T) {
			test.op()
			if exp, got := test.len, b.Len(); exp != got {
				t.Errorf("len: %v -%v", got, exp)
			}
			if exp, got := test.empty, b.Empty(); exp != got {
				t.Errorf("empty: %v -%v", got, exp)
			}
			if exp, got := test.full, b.Full(); exp != got {
				t.Errorf("full: %v -%v", got, exp)
			}
		})
	}
}

func TestFromSlice(t *testing.T) {
	vals := []int{1, 2, 3, 4, 5}
	b := ring.FromSlice(vals)
	{
		exp := 5
		got, _ := b.Head()
		if exp != got {
			t.Errorf("head: %v -%v", got, exp)
		}
	}
	{
		exp := 1
		got, _ := b.Tail()
		if exp != got {
			t.Errorf("tail: %v -%v", got, exp)
		}
	}
	if exp, got := 5, b.Len(); exp != got {
		t.Errorf("len: %v -%v", got, exp)
	}
	if exp, got := true, b.Full(); exp != got {
		t.Errorf("full: %v -%v", got, exp)
	}
	if exp, got := false, b.Empty(); exp != got {
		t.Errorf("empty: %v -%v", got, exp)
	}
	b.Push(6)
	{
		exp := 6
		got, _ := b.Head()
		if exp != got {
			t.Errorf("head: %v -%v", got, exp)
		}
	}
	{
		exp := 2
		got, _ := b.Tail()
		if exp != got {
			t.Errorf("tail: %v -%v", got, exp)
		}
	}
}
