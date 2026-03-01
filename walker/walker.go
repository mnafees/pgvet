package walker

import (
	pg_query "github.com/pganalyze/pg_query_go/v6"
)

// Visitor is called for each node in the AST. Return false to stop traversal.
type Visitor func(node *pg_query.Node) bool

// Walk recursively visits every node in the parse tree rooted at node.
func Walk(node *pg_query.Node, visit Visitor) {
	if node == nil {
		return
	}
	if !visit(node) {
		return
	}
	for _, child := range Children(node) {
		Walk(child, visit)
	}
}

// Children returns the direct child nodes of a given node.
func Children(node *pg_query.Node) []*pg_query.Node {
	if node == nil {
		return nil
	}

	var children []*pg_query.Node
	add := func(n *pg_query.Node) {
		if n != nil {
			children = append(children, n)
		}
	}
	addAll := func(nodes []*pg_query.Node) {
		for _, n := range nodes {
			add(n)
		}
	}

	switch {
	case node.GetSelectStmt() != nil:
		s := node.GetSelectStmt()
		addAll(s.DistinctClause)
		addAll(s.TargetList)
		addAll(s.FromClause)
		add(s.WhereClause)
		addAll(s.GroupClause)
		add(s.HavingClause)
		addAll(s.WindowClause)
		addAll(s.ValuesLists)
		addAll(s.SortClause)
		add(s.LimitOffset)
		add(s.LimitCount)
		addAll(s.LockingClause)
		if s.WithClause != nil {
			for _, cte := range s.WithClause.Ctes {
				add(cte)
			}
		}
		if s.Larg != nil {
			children = append(children, makeSelectStmtNode(s.Larg))
		}
		if s.Rarg != nil {
			children = append(children, makeSelectStmtNode(s.Rarg))
		}

	case node.GetInsertStmt() != nil:
		s := node.GetInsertStmt()
		if s.Relation != nil {
			children = append(children, makeRangeVarNode(s.Relation))
		}
		addAll(s.Cols)
		add(s.SelectStmt)
		// OnConflictClause is a single struct, not a slice of nodes.
		// We don't recurse into it for now.
		addAll(s.ReturningList)
		if s.WithClause != nil {
			for _, cte := range s.WithClause.Ctes {
				add(cte)
			}
		}

	case node.GetUpdateStmt() != nil:
		s := node.GetUpdateStmt()
		if s.Relation != nil {
			children = append(children, makeRangeVarNode(s.Relation))
		}
		addAll(s.TargetList)
		add(s.WhereClause)
		addAll(s.FromClause)
		addAll(s.ReturningList)
		if s.WithClause != nil {
			for _, cte := range s.WithClause.Ctes {
				add(cte)
			}
		}

	case node.GetDeleteStmt() != nil:
		s := node.GetDeleteStmt()
		if s.Relation != nil {
			children = append(children, makeRangeVarNode(s.Relation))
		}
		addAll(s.UsingClause)
		add(s.WhereClause)
		addAll(s.ReturningList)
		if s.WithClause != nil {
			for _, cte := range s.WithClause.Ctes {
				add(cte)
			}
		}

	case node.GetCommonTableExpr() != nil:
		cte := node.GetCommonTableExpr()
		add(cte.Ctequery)
		addAll(cte.Aliascolnames)

	case node.GetJoinExpr() != nil:
		j := node.GetJoinExpr()
		add(j.Larg)
		add(j.Rarg)
		add(j.Quals)
		addAll(j.UsingClause)

	case node.GetRangeSubselect() != nil:
		rs := node.GetRangeSubselect()
		add(rs.Subquery)

	case node.GetSubLink() != nil:
		sl := node.GetSubLink()
		add(sl.Testexpr)
		add(sl.Subselect)

	case node.GetBoolExpr() != nil:
		be := node.GetBoolExpr()
		addAll(be.Args)

	case node.GetAExpr() != nil:
		ae := node.GetAExpr()
		add(ae.Lexpr)
		add(ae.Rexpr)
		addAll(ae.Name)

	case node.GetFuncCall() != nil:
		fc := node.GetFuncCall()
		addAll(fc.Funcname)
		addAll(fc.Args)
		addAll(fc.AggOrder)
		add(fc.AggFilter)

	case node.GetResTarget() != nil:
		rt := node.GetResTarget()
		add(rt.Val)
		addAll(rt.Indirection)

	case node.GetTypeCast() != nil:
		tc := node.GetTypeCast()
		add(tc.Arg)

	case node.GetColumnRef() != nil:
		cr := node.GetColumnRef()
		addAll(cr.Fields)

	case node.GetCaseExpr() != nil:
		ce := node.GetCaseExpr()
		add(ce.Arg)
		addAll(ce.Args)
		add(ce.Defresult)

	case node.GetCaseWhen() != nil:
		cw := node.GetCaseWhen()
		add(cw.Expr)
		add(cw.Result)

	case node.GetNullTest() != nil:
		nt := node.GetNullTest()
		add(nt.Arg)

	case node.GetCoalesceExpr() != nil:
		ce := node.GetCoalesceExpr()
		addAll(ce.Args)

	case node.GetSortBy() != nil:
		sb := node.GetSortBy()
		add(sb.Node)
		addAll(sb.UseOp)

	case node.GetLockingClause() != nil:
		lc := node.GetLockingClause()
		addAll(lc.LockedRels)

	case node.GetList() != nil:
		l := node.GetList()
		addAll(l.Items)

	case node.GetRowExpr() != nil:
		re := node.GetRowExpr()
		addAll(re.Args)

	case node.GetMinMaxExpr() != nil:
		mm := node.GetMinMaxExpr()
		addAll(mm.Args)

	case node.GetXmlExpr() != nil:
		xe := node.GetXmlExpr()
		addAll(xe.NamedArgs)
		addAll(xe.Args)

	case node.GetGroupingFunc() != nil:
		gf := node.GetGroupingFunc()
		addAll(gf.Args)
	}

	return children
}

func makeRangeVarNode(rv *pg_query.RangeVar) *pg_query.Node {
	return &pg_query.Node{
		Node: &pg_query.Node_RangeVar{RangeVar: rv},
	}
}

func makeSelectStmtNode(ss *pg_query.SelectStmt) *pg_query.Node {
	return &pg_query.Node{
		Node: &pg_query.Node_SelectStmt{SelectStmt: ss},
	}
}
