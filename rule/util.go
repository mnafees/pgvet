package rule

// offsetToLineCol converts a byte offset in sql to a 1-based line and column.
func offsetToLineCol(sql string, offset int) (line, col int) {
	if offset < 0 || offset > len(sql) {
		return 1, 1
	}
	line = 1
	col = 1
	for i := 0; i < offset && i < len(sql); i++ {
		if sql[i] == '\n' {
			line++
			col = 1
		} else {
			col++
		}
	}
	return line, col
}
