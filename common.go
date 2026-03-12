package reveal

import (
	"reflect"
	"sort"
	"strconv"
)

// pointerTracker tracks visited pointers for circular reference detection.
type pointerTracker struct {
	seen map[uintptr]bool
}

func newPointerTracker() *pointerTracker {
	return &pointerTracker{seen: make(map[uintptr]bool)}
}

// visit returns true if this pointer has already been visited (circular reference).
func (pt *pointerTracker) visit(ptr uintptr) bool {
	if pt.seen[ptr] {
		return true
	}
	pt.seen[ptr] = true
	return false
}

// leave removes a pointer from the visited set when leaving a scope.
func (pt *pointerTracker) leave(ptr uintptr) {
	delete(pt.seen, ptr)
}

// isPointerType returns true if the kind is a type that has a pointer address
// worth tracking for circular reference detection.
func isPointerType(k reflect.Kind) bool {
	switch k {
	case reflect.Ptr, reflect.Map, reflect.Slice:
		return true
	default:
		return false
	}
}

// valueLen returns a printable length/cap string for array-like types.
func valueLen(v reflect.Value) (int, int) {
	return v.Len(), v.Cap()
}

// sortMapKeys sorts map keys for deterministic output.
// Keys must all be the same type. Supports string, int, uint, float, and bool keys.
func sortMapKeys(keys []reflect.Value) {
	if len(keys) == 0 {
		return
	}
	sort.SliceStable(keys, func(i, j int) bool {
		return mapKeyLess(keys[i], keys[j])
	})
}

// mapKeyLess compares two map keys for sorting.
func mapKeyLess(a, b reflect.Value) bool {
	if a.Kind() != b.Kind() {
		return false
	}
	switch a.Kind() {
	case reflect.String:
		return a.String() < b.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return a.Int() < b.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return a.Uint() < b.Uint()
	case reflect.Float32, reflect.Float64:
		return a.Float() < b.Float()
	case reflect.Bool:
		return !a.Bool() && b.Bool()
	default:
		return false
	}
}

// typeString returns a human-friendly type name for a reflect.Type.
func typeString(t reflect.Type) string {
	return t.String()
}

// hexPointer returns a hex representation of a pointer address.
func hexPointer(v reflect.Value) string {
	var ptr uintptr
	switch v.Kind() {
	case reflect.Ptr, reflect.Map, reflect.Func, reflect.Chan, reflect.Slice, reflect.UnsafePointer:
		ptr = v.Pointer()
	default:
		return ""
	}
	return "0x" + strconv.FormatUint(uint64(ptr), 16)
}

// handleMethods attempts to invoke error/Stringer interfaces on the value
// and returns the result string and true if a method was called.
func handleMethods(cs *ConfigState, v reflect.Value) (string, bool) {
	if cs.DisableMethods {
		return "", false
	}

	if !v.IsValid() {
		return "", false
	}

	// If the value is not a pointer, and DisablePointerMethods is false,
	// try to get an addressable copy to check for pointer receiver methods.
	if !cs.DisablePointerMethods && v.CanAddr() {
		pv := v.Addr()
		if pv.CanInterface() {
			if e, ok := pv.Interface().(error); ok {
				return e.Error(), true
			}
			if s, ok := pv.Interface().(interface{ String() string }); ok {
				return s.String(), true
			}
		}
	}

	if v.CanInterface() {
		if e, ok := v.Interface().(error); ok {
			return e.Error(), true
		}
		if s, ok := v.Interface().(interface{ String() string }); ok {
			return s.String(), true
		}
	}

	return "", false
}
