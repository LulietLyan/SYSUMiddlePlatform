package types

import (
	"regexp"

	"github.com/pingcap/parser/ast"
)

// AggregateFuncExpr like *ast.AggregateFuncExpr
type AggregateFuncExpr struct {
	// F is the function name.
	Function string
	// Args is the function args.
	Args []interface{} // []ExprNode
	// Distinct is true, function hence only aggregate distinct values.
	// For example, column c1 values are "1", "2", "2",  "sum(c1)" is "5", but "sum(distinct c1)" is "3".
	Distinct bool
	// Order is only used in GROUP_CONCAT
	Order *OrderByClause
}

// ParseAggregateFuncExpr 解析Aggregate Function表达式
func (afe *AggregateFuncExpr) ParseAggregateFuncExpr(aafe *ast.AggregateFuncExpr) *AggregateFuncExpr {
	if aafe == nil {
		return nil
	}

	// F string
	afe.Function = aafe.F

	// Args []ExprNode
	for _, aarg := range aafe.Args {
		if val, bv := classifyExprNode(aarg); bv == true {
			argitem := val

			afe.Args = append(afe.Args, argitem)
		}
	}

	// Distinct bool
	afe.Distinct = aafe.Distinct

	// Order *OrderByClause
	oc := new(OrderByClause)
	afe.Order = oc.ParseOrderByClause(aafe.Order)

	return afe
}

// BetweenExpr like *ast.BetweenExpr
type BetweenExpr struct {
	// Expr is the expression to be checked.
	Expr interface{} // ExprNode
	// Left is the expression for minimal value in the range.
	Left interface{} // ExprNode
	// Right is the expression for maximum value in the range.
	Right interface{} // ExprNode
	// Not is true, the expression is "not between and".
	Not bool
}

// ParseBetweenExpr 解析Between表达式
func (bwe *BetweenExpr) ParseBetweenExpr(abwe *ast.BetweenExpr) *BetweenExpr {
	if abwe == nil {
		return nil
	}

	// Expr ExprNode
	if val, bv := classifyExprNode(abwe.Expr); bv == true {
		bwe.Expr = val
	}

	// Left ExprNode
	if val, bv := classifyExprNode(abwe.Left); bv == true {
		bwe.Left = val
	}

	// Right ExprNode
	if val, bv := classifyExprNode(abwe.Right); bv == true {
		bwe.Right = val
	}

	// Not bool
	bwe.Not = abwe.Not

	return bwe
}

// BinaryOperationExpr like *ast.BinaryOperationExpr
type BinaryOperationExpr struct {
	// Op is the operator code for BinaryOperation.
	Operator string // opcode.Op
	// L is the left expression in BinaryOperation.
	Left interface{} // ExprNode
	// R is the right expression in BinaryOperation.
	Right interface{} // ExprNode
}

// ParseBinaryOperationExpr 解析Binary Operation表达式，where
func (boe *BinaryOperationExpr) ParseBinaryOperationExpr(aboe *ast.BinaryOperationExpr) *BinaryOperationExpr {
	if aboe == nil {
		return nil
	}

	// Op  opcode.Op
	if val, bv := classifyOpcode(aboe.Op); bv == true {
		boe.Operator = val
	}

	// L   ExprNode
	if val, bv := classifyExprNode(aboe.L); bv == true {
		boe.Left = val
	}

	// R   ExprNode
	if val, bv := classifyExprNode(aboe.R); bv == true {
		boe.Right = val
	}

	return boe
}

// CaseExpr like *ast.CaseExpr
type CaseExpr struct {
	// Value is the compare value expression.
	Value interface{} // ExprNode
	// WhenClauses is the condition check expression.
	WhenClauses []interface{} // []*WhenClause
	// ElseClause is the else result expression.
	ElseClause interface{} // ExprNode
	// contains filtered or unexported fields
}

