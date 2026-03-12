package reveal

import (
	"strings"
	"testing"
)

// compat_test.go — golden file compatibility tests comparing reveal output
// against the expected go-spew v1.1.1 format for key types.
//
// These tests verify that reveal's Sdump output matches the structural format
// that go-spew v1.1.1 produces (type annotations, indentation, pointer
// handling, etc.) so that reveal can serve as a drop-in replacement.

// spewLike returns a ConfigState that mimics go-spew's default config.
func spewLike() ConfigState {
	return ConfigState{
		Indent:                  " ",
		DisablePointerAddresses: true, // addresses are non-deterministic
		SortKeys:                true,
	}
}

func TestCompat_Int(t *testing.T) {
	cs := spewLike()
	got := cs.Sdump(42)
	expect := "(int) 42\n"
	if got != expect {
		t.Errorf("int:\ngot:    %q\nexpect: %q", got, expect)
	}
}

func TestCompat_String(t *testing.T) {
	cs := spewLike()
	got := cs.Sdump("hello")
	if !strings.Contains(got, `(string)`) {
		t.Errorf("string: missing type annotation, got: %q", got)
	}
	if !strings.Contains(got, `"hello"`) {
		t.Errorf("string: missing quoted value, got: %q", got)
	}
	if !strings.Contains(got, `len=5`) {
		t.Errorf("string: missing length, got: %q", got)
	}
}

func TestCompat_Bool(t *testing.T) {
	cs := spewLike()
	got := cs.Sdump(true)
	expect := "(bool) true\n"
	if got != expect {
		t.Errorf("bool:\ngot:    %q\nexpect: %q", got, expect)
	}
}

func TestCompat_Float(t *testing.T) {
	cs := spewLike()
	got := cs.Sdump(3.14)
	if !strings.Contains(got, "(float64)") {
		t.Errorf("float: missing type, got: %q", got)
	}
	if !strings.Contains(got, "3.14") {
		t.Errorf("float: missing value, got: %q", got)
	}
}

func TestCompat_NilPointer(t *testing.T) {
	cs := spewLike()
	var p *int
	got := cs.Sdump(p)
	if !strings.Contains(got, "(*int)") {
		t.Errorf("nil ptr: missing type, got: %q", got)
	}
	if !strings.Contains(got, "<nil>") {
		t.Errorf("nil ptr: missing <nil>, got: %q", got)
	}
}

func TestCompat_Pointer(t *testing.T) {
	cs := spewLike()
	x := 42
	got := cs.Sdump(&x)
	if !strings.Contains(got, "(*int)") {
		t.Errorf("ptr: missing type, got: %q", got)
	}
	if !strings.Contains(got, "(int) 42") {
		t.Errorf("ptr: missing pointed value, got: %q", got)
	}
}

func TestCompat_Slice(t *testing.T) {
	cs := spewLike()
	got := cs.Sdump([]int{1, 2, 3})
	if !strings.Contains(got, "[]int") {
		t.Errorf("slice: missing type, got: %q", got)
	}
	if !strings.Contains(got, "len=3") {
		t.Errorf("slice: missing len, got: %q", got)
	}
	if !strings.Contains(got, "cap=3") {
		t.Errorf("slice: missing cap, got: %q", got)
	}
	if !strings.Contains(got, "(int) 1") {
		t.Errorf("slice: missing element, got: %q", got)
	}
}

func TestCompat_Map(t *testing.T) {
	cs := spewLike()
	got := cs.Sdump(map[string]int{"a": 1, "b": 2})
	if !strings.Contains(got, "map[string]int") {
		t.Errorf("map: missing type, got: %q", got)
	}
	if !strings.Contains(got, "len=2") {
		t.Errorf("map: missing len, got: %q", got)
	}
	// With SortKeys, "a" should appear before "b".
	idxA := strings.Index(got, "(a)")
	idxB := strings.Index(got, "(b)")
	if idxA == -1 || idxB == -1 || idxA >= idxB {
		t.Errorf("map: keys not sorted, got: %q", got)
	}
}

func TestCompat_Struct(t *testing.T) {
	type Point struct {
		X int
		Y int
	}
	cs := spewLike()
	got := cs.Sdump(Point{X: 10, Y: 20})
	if !strings.Contains(got, "reveal.Point") {
		t.Errorf("struct: missing type, got: %q", got)
	}
	if !strings.Contains(got, "X:") {
		t.Errorf("struct: missing field name X, got: %q", got)
	}
	if !strings.Contains(got, "(int) 10") {
		t.Errorf("struct: missing field value 10, got: %q", got)
	}
	if !strings.Contains(got, "Y:") {
		t.Errorf("struct: missing field name Y, got: %q", got)
	}
}

func TestCompat_NilSlice(t *testing.T) {
	cs := spewLike()
	var s []int
	got := cs.Sdump(s)
	if !strings.Contains(got, "[]int") {
		t.Errorf("nil slice: missing type, got: %q", got)
	}
	if !strings.Contains(got, "<nil>") {
		t.Errorf("nil slice: missing <nil>, got: %q", got)
	}
}

func TestCompat_Nil(t *testing.T) {
	cs := spewLike()
	got := cs.Sdump(nil)
	if !strings.Contains(got, "<nil>") {
		t.Errorf("nil: missing <nil>, got: %q", got)
	}
}

func TestCompat_CircularRef(t *testing.T) {
	type Node struct {
		V    int
		Next *Node
	}
	cs := spewLike()
	a := &Node{V: 1}
	b := &Node{V: 2}
	a.Next = b
	b.Next = a
	got := cs.Sdump(a)
	if !strings.Contains(got, "<already shown>") {
		t.Errorf("circular: missing detection, got: %q", got)
	}
}
