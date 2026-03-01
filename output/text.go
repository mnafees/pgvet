package output

import (
	"fmt"
	"io"

	"github.com/mnafees/pgvet/rule"
)

// Text writes diagnostics in a human-readable format.
func Text(w io.Writer, diags []rule.Diagnostic) {
	for _, d := range diags {
		fmt.Fprintf(w, "%s:%d:%d: %s: [%s] %s\n",
			d.File, d.Line, d.Col, d.Severity, d.Rule, d.Message)
	}
}