// ParseCaseExpr 解析case表达式
func (ce *CaseExpr) ParseCaseExpr(ace *ast.CaseExpr) *CaseExpr {
	if ace == nil {
		return nil
	}

	// Value ExprNode
	if val, bv := classifyExprNode(ace.Value); bv == true {
		ce.Value = val
	}

	// WhenClauses []*WhenClause
	for _, awc := range ace.WhenClauses {
		wc := new(WhenClause)
		wcitem := wc.ParseWhenClause(awc)

		ce.WhenClauses = append(ce.WhenClauses, wcitem)
	}

	// ElseClause ExprNode
	if val, bv := classifyExprNode(ace.ElseClause); bv == true {
		ce.ElseClause = val
	}

	return ce
}

// SubqueryExpr like *ast.SubqueryExpr
type SubqueryExpr struct {
	// Query is the query SelectNode.
	Query      interface{} // ResultSetNode
	Evaluated  bool
	Correlated bool
	MultiRows  bool
	Exists     bool
}

// ParseSubqueryExpr 解析Subquery表达式
func (sqe *SubqueryExpr) ParseSubqueryExpr(asqe *ast.SubqueryExpr) *SubqueryExpr {
	if asqe == nil {
		return nil
	}

	// Query      ResultSetNode
	if val, bv := classifyResultSetNode(asqe.Query); bv == true {
		sqe.Query = val
	}

	// Evaluated  bool
	sqe.Evaluated = asqe.Evaluated

	// Correlated bool
	sqe.Correlated = asqe.Correlated

	// MultiRows  bool
	sqe.MultiRows = asqe.MultiRows

	// Exists     bool
	sqe.Exists = asqe.Exists

	return sqe
}

// ExistsSubqueryExpr like *ast.ExistsSubqueryExpr
type ExistsSubqueryExpr struct {
	// Sel is the subquery, may be rewritten to other type of expression.
	Sel interface{} // ExprNode
	// Not is true, the expression is "not exists".
	Not bool
}

// ParseExistsSubqueryExpr 解析Exists Subquery表达式
func (esqe *ExistsSubqueryExpr) ParseExistsSubqueryExpr(aesqe *ast.ExistsSubqueryExpr) *ExistsSubqueryExpr {
	if aesqe == nil {
		return nil
	}

	// Sel ExprNode
	if val, bv := classifyExprNode(aesqe.Sel); bv == true {
		esqe.Sel = val
	}

	// Not bool
	esqe.Not = aesqe.Not
	return esqe
}

// FuncCallExpr like *ast.FuncCallExpr
type FuncCallExpr struct {
	// Type   string            // FuncCallExprType
	// Schema map[string]string // model.CIStr
	// FnName is the function name.
	FunctionName map[string]string // model.CIStr
	// Args is the function args.
	Args []interface{} // []ExprNode
}

// ParseFuncCallExpr 解析Function Call表达式
func (fce *FuncCallExpr) ParseFuncCallExpr(afce *ast.FuncCallExpr) *FuncCallExpr {
	if afce == nil {
		return nil
	}

	// Tp     FuncCallExprType
	// Schema model.CIStr

	// FnName model.CIStr
	fce.FunctionName = beautifyModelCIStr(&afce.FnName)

	// Args []ExprNode
	for _, aarg := range afce.Args {
		if val, bv := classifyExprNode(aarg); bv == true {
			argitem := val

			fce.Args = append(fce.Args, argitem)
		}
	}

	return fce
}

// FuncCastExpr like *ast.FuncCastExpr
type FuncCastExpr struct {
	// Expr is the expression to be converted.
	Expr interface{} // ExprNode
	// Tp is the conversion type.
	Type *FieldType
	// FunctionType is either Cast, Convert or Binary.
	FunctionType string // CastFunctionType
	// ExplicitCharSet is true when charset is explicit indicated.
	//ExplicitCharSet bool
}

