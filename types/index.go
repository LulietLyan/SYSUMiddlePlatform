package types

import (
	"strconv"

	"github.com/pingcap/parser/ast"
)

// IndexOption like *ast.IndexOption
type IndexOption struct {
	KeyBlockSize string // uint64
	Type         string // model.IndexType
	Comment      string
	ParserName   map[string]string // model.CIStr
	Visibility   string            // IndexVisibility
}

// ParseIndexOption 解析索引选项
func (io *IndexOption) ParseIndexOption(aio *ast.IndexOption) *IndexOption {
	if aio == nil {
		return nil
	}

	// KeyBlockSize uint64
	io.KeyBlockSize = strconv.FormatUint(aio.KeyBlockSize, 10)

	// Tp model.IndexType
	if val, bv := classifyIndexType(aio.Tp); bv == true {
		io.Type = val
	}

	// Comment string
	io.Comment = aio.Comment

	// ParserName model.CIStr
	io.ParserName = beautifyModelCIStr(&aio.ParserName)

	// Visibility IndexVisibility
	if val, bv := classifyIndexVisibility(aio.Visibility); bv == true {
		io.Visibility = val
	}
	return io
}

// IndexPartSpecification like *ast.IndexPartSpecification
type IndexPartSpecification struct {
	Column *ColumnName
	Length int
	Expr   interface{} // ExprNode
}

// ParseIndexPartSpecification 解析索引分区内容
func (ips *IndexPartSpecification) ParseIndexPartSpecification(aips *ast.IndexPartSpecification) *IndexPartSpecification {
	if aips == nil {
		return nil
	}

	// Column *ColumnName
	cn := new(ColumnName)
	ips.Column = cn.ParseColumnName(aips.Column)

	// Length int
	ips.Length = aips.Length

	// Expr ExprNode
	if val, bv := classifyExprNode(aips.Expr); bv == true {
		ips.Expr = val
	}

	return ips
}

// IndexLockAndAlgorithm like *ast.IndexLockAndAlgorithm
type IndexLockAndAlgorithm struct {
	LockType      string // LockType
	AlgorithmType string // AlgorithmType
}

// ParseIndexLockAndAlgorithm 解析索引Lock And Algorithm内容
func (ilaa *IndexLockAndAlgorithm) ParseIndexLockAndAlgorithm(ailaa *ast.IndexLockAndAlgorithm) *IndexLockAndAlgorithm {
	if ailaa == nil {
		return nil
	}

	// LockTp LockType
	if val, bv := classifyLockType(ailaa.LockTp); bv == true {
		ilaa.LockType = val
	}

	// AlgorithmTp AlgorithmType
	if val, bv := classifyAlgorithmType(ailaa.AlgorithmTp); bv == true {
		ilaa.AlgorithmType = val
	}

	return ilaa
}

// IndexHint like *ast.IndexHint
type IndexHint struct {
	IndexNames []map[string]string // model.CIStr
	HintType   string              // IndexHintType
	HintScope  string              // IndexHintScope
}

// ParseIndexHint 解析Index Hint
func (ih *IndexHint) ParseIndexHint(aih *ast.IndexHint) *IndexHint {
	if aih == nil {
		return nil
	}

	// IndexNames []model.CIStr
	for _, ain := range aih.IndexNames {
		initem := beautifyModelCIStr(&ain)

		ih.IndexNames = append(ih.IndexNames, initem)
	}

	// HintType   IndexHintType
	if val, bv := classifyIndexHintType(aih.HintType); bv == true {
		ih.HintType = val
	}

	// HintScope  IndexHintScope
	if val, bv := classifyIndexHintScope(aih.HintScope); bv == true {
		ih.HintScope = val
	}

	return ih
}
