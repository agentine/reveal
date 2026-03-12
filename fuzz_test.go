package reveal

import (
	"strings"
	"testing"
)

// TestRecoverPanics verifies that a panicking Stringer is caught.
func TestRecoverPanics(t *testing.T) {
	cs := ConfigState{Indent: " ", RecoverPanics: true}
	result := cs.Sdump(panicStringer{})
	if !strings.Contains(result, "<panic in method:") {
		t.Errorf("expected panic recovery message, got: %s", result)
	}
}

type panicStringer struct{}

func (panicStringer) String() string {
	panic("intentional panic in Stringer")
}

// TestRecoverPanicsDisabled verifies panics propagate when recovery is off.
func TestRecoverPanicsDisabled(t *testing.T) {
	cs := ConfigState{Indent: " ", RecoverPanics: false}
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic to propagate when RecoverPanics is false")
		}
	}()
	_ = cs.Sdump(panicStringer{})
}

// TestRecoverPanicsError verifies panic recovery for Error() method.
func TestRecoverPanicsError(t *testing.T) {
	cs := ConfigState{Indent: " ", RecoverPanics: true}
	result := cs.Sdump(panicError{})
	if !strings.Contains(result, "<panic in method:") {
		t.Errorf("expected panic recovery for Error(), got: %s", result)
	}
}

type panicError struct{}

func (panicError) Error() string {
	panic("intentional panic in Error")
}

// TestMaxSize verifies output truncation.
func TestMaxSize(t *testing.T) {
	cs := ConfigState{Indent: " ", MaxSize: 50}
	bigSlice := make([]int, 100)
	result := cs.Sdump(bigSlice)
	if !strings.Contains(result, "<output truncated") {
		t.Errorf("expected truncation message, got: %s", result)
	}
	// Output should not exceed MaxSize by too much (truncation message adds some).
	if len(result) > 200 {
		t.Errorf("output too large (%d bytes), MaxSize should limit it", len(result))
	}
}

// TestMaxSizeUnlimited verifies no truncation when MaxSize is 0.
func TestMaxSizeUnlimited(t *testing.T) {
	cs := ConfigState{Indent: " ", MaxSize: 0}
	bigSlice := make([]int, 100)
	result := cs.Sdump(bigSlice)
	if strings.Contains(result, "<output truncated") {
		t.Errorf("unexpected truncation with MaxSize=0")
	}
}

// FuzzDump fuzz tests the Dump function with string inputs.
func FuzzDump(f *testing.F) {
	// Seed corpus.
	f.Add("")
	f.Add("hello")
	f.Add("a very long string with special chars: \x00\x01\xff")
	f.Add(strings.Repeat("x", 10000))

	f.Fuzz(func(t *testing.T, s string) {
		// Should never panic.
		cs := ConfigState{Indent: " ", MaxSize: 1024, RecoverPanics: true}
		_ = cs.Sdump(s)
	})
}

// TestDumpNilVariants ensures no panics with nil values of various types.
func TestDumpNilVariants(t *testing.T) {
	cs := ConfigState{Indent: " ", RecoverPanics: true}

	var nilSlice []int
	var nilMap map[string]int
	var nilChan chan int
	var nilFunc func()
	var nilPtr *int
	var nilInterface interface{}

	// None of these should panic.
	_ = cs.Sdump(nilSlice)
	_ = cs.Sdump(nilMap)
	_ = cs.Sdump(nilChan)
	_ = cs.Sdump(nilFunc)
	_ = cs.Sdump(nilPtr)
	_ = cs.Sdump(nilInterface)
	_ = cs.Sdump(nil)
}

// TestDumpCircularMap ensures no panic/infinite loop with self-referencing map.
func TestDumpCircularMap(t *testing.T) {
	cs := ConfigState{Indent: " ", MaxDepth: 5}
	m := make(map[string]interface{})
	m["self"] = m
	result := cs.Sdump(m)
	// Should complete without hanging. Max depth should stop it.
	if result == "" {
		t.Error("expected non-empty output for circular map")
	}
}

// TestDumpDeepNesting ensures MaxDepth prevents stack overflow.
func TestDumpDeepNesting(t *testing.T) {
	type Node struct {
		Next *Node
	}
	// Build a chain of 1000 nodes.
	var head *Node
	for i := 0; i < 1000; i++ {
		head = &Node{Next: head}
	}
	cs := ConfigState{Indent: " ", MaxDepth: 10}
	result := cs.Sdump(head)
	if !strings.Contains(result, "<max depth reached>") {
		t.Error("expected max depth reached for deep nesting")
	}
}