// ParseFuncCastExpr 解析Function Cast表达式
func (fcae *FuncCastExpr) ParseFuncCastExpr(afcae *ast.FuncCastExpr) *FuncCastExpr {
	if afcae == nil {
		return nil
	}

	// Expr ExprNode
	if val, bv := classifyExprNode(afcae.Expr); bv == true {
		fcae.Expr = val
	}

	// Type *types.FieldType
	ft := new(FieldType)
	fcae.Type = ft.ParseFieldType(afcae.Tp)

	// FunctionType CastFunctionType
	if val, bv := classifyCastFunctionType(afcae.FunctionType); bv == true {
		fcae.FunctionType = val
	}

	// ExplicitCharSet bool
	// fcae.ExplicitCharSet = afcae.ExplicitCharSet

	return fcae
}

// TableNameExpr like *ast.TableNameExpr
type TableNameExpr struct {
	// Name is the referenced object name expression.
	Name *TableName
}

// ParseTableNameExpr 解析Table Name表达式
func (tne *TableNameExpr) ParseTableNameExpr(atne *ast.TableNameExpr) *TableNameExpr {
	if atne == nil {
		return nil
	}

	// Name *ColumnName
	tn := new(TableName)
	tne.Name = tn.ParseTableName(atne.Name)

	return tne
}

// ColumnNameExpr like *ast.ColumnNameExpr
type ColumnNameExpr struct {
	// Name is the referenced column name.
	Name *ColumnName
	// Refer is the result field the column name refers to.
	// The value of Refer.Expr is used as the value of the expression.
	Refer *ResultField
}

// ParseColumnNameExpr 解析Column Name表达式
func (cne *ColumnNameExpr) ParseColumnNameExpr(acne *ast.ColumnNameExpr) *ColumnNameExpr {
	if acne == nil {
		return nil
	}

	// Name *ColumnName
	cn := new(ColumnName)
	cne.Name = cn.ParseColumnName(acne.Name)

	// Refer *ResultField
	rf := new(ResultField)
	cne.Refer = rf.ParseResultField(acne.Refer)

	return cne
}

// CompareSubqueryExpr like *ast.CompareSubqueryExpr
type CompareSubqueryExpr struct {
	// L is the left expression
	Left interface{} // ExprNode
	// Op is the comparison opcode.
	Opcode string // opcode.Op
	// R is the subquery for right expression, may be rewritten to other type of expression.
	Right interface{} // ExprNode
	// All is true, we should compare all records in subquery.
	All bool
}

// ParseCompareSubqueryExpr 解析Compare Subquery表达式
func (csqe *CompareSubqueryExpr) ParseCompareSubqueryExpr(acsqe *ast.CompareSubqueryExpr) *CompareSubqueryExpr {
	if acsqe == nil {
		return nil
	}

	// L ExprNode
	if val, bv := classifyExprNode(acsqe.L); bv == true {
		csqe.Left = val
	}

	// Op opcode.Op
	if val, bv := classifyOpcode(acsqe.Op); bv == true {
		csqe.Opcode = val
	}

	// R ExprNode
	if val, bv := classifyExprNode(acsqe.R); bv == true {
		csqe.Right = val
	}

	// All bool
	csqe.All = acsqe.All

	return csqe
}

// DefaultExpr like *ast.DefaultExpr
type DefaultExpr struct {
	// Name is the column name.
	Name *ColumnName
}

// ParseDefaultExpr 解析Default表达式
func (de *DefaultExpr) ParseDefaultExpr(ade *ast.DefaultExpr) *DefaultExpr {
	if ade == nil {
		return nil
	}

	// Name *ColumnName
	cn := new(ColumnName)
	de.Name = cn.ParseColumnName(ade.Name)

	return de
}

// IsNullExpr like *ast.IsNullExpr
type IsNullExpr struct {
	// Expr is the expression to be checked.
	Expr interface{} // ExprNode
	// Not is true, the expression is "is not null".
	Not bool
}

// ParseIsNullExpr 解析Is Null表达式
func (ine *IsNullExpr) ParseIsNullExpr(aine *ast.IsNullExpr) *IsNullExpr {
	if aine == nil {
		return nil
	}

	// Expr ExprNode
	if val, bv := classifyExprNode(aine.Expr); bv == true {
		ine.Expr = val
	}

	// Not bool
	ine.Not = aine.Not

	return ine
}

