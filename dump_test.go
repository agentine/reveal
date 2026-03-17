package reveal

import (
	"strings"
	"testing"
	"unsafe"
)

func TestDumpUnexportedFields(t *testing.T) {
	type secret struct {
		Exported   int
		unexported string
	}
	s := secret{Exported: 42, unexported: "hidden"}
	result := Sdump(s)

	if !strings.Contains(result, "Exported:") {
		t.Errorf("expected Exported field in output, got: %s", result)
	}
	if !strings.Contains(result, "(int) 42") {
		t.Errorf("expected (int) 42 in output, got: %s", result)
	}
	if unsafeEnabled {
		if !strings.Contains(result, "unexported:") {
			t.Errorf("expected unexported field name in output, got: %s", result)
		}
	}
}

func TestDumpNestedStruct(t *testing.T) {
	type Inner struct {
		Value int
	}
	type Outer struct {
		Name  string
		Inner Inner
	}
	o := Outer{Name: "test", Inner: Inner{Value: 99}}
	result := Sdump(o)
	if !strings.Contains(result, "Outer") {
		t.Errorf("expected Outer in output, got: %s", result)
	}
	if !strings.Contains(result, "Inner") {
		t.Errorf("expected Inner in output, got: %s", result)
	}
	if !strings.Contains(result, "(int) 99") {
		t.Errorf("expected (int) 99 in output, got: %s", result)
	}
}

func TestDumpSliceOfStructs(t *testing.T) {
	type Item struct {
		ID   int
		Name string
	}
	items := []Item{{1, "foo"}, {2, "bar"}}
	result := Sdump(items)
	if !strings.Contains(result, "[]reveal.Item") {
		t.Errorf("expected []reveal.Item in output, got: %s", result)
	}
	if !strings.Contains(result, `"foo"`) {
		t.Errorf("expected foo in output, got: %s", result)
	}
	if !strings.Contains(result, `"bar"`) {
		t.Errorf("expected bar in output, got: %s", result)
	}
}

func TestDumpMapOfStructs(t *testing.T) {
	type Entry struct {
		Value int
	}
	cs := ConfigState{Indent: " ", SortKeys: true}
	m := map[string]Entry{"a": {1}, "b": {2}}
	result := cs.Sdump(m)
	if !strings.Contains(result, "map[string]reveal.Entry") {
		t.Errorf("expected map type in output, got: %s", result)
	}
}

func TestDumpUintptr(t *testing.T) {
	v := uintptr(0xDEADBEEF)
	result := Sdump(v)
	if !strings.Contains(result, "(uintptr)") {
		t.Errorf("expected (uintptr) in output, got: %s", result)
	}
}

func TestDumpUnsafePointer(t *testing.T) {
	x := 42
	p := unsafe.Pointer(&x)
	result := Sdump(p)
	if !strings.Contains(result, "unsafe.Pointer") {
		t.Errorf("expected unsafe.Pointer in output, got: %s", result)
	}
}

func TestDumpEmptySlice(t *testing.T) {
	s := []int{}
	result := Sdump(s)
	if !strings.Contains(result, "len=0") {
		t.Errorf("expected len=0 in output, got: %s", result)
	}
}

func TestDumpEmptyMap(t *testing.T) {
	m := map[string]int{}
	result := Sdump(m)
	if !strings.Contains(result, "len=0") {
		t.Errorf("expected len=0 in output, got: %s", result)
	}
}

func TestDumpEmptyStruct(t *testing.T) {
	type Empty struct{}
	result := Sdump(Empty{})
	if !strings.Contains(result, "Empty") {
		t.Errorf("expected Empty in output, got: %s", result)
	}
}

func TestDumpCustomIndent(t *testing.T) {
	type S struct {
		X int
	}
	cs := ConfigState{Indent: "    "}
	result := cs.Sdump(S{X: 1})
	if !strings.Contains(result, "    X:") {
		t.Errorf("expected 4-space indent, got: %s", result)
	}
}

