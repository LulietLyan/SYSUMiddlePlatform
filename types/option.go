package types

import "github.com/pingcap/parser/ast"

// SelectIntoOption like *ast.SelectIntoOption
type SelectIntoOption struct {
	Type       string // SelectIntoType
	FileName   string
	FieldsInfo *FieldsClause
	LinesInfo  *LinesClause
}

// ParseSelectIntoOption 解析Select Into选项
func (sio *SelectIntoOption) ParseSelectIntoOption(asio *ast.SelectIntoOption) *SelectIntoOption {
	if asio == nil {
		return nil
	}

	// Tp         SelectIntoType
	if val, bv := classifySelectIntoType(asio.Tp); bv == true {
		sio.Type = val
	}

	// FileName   string
	sio.FileName = asio.FileName

	// FieldsInfo *FieldsClause
	fc := new(FieldsClause)
	sio.FieldsInfo = fc.ParseFieldsClause(asio.FieldsInfo)

	// LinesInfo  *LinesClause
	lc := new(LinesClause)
	sio.LinesInfo = lc.ParseLinesClause(asio.LinesInfo)

	return sio
}

// SequenceOption like *ast.SequenceOption
type SequenceOption struct {
	Type     string // SequenceOptionType
	IntValue int64
}

// ParseSequenceOption 解析Sequence选项
func (seqo *SequenceOption) ParseSequenceOption(aseqo *ast.SequenceOption) *SequenceOption {
	// Tp     SequenceOptionType
	if val, bv := classifySequenceOptionType(aseqo.Tp); bv == true {
		seqo.Type = val
	}

	// IntValue int64
	seqo.IntValue = aseqo.IntValue

	return seqo
}

// AuthOption like *ast.AuthOption
type AuthOption struct {
	// ByAuthString set as true, if AuthString is used for authorization. Otherwise, authorization is done by HashString.
	ByAuthString bool
	AuthString   string
	HashString   string
}

// ParseAuthOption 解析Auth Option
func (ao *AuthOption) ParseAuthOption(aao *ast.AuthOption) *AuthOption {
	if aao == nil {
		return nil
	}

	ao.ByAuthString = aao.ByAuthString
	ao.AuthString = aao.AuthString
	ao.HashString = aao.HashString

	return ao
}

// TLSOption like *ast.TLSOption
type TLSOption struct {
	Type  int
	Value string
}

// ParseTLSOption 解析TLS选项
func (to *TLSOption) ParseTLSOption(ato *ast.TLSOption) *TLSOption {
	if ato == nil {
		return nil
	}

	to.Type = ato.Type
	to.Value = ato.Value

	return to
}

// ResourceOption like *ast.ResourceOption
type ResourceOption struct {
	Type  int
	Count int64
}

// ParseResourceOption 解析Resource选项
func (ro *ResourceOption) ParseResourceOption(aro *ast.ResourceOption) *ResourceOption {
	if aro == nil {
		return nil
	}

	ro.Type = aro.Type
	ro.Count = aro.Count

	return ro
}

// PasswordOrLockOption like *ast.PasswordOrLockOption
type PasswordOrLockOption struct {
	Type  int
	Count int64
}

// ParsePasswordOrLockOption 解析Password Or Lock选项
func (plo *PasswordOrLockOption) ParsePasswordOrLockOption(aplo *ast.PasswordOrLockOption) *PasswordOrLockOption {
	if aplo == nil {
		return nil
	}

	plo.Type = aplo.Type
	plo.Count = aplo.Count

	return plo
}

// SplitSyntaxOption like *ast.SplitSyntaxOption
type SplitSyntaxOption struct {
	HasRegionFor bool
	HasPartition bool
}

// ParseSplitSyntaxOption 解析Split Syntax选项
func (sso *SplitSyntaxOption) ParseSplitSyntaxOption(asso *ast.SplitSyntaxOption) *SplitSyntaxOption {
	if asso == nil {
		return nil
	}

	// HasRegionFor bool
	sso.HasRegionFor = asso.HasRegionFor

	// HasPartition bool
	sso.HasPartition = asso.HasPartition

	return sso
}

