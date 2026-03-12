package reveal

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestNewFormatter(t *testing.T) {
	f := NewFormatter(42)
	if f == nil {
		t.Fatal("NewFormatter returned nil")
	}
}

func TestFormatterDumpVerb(t *testing.T) {
	result := fmt.Sprintf("%d", NewFormatter(42))
	if !strings.Contains(result, "(int) 42") {
		t.Errorf("expected dump format with %%d, got: %s", result)
	}
}

func TestFormatterVVerb(t *testing.T) {
	result := fmt.Sprintf("%v", NewFormatter(42))
	if !strings.Contains(result, "42") {
		t.Errorf("expected value with %%v, got: %s", result)
	}
}

func TestFormatterPlusVVerb(t *testing.T) {
	type S struct {
		X int
	}
	result := fmt.Sprintf("%+v", NewFormatter(S{X: 1}))
	if !strings.Contains(result, "X:") {
		t.Errorf("expected field names with %%+v, got: %s", result)
	}
}

func TestFormatterHashVVerb(t *testing.T) {
	type S struct {
		X int
	}
	result := fmt.Sprintf("%#v", NewFormatter(S{X: 1}))
	if !strings.Contains(result, "reveal.S") {
		t.Errorf("expected Go syntax with %%#v, got: %s", result)
	}
}

func TestFormatterSVerb(t *testing.T) {
	result := fmt.Sprintf("%s", NewFormatter(42))
	if !strings.Contains(result, "(int) 42") {
		t.Errorf("expected dump format with %%s, got: %s", result)
	}
}

func TestSprint(t *testing.T) {
	result := Sprint(42)
	if !strings.Contains(result, "(int) 42") {
		t.Errorf("Sprint(42) = %q, expected to contain (int) 42", result)
	}
}

func TestSprintf(t *testing.T) {
	result := Sprintf("value: %d", 42)
	if !strings.Contains(result, "(int) 42") {
		t.Errorf("Sprintf = %q, expected dump format", result)
	}
}

func TestSprintln(t *testing.T) {
	result := Sprintln(42)
	if !strings.Contains(result, "(int) 42") {
		t.Errorf("Sprintln = %q, expected dump format", result)
	}
	if !strings.HasSuffix(result, "\n") {
		t.Errorf("Sprintln should end with newline, got: %q", result)
	}
}

func TestFprint(t *testing.T) {
	var buf bytes.Buffer
	n, err := Fprint(&buf, 42)
	if err != nil {
		t.Fatalf("Fprint error: %v", err)
	}
	if n == 0 {
		t.Error("Fprint wrote 0 bytes")
	}
	if !strings.Contains(buf.String(), "(int) 42") {
		t.Errorf("Fprint = %q, expected dump format", buf.String())
	}
}

func TestFprintf(t *testing.T) {
	var buf bytes.Buffer
	_, err := Fprintf(&buf, "val: %d", 42)
	if err != nil {
		t.Fatalf("Fprintf error: %v", err)
	}
	if !strings.Contains(buf.String(), "(int) 42") {
		t.Errorf("Fprintf = %q, expected dump format", buf.String())
	}
}

func TestFprintln(t *testing.T) {
	var buf bytes.Buffer
	_, err := Fprintln(&buf, 42)
	if err != nil {
		t.Fatalf("Fprintln error: %v", err)
	}
	if !strings.Contains(buf.String(), "(int) 42") {
		t.Errorf("Fprintln = %q, expected dump format", buf.String())
	}
}

func TestConfigStateSprint(t *testing.T) {
	cs := ConfigState{Indent: " "}
	result := cs.Sprint(42)
	if !strings.Contains(result, "(int) 42") {
		t.Errorf("ConfigState.Sprint = %q, expected dump format", result)
	}
}

func TestConfigStateSprintf(t *testing.T) {
	cs := ConfigState{Indent: " "}
	result := cs.Sprintf("val: %d", 42)
	if !strings.Contains(result, "(int) 42") {
		t.Errorf("ConfigState.Sprintf = %q, expected dump format", result)
	}
}

func TestConfigStateNewFormatter(t *testing.T) {
	cs := ConfigState{Indent: "  "}
	f := cs.NewFormatter(42)
	result := fmt.Sprintf("%d", f)
	if !strings.Contains(result, "(int) 42") {
		t.Errorf("ConfigState.NewFormatter %%d = %q", result)
	}
}

func TestSprintMultipleArgs(t *testing.T) {
	result := Sprint(1, "hello")
	if !strings.Contains(result, "(int) 1") {
		t.Errorf("expected (int) 1, got: %s", result)
	}
	if !strings.Contains(result, "(string)") {
		t.Errorf("expected (string), got: %s", result)
	}
}

func TestFormatterNil(t *testing.T) {
	result := fmt.Sprintf("%d", NewFormatter(nil))
	if !strings.Contains(result, "<nil>") {
		t.Errorf("expected <nil> for nil formatter, got: %s", result)
	}
}

func TestFormatterSlice(t *testing.T) {
	s := []int{1, 2, 3}
	result := fmt.Sprintf("%d", NewFormatter(s))
	if !strings.Contains(result, "[]int") {
		t.Errorf("expected []int in dump format, got: %s", result)
	}
}