// IsTruthExpr like *ast.IsTruthExpr
type IsTruthExpr struct {
	// Expr is the expression to be checked.
	Expr interface{} // ExprNode
	// Not is true, the expression is "is not true/false".
	Not bool
	// True indicates checking true or false.
	True int64
}

// ParseIsTruthExpr 解析Is Truth表达式
func (ite *IsTruthExpr) ParseIsTruthExpr(aite *ast.IsTruthExpr) *IsTruthExpr {
	if aite == nil {
		return nil
	}

	// Expr ExprNode
	if val, bv := classifyExprNode(aite.Expr); bv == true {
		ite.Expr = val
	}

	// Not bool
	ite.Not = aite.Not

	// True int64
	ite.True = aite.True

	return ite
}

// ParenthesesExpr like *ast.ParenthesesExpr
type ParenthesesExpr struct {
	// Expr is the expression in parentheses.
	Expr interface{} // ExprNode
}

// ParseParenthesesExpr 解析Parentheses表达式
func (pte *ParenthesesExpr) ParseParenthesesExpr(apte *ast.ParenthesesExpr) *ParenthesesExpr {
	if apte == nil {
		return nil
	}

	// Expr ExprNode
	if val, bv := classifyExprNode(apte.Expr); bv == true {
		pte.Expr = val
	}

	return pte
}

// PatternInExpr like *ast.PatternInExpr
type PatternInExpr struct {
	// Expr is the value expression to be compared.
	Expr interface{} // ExprNode
	// List is the list expression in compare list.
	List []interface{} // []ExprNode
	// Not is true, the expression is "not in".
	Not bool
	// Sel is the subquery, may be rewritten to other type of expression.
	Sel interface{} // ExprNode
}

// ParsePatternInExpr 解析PatternIn表达式
func (pie *PatternInExpr) ParsePatternInExpr(apie *ast.PatternInExpr) *PatternInExpr {
	if apie == nil {
		return nil
	}

	// Expr ExprNode
	if val, bv := classifyExprNode(apie.Expr); bv == true {
		pie.Expr = val
	}

	// List []ExprNode
	for _, al := range apie.List {
		if val, bv := classifyExprNode(al); bv == true {
			alitem := val

			pie.List = append(pie.List, alitem)
		}
	}

	// Not bool
	pie.Not = apie.Not

	// Sel ExprNode
	if val, bv := classifyExprNode(apie.Sel); bv == true {
		pie.Sel = val
	}

	return pie
}

// PatternLikeExpr like *ast.PatternLikeExpr
type PatternLikeExpr struct {
	// Expr is the expression to be checked.
	Expr interface{} // ExprNode
	// Pattern is the like expression.
	Pattern interface{} // ExprNode
	// Not is true, the expression is "not like".
	Not      bool
	Escape   byte
	PatChars []byte
	PatTypes []byte
}

// ParsePatternLikeExpr 解析PatternLike表达式
func (ple *PatternLikeExpr) ParsePatternLikeExpr(aple *ast.PatternLikeExpr) *PatternLikeExpr {
	if aple == nil {
		return nil
	}

	// Expr ExprNode
	if val, bv := classifyExprNode(aple.Expr); bv == true {
		ple.Expr = val
	}

	// Pattern ExprNode
	if val, bv := classifyExprNode(aple.Pattern); bv == true {
		ple.Pattern = val
	}

	// Not      bool
	ple.Not = aple.Not

	// Escape   byte
	ple.Escape = aple.Escape

	// PatChars []byte
	ple.PatChars = aple.PatChars

	// PatTypes []byte
	ple.PatTypes = aple.PatTypes

	return ple
}

// PatternRegexpExpr like *ast.PatternRegexpExpr
type PatternRegexpExpr struct {
	// Expr is the expression to be checked.
	Expr interface{} // ExprNode
	// Pattern is the expression for pattern.
	Pattern interface{} // ExprNode
	// Not is true, the expression is "not rlike",
	Not bool
	// Re is the compiled regexp.
	Regexp *regexp.Regexp
	// Sexpr is the string for Expr expression.
	Stringexpr *string
}