// SplitOption like *ast.SplitOption
type SplitOption struct {
	Lower      []interface{} // []ExprNode
	Upper      []interface{} // []ExprNode
	Num        int64
	ValueLists [][]interface{} // [][]ExprNode
}

// ParseSplitOption 解析Split选项
func (so *SplitOption) ParseSplitOption(aso *ast.SplitOption) *SplitOption {
	if aso == nil {
		return nil
	}

	// Lower      []ExprNode
	for _, al := range aso.Lower {
		if val, bv := classifyExprNode(al); bv == true {
			litem := val

			so.Lower = append(so.Lower, litem)
		}
	}

	// Upper      []ExprNode
	for _, al := range aso.Upper {
		if val, bv := classifyExprNode(al); bv == true {
			ritem := val

			so.Upper = append(so.Upper, ritem)
		}
	}

	// Num        int64
	so.Num = aso.Num

	// ValueLists [][]ExprNode
	for _, avl := range aso.ValueLists {
		var mid []interface{}
		for _, ae := range avl {
			if val, bv := classifyExprNode(ae); bv == true {
				eitem := val

				mid = append(mid, eitem)
			}
		}

		so.ValueLists = append(so.ValueLists, mid)
	}

	return so
}

// SelectStmtOpts like *ast.SelectStmtOpts
type SelectStmtOpts struct {
	Distinct        bool
	SQLBigResult    bool
	SQLBufferResult bool
	SQLCache        bool
	SQLSmallResult  bool
	CalcFoundRows   bool
	StraightJoin    bool
	Priority        string // mysql.PriorityEnum
	TableHints      []*TableOptimizerHint
	ExplicitAll     bool
}

// ParseSelectStmtOpts 解析SELECT语句选项
func (sso *SelectStmtOpts) ParseSelectStmtOpts(asso *ast.SelectStmtOpts) *SelectStmtOpts {
	if asso == nil {
		return nil
	}

	// Distinct        bool
	sso.Distinct = asso.Distinct

	// SQLBigResult    bool
	sso.SQLBigResult = asso.SQLBigResult

	// SQLBufferResult bool
	sso.SQLBufferResult = asso.SQLBufferResult

	// SQLCache        bool
	sso.SQLCache = asso.SQLCache

	// SQLSmallResult  bool
	sso.SQLSmallResult = asso.SQLSmallResult

	// CalcFoundRows   bool
	sso.CalcFoundRows = asso.CalcFoundRows

	// StraightJoin    bool
	sso.StraightJoin = asso.StraightJoin

	// Priority        mysql.PriorityEnum
	if val, bv := classifyPriorityEnum(asso.Priority); bv == true {
		sso.Priority = val
	}

	// TableHints      []*TableOptimizerHint
	for _, ath := range asso.TableHints {
		toh := new(TableOptimizerHint)
		tohitem := toh.ParseTableOptimizerHint(ath)

		sso.TableHints = append(sso.TableHints, tohitem)
	}

	// ExplicitAll     bool
	//sso.ExplicitAll = asso.ExplicitAll

	return sso
}

// OnDeleteOpt like *ast.OnDeleteOpt
type OnDeleteOpt struct {
	ReferOpt string // ReferOptionType
}

// ParseOnDeleteOpt 解析OnDelete Option
func (odo *OnDeleteOpt) ParseOnDeleteOpt(aodo *ast.OnDeleteOpt) *OnDeleteOpt {
	if aodo == nil {
		return nil
	}

	// ReferOpt ReferOptionType
	if val, bv := classifyReferOptionType(aodo.ReferOpt); bv == true {
		odo.ReferOpt = val
	}

	return odo
}

// OnUpdateOpt like *ast.OnUpdateOpt
type OnUpdateOpt struct {
	ReferOpt string // ReferOptionType
}

// ParseOnUpdateOpt 解析OnUpdate Option
func (ouo *OnUpdateOpt) ParseOnUpdateOpt(aouo *ast.OnUpdateOpt) *OnUpdateOpt {
	if aouo == nil {
		return nil
	}

	// ReferOpt ReferOptionType
	if val, bv := classifyReferOptionType(aouo.ReferOpt); bv == true {
		ouo.ReferOpt = val
	}

	return ouo
}
