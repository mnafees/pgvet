package analyzer

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	pg_query "github.com/pganalyze/pg_query_go/v6"

	"github.com/mnafees/pgvet/rule"
)

// Analyzer runs rules against SQL files.
type Analyzer struct {
	Rules []rule.Rule
}

// New creates an Analyzer with the given rules.
func New(rules []rule.Rule) *Analyzer {
	return &Analyzer{Rules: rules}
}

// AnalyzePaths analyzes one or more file or directory paths.
func (a *Analyzer) AnalyzePaths(paths []string) ([]rule.Diagnostic, error) {
	var allDiags []rule.Diagnostic

	for _, p := range paths {
		info, err := os.Stat(p)
		if err != nil {
			return nil, err
		}
		if info.IsDir() {
			diags, err := a.analyzeDir(p)
			if err != nil {
				return nil, err
			}
			allDiags = append(allDiags, diags...)
		} else {
			diags, err := a.analyzeFile(p)
			if err != nil {
				return nil, err
			}
			allDiags = append(allDiags, diags...)
		}
	}

	sort.Slice(allDiags, func(i, j int) bool {
		if allDiags[i].File != allDiags[j].File {
			return allDiags[i].File < allDiags[j].File
		}
		if allDiags[i].Line != allDiags[j].Line {
			return allDiags[i].Line < allDiags[j].Line
		}
		return allDiags[i].Col < allDiags[j].Col
	})

	return allDiags, nil
}

// AnalyzeStdin analyzes SQL from stdin.
func (a *Analyzer) AnalyzeStdin(sql string) ([]rule.Diagnostic, error) {
	return a.analyzeSQL("<stdin>", sql)
}

func (a *Analyzer) analyzeDir(dir string) ([]rule.Diagnostic, error) {
	var allDiags []rule.Diagnostic

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".sql") {
			return nil
		}
		diags, err := a.analyzeFile(path)
		if err != nil {
			return err
		}
		allDiags = append(allDiags, diags...)
		return nil
	})

	return allDiags, err
}

func (a *Analyzer) analyzeFile(path string) ([]rule.Diagnostic, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	diags, err := a.analyzeSQL(path, string(data))
	if err != nil {
		return nil, err
	}
	return diags, nil
}

func (a *Analyzer) analyzeSQL(file, sql string) ([]rule.Diagnostic, error) {
	// Split into individual statements.
	stmtTexts, err := pg_query.SplitWithParser(sql, true)
	if err != nil {
		// If splitting fails, try parsing the whole thing.
		stmtTexts = []string{sql}
	}

	var allDiags []rule.Diagnostic

	// Check multi-statement rule on the full SQL.
	fullResult, fullErr := pg_query.Parse(sql)
	if fullErr == nil && fullResult != nil {
		for _, r := range a.Rules {
			if ms, ok := r.(*rule.MultiStatement); ok {
				diags := ms.CheckMulti(fullResult.Stmts, sql)
				for i := range diags {
					diags[i].File = file
				}
				allDiags = append(allDiags, diags...)
			}
		}
	}

	// Run per-statement rules.
	for _, stmtSQL := range stmtTexts {
		stmtSQL = strings.TrimSpace(stmtSQL)
		if stmtSQL == "" {
			continue
		}

		result, err := pg_query.Parse(stmtSQL)
		if err != nil {
			// Report parse errors as diagnostics.
			allDiags = append(allDiags, rule.Diagnostic{
				Rule:     "parse-error",
				Message:  err.Error(),
				File:     file,
				Line:     1,
				Col:      1,
				Severity: rule.SeverityError,
			})
			continue
		}

		for _, stmt := range result.Stmts {
			for _, r := range a.Rules {
				// Skip multi-statement here — handled above.
				if _, ok := r.(*rule.MultiStatement); ok {
					continue
				}
				diags := r.Check(stmt, stmtSQL)
				for i := range diags {
					diags[i].File = file
				}
				allDiags = append(allDiags, diags...)
			}
		}
	}

	return allDiags, nil
}