// ParsePatternRegexpExpr 解析PatternRegexp表达式
func (pre *PatternRegexpExpr) ParsePatternRegexpExpr(apre *ast.PatternRegexpExpr) *PatternRegexpExpr {
	if apre == nil {
		return nil
	}

	// Expr ExprNode
	if val, bv := classifyExprNode(apre.Expr); bv == true {
		pre.Expr = val
	}

	// Pattern ExprNode
	if val, bv := classifyExprNode(apre.Pattern); bv == true {
		pre.Pattern = val
	}

	// Not      bool
	pre.Not = apre.Not

	// Re *regexp.Regexp
	pre.Regexp = apre.Re

	// Sexpr *string
	pre.Stringexpr = apre.Sexpr

	return pre
}

// VariableExpr like *ast.VariableExpr
type VariableExpr struct {
	// Name is the variable name.
	Name string
	// IsGlobal indicates whether this variable is global.
	IsGlobal bool
	// IsSystem indicates whether this variable is a system variable in current session.
	IsSystem bool
	// ExplicitScope indicates whether this variable scope is set explicitly.
	ExplicitScope bool
	// Value is the variable value.
	Value interface{} // ExprNode
}

// ParseVariableExpr 解析Variable表达式
func (ve *VariableExpr) ParseVariableExpr(ave *ast.VariableExpr) *VariableExpr {
	if ave == nil {
		return nil
	}

	// Name string
	ve.Name = ave.Name

	// IsGlobal bool
	ve.IsGlobal = ave.IsGlobal

	// IsSystem bool
	ve.IsSystem = ave.IsSystem

	// ExplicitScope bool
	ve.ExplicitScope = ave.ExplicitScope

	// Value ExprNode
	if val, bv := classifyExprNode(ave.Value); bv == true {
		ve.Value = val
	}

	return ve
}

// ValuesExpr like *ast.ValuesExpr
type ValuesExpr struct {
	// Column is column name.
	Column *ColumnNameExpr
}

// ParseValuesExpr 解析Values表达式
func (vse *ValuesExpr) ParseValuesExpr(avse *ast.ValuesExpr) *ValuesExpr {
	if avse == nil {
		return nil
	}

	// Column *ColumnNameExpr
	cne := new(ColumnNameExpr)
	cne.ParseColumnNameExpr(avse.Column)
	vse.Column = cne

	return vse
}

// TimeUnitExpr like *ast.TimeUnitExpr
type TimeUnitExpr struct {
	// Unit is the time or timestamp unit.
	Unit string // TimeUnitType
}

// ParseTimeUnitExpr 解析Time Unit表达式
func (tue *TimeUnitExpr) ParseTimeUnitExpr(atue *ast.TimeUnitExpr) *TimeUnitExpr {
	if atue == nil {
		return nil
	}

	// Unit TimeUnitType
	if val, bv := classifyTimeUnitType(atue.Unit); bv == true {
		tue.Unit = val
	}

	return tue
}

// TrimDirectionExpr like *ast.TrimDirectionExpr
type TrimDirectionExpr struct {
	// Direction is the trim direction
	Direction string // TrimDirectionType
}

// ParseTrimDirectionExpr 解析Trim Direction表达式
func (tde *TrimDirectionExpr) ParseTrimDirectionExpr(atde *ast.TrimDirectionExpr) *TrimDirectionExpr {
	if atde == nil {
		return nil
	}

	// Direction TrimDirectionType
	if val, bv := classifyTrimDirectionType(atde.Direction); bv == true {
		tde.Direction = val
	}

	return tde
}

// UnaryOperationExpr like *ast.UnaryOperationExpr
type UnaryOperationExpr struct {
	// Op is the operator opcode.
	Opcode string // opcode.Op
	// V is the unary expression.
	V interface{} // ExprNode
}

