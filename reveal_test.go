package reveal

import (
	"bytes"
	"strings"
	"testing"
)

func TestDumpBasicTypes(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		contains []string
	}{
		{"nil", nil, []string{"<nil>"}},
		{"bool true", true, []string{"(bool)", "true"}},
		{"bool false", false, []string{"(bool)", "false"}},
		{"int", 42, []string{"(int)", "42"}},
		{"int8", int8(-1), []string{"(int8)", "-1"}},
		{"int16", int16(256), []string{"(int16)", "256"}},
		{"int32", int32(65536), []string{"(int32)", "65536"}},
		{"int64", int64(1<<40), []string{"(int64)"}},
		{"uint", uint(42), []string{"(uint)", "42"}},
		{"uint8", uint8(255), []string{"(uint8)", "255"}},
		{"uint16", uint16(65535), []string{"(uint16)", "65535"}},
		{"uint32", uint32(100), []string{"(uint32)", "100"}},
		{"uint64", uint64(100), []string{"(uint64)", "100"}},
		{"float32", float32(3.14), []string{"(float32)"}},
		{"float64", float64(3.14159), []string{"(float64)", "3.14159"}},
		{"complex64", complex64(1 + 2i), []string{"(complex64)"}},
		{"complex128", complex128(3 + 4i), []string{"(complex128)"}},
		{"string", "hello", []string{"(string)", `"hello"`, "len=5"}},
		{"empty string", "", []string{"(string)", `""`, "len=0"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Sdump(tt.input)
			for _, want := range tt.contains {
				if !strings.Contains(result, want) {
					t.Errorf("Sdump(%v) = %q, want to contain %q", tt.input, result, want)
				}
			}
		})
	}
}

func TestDumpSlice(t *testing.T) {
	s := []int{1, 2, 3}
	result := Sdump(s)
	if !strings.Contains(result, "[]int") {
		t.Errorf("expected []int in output, got: %s", result)
	}
	if !strings.Contains(result, "len=3") {
		t.Errorf("expected len=3 in output, got: %s", result)
	}
	if !strings.Contains(result, "cap=3") {
		t.Errorf("expected cap=3 in output, got: %s", result)
	}
	if !strings.Contains(result, "(int) 1") {
		t.Errorf("expected (int) 1 in output, got: %s", result)
	}
}

func TestDumpNilSlice(t *testing.T) {
	var s []int
	result := Sdump(s)
	if !strings.Contains(result, "<nil>") {
		t.Errorf("expected <nil> in output, got: %s", result)
	}
}

func TestDumpArray(t *testing.T) {
	a := [3]int{10, 20, 30}
	result := Sdump(a)
	if !strings.Contains(result, "[3]int") {
		t.Errorf("expected [3]int in output, got: %s", result)
	}
	if !strings.Contains(result, "len=3") {
		t.Errorf("expected len=3 in output, got: %s", result)
	}
	// Arrays should not show cap.
	if strings.Contains(result, "cap=") {
		t.Errorf("expected no cap for array, got: %s", result)
	}
}

func TestDumpMap(t *testing.T) {
	// Use SortKeys to get deterministic output.
	cs := ConfigState{Indent: " ", SortKeys: true}
	m := map[string]int{"a": 1, "b": 2}
	result := cs.Sdump(m)
	if !strings.Contains(result, "map[string]int") {
		t.Errorf("expected map[string]int in output, got: %s", result)
	}
	if !strings.Contains(result, "len=2") {
		t.Errorf("expected len=2 in output, got: %s", result)
	}
}

func TestDumpNilMap(t *testing.T) {
	var m map[string]int
	result := Sdump(m)
	if !strings.Contains(result, "<nil>") {
		t.Errorf("expected <nil> in output, got: %s", result)
	}
}

func TestDumpStruct(t *testing.T) {
	type Point struct {
		X int
		Y int
	}
	p := Point{X: 10, Y: 20}
	result := Sdump(p)
	if !strings.Contains(result, "Point") {
		t.Errorf("expected Point in output, got: %s", result)
	}
	if !strings.Contains(result, "X:") {
		t.Errorf("expected X: in output, got: %s", result)
	}
	if !strings.Contains(result, "(int) 10") {
		t.Errorf("expected (int) 10 in output, got: %s", result)
	}
}

func TestDumpPointer(t *testing.T) {
	x := 42
	result := Sdump(&x)
	if !strings.Contains(result, "(*int)") {
		t.Errorf("expected (*int) in output, got: %s", result)
	}
	if !strings.Contains(result, "0x") {
		t.Errorf("expected pointer address in output, got: %s", result)
	}
	if !strings.Contains(result, "(int) 42") {
		t.Errorf("expected (int) 42 in output, got: %s", result)
	}
}

