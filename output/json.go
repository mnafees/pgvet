package output

import (
	"encoding/json"
	"io"

	"github.com/mnafees/pgvet/rule"
)

// JSON writes diagnostics as a JSON array.
func JSON(w io.Writer, diags []rule.Diagnostic) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(diags)
}