// ParseUnaryOperationExpr 解析Unary Operation表达式
func (uoe *UnaryOperationExpr) ParseUnaryOperationExpr(auoe *ast.UnaryOperationExpr) *UnaryOperationExpr {
	if auoe == nil {
		return nil
	}

	// Op  opcode.Op
	if val, bv := classifyOpcode(auoe.Op); bv == true {
		uoe.Opcode = val
	}

	// V   ExprNode
	if val, bv := classifyExprNode(auoe.V); bv == true {
		uoe.V = val
	}

	return uoe
}

// PositionExpr like *ast.PositionExpr
type PositionExpr struct {
	// N is the position, started from 1 now.
	Position int
	// P is the parameterized position.
	ParaPosition interface{} // ExprNode
	// Refer is the result field the position refers to.
	Refer *ResultField
}

// ParsePositionExpr 解析Position表达式
func (pe *PositionExpr) ParsePositionExpr(ape *ast.PositionExpr) *PositionExpr {
	if ape == nil {
		return nil
	}

	// N int
	pe.Position = ape.N

	// P ExprNode
	if val, bv := classifyExprNode(ape.P); bv == true {
		pe.ParaPosition = val
	}

	// Refer *ResultField
	rf := new(ResultField)
	pe.Refer = rf.ParseResultField(ape.Refer)

	return pe
}

// RowExpr like *ast.RowExpr
type RowExpr struct {
	Values []interface{} // []ExprNode
}

// ParseRowExpr 解析Row表达式
func (re *RowExpr) ParseRowExpr(are *ast.RowExpr) *RowExpr {
	if are == nil {
		return nil
	}

	// Values []ExprNode
	for _, av := range are.Values {
		if val, bv := classifyExprNode(av); bv == true {
			reitem := val

			re.Values = append(re.Values, reitem)
		}
	}

	return re
}

// SetCollationExpr like *ast.SetCollationExpr
type SetCollationExpr struct {
	// Expr is the expression to be set.
	Expr interface{} // ExprNode
	// Collate is the name of collation to set.
	Collate string
}

// ParseSetCollationExpr 解析Set Collation表达式
func (sce *SetCollationExpr) ParseSetCollationExpr(asce *ast.SetCollationExpr) *SetCollationExpr {
	if asce == nil {
		return nil
	}

	// Expr ExprNode
	if val, bv := classifyExprNode(asce.Expr); bv == true {
		sce.Expr = val
	}

	// Collate string
	sce.Collate = asce.Collate

	return sce
}

// WindowFuncExpr like *ast.WindowFuncExpr
type WindowFuncExpr struct {
	// F is the function name.
	Function string
	// Args is the function args.
	Args []interface{} // []ExprNode
	// Distinct cannot be true for most window functions, except `max` and `min`.
	// We need to raise error if it is not allowed to be true.
	Distinct bool
	// IgnoreNull indicates how to handle null value.
	// MySQL only supports `RESPECT NULLS`, so we need to raise error if it is true.
	IgnoreNull bool
	// FromLast indicates the calculation direction of this window function.
	// MySQL only supports calculation from first, so we need to raise error if it is true.
	FromLast bool
	// Spec is the specification of this window.
	Spec *WindowSpec
}

// ParseWindowFuncExpr 解析Window Function表达式
func (wfe *WindowFuncExpr) ParseWindowFuncExpr(awfe *ast.WindowFuncExpr) *WindowFuncExpr {
	if awfe == nil {
		return nil
	}
	// F string
	wfe.Function = awfe.F

	// Args []ExprNode
	for _, aa := range awfe.Args {
		if val, bv := classifyExprNode(aa); bv == true {
			aaitem := val

			wfe.Args = append(wfe.Args, aaitem)
		}
	}

	// Distinct bool
	wfe.Distinct = awfe.Distinct

	// IgnoreNull bool
	wfe.IgnoreNull = awfe.IgnoreNull

	// FromLast bool
	wfe.FromLast = awfe.FromLast

	// Spec WindowSpec
	ws := new(WindowSpec)
	wfe.Spec = ws.ParseWindowSpec(&awfe.Spec)

	return wfe
}
