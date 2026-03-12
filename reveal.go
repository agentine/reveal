package reveal

import (
	"bytes"
	"io"
	"os"
)

// Dump displays the passed parameters to standard output with type info,
// pointer addresses, and nested structure. It uses the default configuration.
func Dump(a ...interface{}) {
	Config.Fdump(os.Stdout, a...)
}

// Sdump returns a string with the passed parameters formatted using Dump style.
// It uses the default configuration.
func Sdump(a ...interface{}) string {
	var buf bytes.Buffer
	Config.Fdump(&buf, a...)
	return buf.String()
}

// Fdump formats and writes the passed parameters to the given io.Writer.
// It uses the default configuration.
func Fdump(w io.Writer, a ...interface{}) {
	Config.Fdump(w, a...)
}
