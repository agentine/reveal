package reveal

import (
	"testing"
)

func BenchmarkSdumpInt(b *testing.B) {
	cs := ConfigState{Indent: " "}
	for i := 0; i < b.N; i++ {
		cs.Sdump(42)
	}
}

func BenchmarkSdumpString(b *testing.B) {
	cs := ConfigState{Indent: " "}
	s := "hello world this is a benchmark string"
	for i := 0; i < b.N; i++ {
		cs.Sdump(s)
	}
}

func BenchmarkSdumpStruct(b *testing.B) {
	type Point struct {
		X int
		Y int
		Z int
	}
	cs := ConfigState{Indent: " "}
	p := Point{X: 1, Y: 2, Z: 3}
	for i := 0; i < b.N; i++ {
		cs.Sdump(p)
	}
}

func BenchmarkSdumpSlice(b *testing.B) {
	cs := ConfigState{Indent: " "}
	s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for i := 0; i < b.N; i++ {
		cs.Sdump(s)
	}
}

func BenchmarkSdumpMap(b *testing.B) {
	cs := ConfigState{Indent: " ", SortKeys: true}
	m := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4, "e": 5}
	for i := 0; i < b.N; i++ {
		cs.Sdump(m)
	}
}

func BenchmarkSdumpNested(b *testing.B) {
	type Inner struct {
		Value int
		Name  string
	}
	type Outer struct {
		Items []Inner
		Count int
	}
	cs := ConfigState{Indent: " "}
	data := Outer{
		Items: []Inner{{1, "a"}, {2, "b"}, {3, "c"}},
		Count: 3,
	}
	for i := 0; i < b.N; i++ {
		cs.Sdump(data)
	}
}

func BenchmarkSdumpPointer(b *testing.B) {
	cs := ConfigState{Indent: " "}
	x := 42
	p := &x
	for i := 0; i < b.N; i++ {
		cs.Sdump(p)
	}
}
