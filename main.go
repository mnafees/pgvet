package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mnafees/pgvet/analyzer"
	"github.com/mnafees/pgvet/output"
	"github.com/mnafees/pgvet/rule"
)

func main() {
	format := flag.String("format", "text", "Output format: text or json")
	rulesFlag := flag.String("rules", "", "Comma-separated list of rules to run (default: all)")
	excludeFlag := flag.String("exclude", "", "Comma-separated list of rules to exclude")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: pgvet [flags] <file-or-dir>...\n\n")
		fmt.Fprintf(os.Stderr, "A static analysis tool for PostgreSQL SQL files.\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nDefault rules:\n")
		for _, r := range rule.All() {
			fmt.Fprintf(os.Stderr, "  %-25s %s\n", r.Name(), r.Description())
		}
		fmt.Fprintf(os.Stderr, "\nOpt-in rules (use --rules to enable):\n")
		for _, r := range rule.Extra() {
			fmt.Fprintf(os.Stderr, "  %-25s %s\n", r.Name(), r.Description())
		}
	}
	flag.Parse()

	rules := selectRules(*rulesFlag, *excludeFlag)
	a := analyzer.New(rules)

	var diags []rule.Diagnostic
	var err error

	if flag.NArg() == 0 {
		// Read from stdin.
		data, readErr := io.ReadAll(os.Stdin)
		if readErr != nil {
			fmt.Fprintf(os.Stderr, "error reading stdin: %v\n", readErr)
			os.Exit(2)
		}
		diags, err = a.AnalyzeStdin(string(data))
	} else {
		diags, err = a.AnalyzePaths(flag.Args())
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	switch *format {
	case "json":
		if err := output.JSON(os.Stdout, diags); err != nil {
			fmt.Fprintf(os.Stderr, "error writing JSON: %v\n", err)
			os.Exit(2)
		}
	default:
		output.Text(os.Stdout, diags)
	}

	if len(diags) > 0 {
		os.Exit(1)
	}
}

func selectRules(include, exclude string) []rule.Rule {
	includeSet := parseCSV(include)
	excludeSet := parseCSV(exclude)

	// When --rules is specified, search the full pool (including opt-in rules).
	// Otherwise, only use the default set.
	pool := rule.All()
	if len(includeSet) > 0 {
		pool = rule.AllIncludingExtra()
	}

	var selected []rule.Rule
	for _, r := range pool {
		if len(includeSet) > 0 && !includeSet[r.Name()] {
			continue
		}
		if excludeSet[r.Name()] {
			continue
		}
		selected = append(selected, r)
	}
	return selected
}

func parseCSV(s string) map[string]bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	m := make(map[string]bool)
	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(part)
		if part != "" {
			m[part] = true
		}
	}
	return m
}
