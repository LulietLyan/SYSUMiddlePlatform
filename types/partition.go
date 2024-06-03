package types

import (
	"strconv"

	"github.com/pingcap/parser/ast"
)

// PartitionOptions 分区选项
type PartitionOptions struct {
	Sub         *PartitionMethod
	Definitions []*PartitionDefinition
}

// ParsePartitionOptions 解析分区选项
func (paro *PartitionOptions) ParsePartitionOptions(aparo *ast.PartitionOptions) *PartitionOptions {
	if aparo == nil {
		return nil
	}

	// Sub         *PartitionMethod
	pm := new(PartitionMethod)
	paro.Sub = pm.ParsePartitionMethod(aparo.Sub)

	// Definitions []*PartitionDefinition
	for _, apard := range aparo.Definitions {
		pard := new(PartitionDefinition)
		parditem := pard.ParsePartitionDefinition(apard)

		paro.Definitions = append(paro.Definitions, parditem)
	}

	return paro
}

// PartitionMethod like *ast.PartitionMethod
type PartitionMethod struct {
	// Tp is the type of the partition function
	Type string // model.PartitionType
	// Linear is a modifier to the HASH and KEY type for choosing a different
	// algorithm
	Linear bool
	// Expr is an expression used as argument of HASH, RANGE, LIST and
	// SYSTEM_TIME types
	Expr interface{} // ExprNode
	// ColumnNames is a list of column names used as argument of KEY,
	// RANGE COLUMNS and LIST COLUMNS types
	ColumnNames []*ColumnName
	// Unit is a time unit used as argument of SYSTEM_TIME type
	Unit string // TimeUnitType
	// Limit is a row count used as argument of the SYSTEM_TIME type
	Limit string // uint64
	// Num is the number of (sub)partitions required by the method.
	Num string // uint64
}

// ParsePartitionMethod 解析Partition Method
func (pm *PartitionMethod) ParsePartitionMethod(apm *ast.PartitionMethod) *PartitionMethod {
	if apm == nil {
		return nil
	}

	// Type model.PartitionType
	// Linear bool
	pm.Linear = apm.Linear

	// Expr ExprNode
	if val, bv := classifyExprNode(apm.Expr); bv == true {
		pm.Expr = val
	}

	// ColumnNames []*ColumnName
	for _, acn := range apm.ColumnNames {
		cn := new(ColumnName)
		cnitem := cn.ParseColumnName(acn)

		pm.ColumnNames = append(pm.ColumnNames, cnitem)
	}

	// Unit TimeUnitType
	if val, bv := classifyTimeUnitType(apm.Unit); bv == true {
		pm.Unit = val
	}

	// Limit uint64
	pm.Limit = strconv.FormatUint(apm.Limit, 10)

	// Num uint64
	pm.Num = strconv.FormatUint(apm.Num, 10)

	return pm
}

// PartitionDefinition 分区定义
type PartitionDefinition struct {
	Name    map[string]string // model.CIStr
	Clause  ast.PartitionDefinitionClause
	Options []*TableOption
	Sub     []*SubPartitionDefinition
}

// ParsePartitionDefinition 解析分区定义
func (pard *PartitionDefinition) ParsePartitionDefinition(apard *ast.PartitionDefinition) *PartitionDefinition {
	if apard == nil {
		return nil
	}

	// Name model.CIStr
	pard.Name = beautifyModelCIStr(&apard.Name)

	// Clause  ast.PartitionDefinitionClause
	pard.Clause = apard.Clause

	// Options []*TableOption
	for _, ao := range apard.Options {
		to := new(TableOption)
		toitem := to.ParseTableOption(ao)

		pard.Options = append(pard.Options, toitem)
	}

	// Sub []*SubPartitionDefinition
	for _, aspard := range apard.Sub {
		spard := new(SubPartitionDefinition)
		sparditem := spard.ParseSubPartitionDefinition(aspard)

		pard.Sub = append(pard.Sub, sparditem)
	}

	return pard
}

// SubPartitionDefinition 子分区定义
type SubPartitionDefinition struct {
	Name    map[string]string // model.CIStr
	Options []*TableOption
}

// ParseSubPartitionDefinition 解析子分区定义
func (spard *SubPartitionDefinition) ParseSubPartitionDefinition(aspard *ast.SubPartitionDefinition) *SubPartitionDefinition {
	if aspard == nil {
		return nil
	}

	// Name model.CIStr
	spard.Name = beautifyModelCIStr(&aspard.Name)

	// Options []*TableOption
	for _, ao := range aspard.Options {
		to := new(TableOption)
		toitem := to.ParseTableOption(ao)

		spard.Options = append(spard.Options, toitem)
	}

	return spard
}
