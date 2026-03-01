package rule

// All returns the default set of built-in rules.
func All() []Rule {
	return []Rule{
		&SelectStar{},
		&LimitWithoutOrder{},
		&NotInSubquery{},
		&ForUpdateNoSkip{},
		&DistinctOnOrder{},
		&NullComparison{},
		&UpdateWithoutWhere{},
		&DeleteWithoutWhere{},
		&InsertWithoutColumns{},
		&BanCharType{},
		&TimestampWithoutTimezone{},
		&OrderByOrdinal{},
		&GroupByOrdinal{},
		&LikeStartsWithWildcard{},
		&OffsetWithoutLimit{},
	}
}

// Extra returns opt-in rules not included in the default set.
// Use --rules multi-statement to enable them explicitly.
func Extra() []Rule {
	return []Rule{
		&MultiStatement{},
	}
}

// AllIncludingExtra returns all rules (default + opt-in).
func AllIncludingExtra() []Rule {
	return append(All(), Extra()...)
}