func TestDumpNilPointer(t *testing.T) {
	var p *int
	result := Sdump(p)
	if !strings.Contains(result, "(*int)") {
		t.Errorf("expected (*int) in output, got: %s", result)
	}
	if !strings.Contains(result, "<nil>") {
		t.Errorf("expected <nil> in output, got: %s", result)
	}
}

func TestDumpCircularReference(t *testing.T) {
	type Node struct {
		Value int
		Next  *Node
	}
	a := &Node{Value: 1}
	b := &Node{Value: 2}
	a.Next = b
	b.Next = a // circular

	result := Sdump(a)
	if !strings.Contains(result, "<already shown>") {
		t.Errorf("expected circular reference detection in output, got: %s", result)
	}
}

func TestDumpMaxDepth(t *testing.T) {
	type Nested struct {
		Value int
		Inner *Nested
	}
	n := &Nested{Value: 1, Inner: &Nested{Value: 2, Inner: &Nested{Value: 3}}}

	cs := ConfigState{Indent: " ", MaxDepth: 2}
	result := cs.Sdump(n)
	if !strings.Contains(result, "<max depth reached>") {
		t.Errorf("expected max depth message in output, got: %s", result)
	}
}

func TestDumpDisablePointerAddresses(t *testing.T) {
	x := 42
	cs := ConfigState{Indent: " ", DisablePointerAddresses: true}
	result := cs.Sdump(&x)
	if strings.Contains(result, "0x") {
		t.Errorf("expected no pointer address in output, got: %s", result)
	}
}

func TestDumpDisableCapacities(t *testing.T) {
	s := []int{1, 2, 3}
	cs := ConfigState{Indent: " ", DisableCapacities: true}
	result := cs.Sdump(s)
	if strings.Contains(result, "cap=") {
		t.Errorf("expected no capacity in output, got: %s", result)
	}
}

func TestDumpChan(t *testing.T) {
	ch := make(chan int, 5)
	result := Sdump(ch)
	if !strings.Contains(result, "chan int") {
		t.Errorf("expected chan int in output, got: %s", result)
	}
}

func TestDumpNilChan(t *testing.T) {
	var ch chan int
	result := Sdump(ch)
	if !strings.Contains(result, "<nil>") {
		t.Errorf("expected <nil> in output, got: %s", result)
	}
}

func TestDumpFunc(t *testing.T) {
	f := func() {}
	result := Sdump(f)
	if !strings.Contains(result, "func()") {
		t.Errorf("expected func() in output, got: %s", result)
	}
}

func TestDumpNilFunc(t *testing.T) {
	var f func()
	result := Sdump(f)
	if !strings.Contains(result, "<nil>") {
		t.Errorf("expected <nil> in output, got: %s", result)
	}
}

func TestFdump(t *testing.T) {
	var buf bytes.Buffer
	Fdump(&buf, 42)
	result := buf.String()
	if !strings.Contains(result, "(int) 42") {
		t.Errorf("Fdump expected (int) 42, got: %s", result)
	}
}

func TestDumpMultipleArgs(t *testing.T) {
	result := Sdump(1, "hello", true)
	if !strings.Contains(result, "(int) 1") {
		t.Errorf("expected (int) 1 in output, got: %s", result)
	}
	if !strings.Contains(result, `"hello"`) {
		t.Errorf("expected hello in output, got: %s", result)
	}
	if !strings.Contains(result, "(bool) true") {
		t.Errorf("expected (bool) true in output, got: %s", result)
	}
}

func TestDumpSortKeys(t *testing.T) {
	cs := ConfigState{Indent: " ", SortKeys: true}
	m := map[string]int{"c": 3, "a": 1, "b": 2}
	result := cs.Sdump(m)

	idxA := strings.Index(result, "(a)")
	idxB := strings.Index(result, "(b)")
	idxC := strings.Index(result, "(c)")
	if idxA == -1 || idxB == -1 || idxC == -1 {
		t.Fatalf("expected all keys in output, got: %s", result)
	}
	if !(idxA < idxB && idxB < idxC) {
		t.Errorf("expected sorted key order a < b < c, got: %s", result)
	}
}

func TestDumpInterface(t *testing.T) {
	var i interface{} = 42
	result := Sdump(i)
	if !strings.Contains(result, "(int) 42") {
		t.Errorf("expected (int) 42 in output, got: %s", result)
	}
}

func TestDumpNilInterface(t *testing.T) {
	result := Sdump(nil)
	if !strings.Contains(result, "<nil>") {
		t.Errorf("expected <nil> in output, got: %s", result)
	}
}