func TestDumpSpewKeys(t *testing.T) {
	cs := ConfigState{Indent: " ", SortKeys: true, SpewKeys: true}
	m := map[int]string{1: "one", 2: "two"}
	result := cs.Sdump(m)
	// SpewKeys should dump keys using full type notation.
	if !strings.Contains(result, "(int)") {
		t.Errorf("expected (int) key format with SpewKeys, got: %s", result)
	}
}

func TestDumpPointerToPointer(t *testing.T) {
	x := 42
	p := &x
	pp := &p
	result := Sdump(pp)
	if !strings.Contains(result, "(**int)") {
		t.Errorf("expected (**int) in output, got: %s", result)
	}
	if !strings.Contains(result, "(int) 42") {
		t.Errorf("expected (int) 42 in output, got: %s", result)
	}
}

func TestDumpByteSlice(t *testing.T) {
	b := []byte{0x48, 0x65, 0x6c}
	result := Sdump(b)
	if !strings.Contains(result, "[]uint8") {
		t.Errorf("expected []uint8 in output, got: %s", result)
	}
}

func TestDumpStringerInterface(t *testing.T) {
	cs := ConfigState{Indent: " "}
	result := cs.Sdump(myStringer{})
	if !strings.Contains(result, "hello from stringer") {
		t.Errorf("expected stringer output, got: %s", result)
	}
}

type myStringer struct{}

func (myStringer) String() string { return "hello from stringer" }

func TestDumpDisableMethods(t *testing.T) {
	cs := ConfigState{Indent: " ", DisableMethods: true}
	result := cs.Sdump(myStringer{})
	if strings.Contains(result, "hello from stringer") {
		t.Errorf("expected no stringer output with DisableMethods, got: %s", result)
	}
}

func TestDumpErrorInterface(t *testing.T) {
	cs := ConfigState{Indent: " "}
	result := cs.Sdump(myError{})
	if !strings.Contains(result, "my error") {
		t.Errorf("expected error output, got: %s", result)
	}
}

type myError struct{}

func (myError) Error() string { return "my error" }

func TestDumpContinueOnMethod(t *testing.T) {
	cs := ConfigState{Indent: " ", ContinueOnMethod: true}
	result := cs.Sdump(myStringer{})
	// Should show both the stringer result and the struct internals.
	if !strings.Contains(result, "hello from stringer") {
		t.Errorf("expected stringer output, got: %s", result)
	}
	if !strings.Contains(result, "myStringer") {
		t.Errorf("expected struct name after ContinueOnMethod, got: %s", result)
	}
}

func TestDumpDeepCircular(t *testing.T) {
	type Node struct {
		Value int
		Next  *Node
	}
	a := &Node{Value: 1}
	b := &Node{Value: 2}
	c := &Node{Value: 3}
	a.Next = b
	b.Next = c
	c.Next = a // circular back to a

	result := Sdump(a)
	if !strings.Contains(result, "<already shown>") {
		t.Errorf("expected circular ref detection, got: %s", result)
	}
	if !strings.Contains(result, "(int) 1") {
		t.Errorf("expected value 1, got: %s", result)
	}
	if !strings.Contains(result, "(int) 2") {
		t.Errorf("expected value 2, got: %s", result)
	}
	if !strings.Contains(result, "(int) 3") {
		t.Errorf("expected value 3, got: %s", result)
	}
}

func TestDumpIntMapKeys(t *testing.T) {
	cs := ConfigState{Indent: " ", SortKeys: true}
	m := map[int]string{3: "c", 1: "a", 2: "b"}
	result := cs.Sdump(m)
	idx1 := strings.Index(result, "(1)")
	idx2 := strings.Index(result, "(2)")
	idx3 := strings.Index(result, "(3)")
	if idx1 == -1 || idx2 == -1 || idx3 == -1 {
		t.Fatalf("expected all keys in output, got: %s", result)
	}
	if idx1 >= idx2 || idx2 >= idx3 {
		t.Errorf("expected sorted int keys 1 < 2 < 3, got: %s", result)
	}
}
