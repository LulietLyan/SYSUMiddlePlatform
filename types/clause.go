package types

import (
	"strconv"

	"github.com/pingcap/parser/ast"
)

// WhenClause like *ast.WhenClause
type WhenClause struct {
	// Expr is the condition expression in WhenClause.
	Expr interface{} // ExprNode
	// Result is the result expression in WhenClause.
	Result interface{} // ExprNode
}

// ParseWhenClause 解析When条件
func (wc *WhenClause) ParseWhenClause(awc *ast.WhenClause) *WhenClause {
	if awc == nil {
		return nil
	}

	// Expr ExprNode
	if val, bv := classifyExprNode(awc.Expr); bv == true {
		wc.Expr = val
	}

	// Result ExprNode
	if val, bv := classifyExprNode(awc.Result); bv == true {
		wc.Result = val
	}

	return wc
}

// OrderByClause like *ast.OrderByClause
type OrderByClause struct {
	Items    []*ByItem
	ForUnion bool
}

// ParseOrderByClause 解析OrderBy条件
func (obc *OrderByClause) ParseOrderByClause(aobc *ast.OrderByClause) *OrderByClause {
	if aobc == nil {
		return nil
	}

	// Items []*ByItem
	for _, abi := range aobc.Items {
		bi := new(ByItem)
		iitem := bi.ParseByItem(abi)

		obc.Items = append(obc.Items, iitem)
	}

	// ForUnion bool
	obc.ForUnion = aobc.ForUnion

	return obc
}

// GroupByClause like *ast.GroupByClause
type GroupByClause struct {
	Items []*ByItem
}

// ParseGroupByClause 解析Group By条件
func (gc *GroupByClause) ParseGroupByClause(agc *ast.GroupByClause) *GroupByClause {
	if agc == nil {
		return nil
	}

	// Items []*ByItem
	for _, abi := range agc.Items {
		bi := new(ByItem)
		biitem := bi.ParseByItem(abi)

		gc.Items = append(gc.Items, biitem)
	}

	return gc
}

// HavingClause like *ast.HavingClause
type HavingClause struct {
	Expr interface{} // ExprNode
}

// ParseHavingClause 解析Having条件
func (hc *HavingClause) ParseHavingClause(ahc *ast.HavingClause) *HavingClause {
	if ahc == nil {
		return nil
	}

	// Expr ExprNode
	if val, bv := classifyExprNode(ahc.Expr); bv == true {
		hc.Expr = val
	}

	return hc
}

// PartitionByClause like *ast.PartitionByClause
type PartitionByClause struct {
	Items []*ByItem
}

// ParsePartitionByClause 解析分区条件
func (parc *PartitionByClause) ParsePartitionByClause(aparc *ast.PartitionByClause) *PartitionByClause {
	if aparc == nil {
		return nil
	}

	// Items []*ByItem
	for _, abi := range aparc.Items {
		bi := new(ByItem)
		biitem := bi.ParseByItem(abi)

		parc.Items = append(parc.Items, biitem)
	}

	return parc
}

// FrameClause like *ast.FrameClause
type FrameClause struct {
	Type   FrameType
	Extent *FrameExtent
}

// ParseFrameClause 解析FrameClause
func (fc *FrameClause) ParseFrameClause(afc *ast.FrameClause) *FrameClause {
	if afc == nil {
		return nil
	}

	// Type   FrameType
	fc.Type = FrameType(afc.Type)

	// Extent FrameExtent
	fe := new(FrameExtent)
	fc.Extent = fe.ParseFrameExtent(&afc.Extent)

	return fc
}

// FieldsClause like *ast.FieldsClause
type FieldsClause struct {
	Terminated  string
	Enclosed    string // byte
	Escaped     string // byte
	OptEnclosed bool
}

// ParseFieldsClause 解析FieldsClause
func (fc *FieldsClause) ParseFieldsClause(afc *ast.FieldsClause) *FieldsClause {
	if afc == nil {
		return nil
	}

	// Terminated  string
	fc.Terminated = afc.Terminated

	// Enclosed    byte
	fc.Enclosed = string(afc.Enclosed)

	// Escaped     byte
	fc.Escaped = string(afc.Escaped)

	// OptEnclosed bool
	fc.OptEnclosed = afc.OptEnclosed

	return fc
}

// LinesClause like *ast.LinesClause
type LinesClause struct {
	Starting   string
	Terminated string
}

// ParseLinesClause 解析LinesClause
func (lc *LinesClause) ParseLinesClause(alc *ast.LinesClause) *LinesClause {
	if alc == nil {
		return nil
	}

	// Starting   string
	lc.Starting = alc.Starting

	// Terminated string
	lc.Terminated = alc.Terminated

	return lc
}

// MaxIndexNumClause like *ast.MaxIndexNumClause
type MaxIndexNumClause struct {
	PerTable string // uint64
	PerDB    string // uint64
}

// ParseMaxIndexNumClause 解析MaxIndexNumClause
func (minc *MaxIndexNumClause) ParseMaxIndexNumClause(aminc *ast.MaxIndexNumClause) *MaxIndexNumClause {
	if aminc == nil {
		return nil
	}

	// PerTable uint64
	minc.PerTable = strconv.FormatUint(aminc.PerTable, 10)

	// PerDB    uint64
	minc.PerDB = strconv.FormatUint(aminc.PerDB, 10)

	return minc
}
