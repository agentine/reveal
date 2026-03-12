# reveal

A modern, maintained, and bug-fixed deep pretty printer for Go data structures. Drop-in replacement for [go-spew](https://github.com/davecgh/go-spew).

## Installation

```bash
go get github.com/agentine/reveal
```

## Quick Start

```go
import "github.com/agentine/reveal"

reveal.Dump(myStruct)          // print to stdout
s := reveal.Sdump(myStruct)    // return as string
reveal.Fdump(w, myStruct)      // print to io.Writer
```

## Migration from go-spew

Change one import:

```go
// Before
import "github.com/davecgh/go-spew/spew"
spew.Dump(myStruct)

// After
import "github.com/agentine/reveal"
reveal.Dump(myStruct)
```

Output format is identical by default. All go-spew API functions are supported.

## API

### Package-level functions

| Function | Description |
|---|---|
| `Dump(a ...)` | Print to stdout |
| `Sdump(a ...) string` | Return as string |
| `Fdump(w, a ...)` | Print to io.Writer |
| `Sprint(a ...) string` | Deep-aware fmt.Sprint |
| `Sprintf(format, a ...) string` | Deep-aware fmt.Sprintf |
| `Sprintln(a ...) string` | Deep-aware fmt.Sprintln |
| `Fprint(w, a ...) (int, error)` | Deep-aware fmt.Fprint |
| `Fprintf(w, format, a ...) (int, error)` | Deep-aware fmt.Fprintf |
| `Fprintln(w, a ...) (int, error)` | Deep-aware fmt.Fprintln |
| `NewFormatter(v) fmt.Formatter` | Wrap value for use with fmt verbs |

### Configuration

```go
reveal.Config.MaxDepth = 5
reveal.Config.SortKeys = true
reveal.Config.Dump(myStruct)
```

| Field | Default | Description |
|---|---|---|
| `Indent` | `" "` | Indentation string |
| `MaxDepth` | `0` | Max recursion depth (0 = unlimited) |
| `DisableMethods` | `false` | Skip Stringer/Error methods |
| `DisablePointerMethods` | `false` | Skip pointer receiver methods |
| `DisablePointerAddresses` | `false` | Hide pointer addresses |
| `DisableCapacities` | `false` | Hide slice/map capacities |
| `ContinueOnMethod` | `false` | Show internals after method call |
| `SortKeys` | `false` | Sort map keys |
| `SpewKeys` | `false` | Dump-format map keys |
| `MaxSize` | `10MB` | Max output bytes (0 = unlimited) |
| `RecoverPanics` | `true` | Recover from Stringer/Error panics |
| `OmitNilPointers` | `false` | Suppress nil pointer struct fields |
| `OmitUnexported` | `false` | Suppress unexported struct fields |
| `HexIntegers` | `false` | Format integers as hex |

## Bug Fixes vs go-spew

- **Memory exhaustion** (#145): Bounded by `MaxSize` (default 10MB)
- **Panic with custom Stringer on maps** (#141, #115): Safe method invocation with `RecoverPanics`
- **Panic with wrapped custom errors** (#144): Panic recovery for Error() methods
- **Panic with unsorted private fields** (#108): Safe unexported field handling
- **Stringer recursion**: Detected and broken automatically

## License

MIT
