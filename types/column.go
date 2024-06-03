package types

import "github.com/pingcap/parser/ast"

// ColumnName 字段名称
type ColumnName struct {
	Schema map[string]string
	Table  map[string]string
	Name   map[string]string
}

// ParseColumnName 解析字段名称
func (cn *ColumnName) ParseColumnName(acn *ast.ColumnName) *ColumnName {
	if acn == nil {
		return nil
	}

	// Schema model.CIStr
	cn.Schema = beautifyModelCIStr(&acn.Schema)

	// Table  model.CIStr
	cn.Table = beautifyModelCIStr(&acn.Table)

	// Name   model.CIStr
	cn.Name = beautifyModelCIStr(&acn.Name)

	return cn
}

// ColumnDef like *ast.ColumnDef 字段属性
type ColumnDef struct {
	Name    *ColumnName
	Type    *FieldType
	Options []*ColumnOption
}

// ParseColumnDef 解析字段定义
func (cd *ColumnDef) ParseColumnDef(acd *ast.ColumnDef) *ColumnDef {
	if acd == nil {
		return nil
	}

	// Name    *ColumnName
	cn := new(ColumnName)
	cd.Name = cn.ParseColumnName(acd.Name)

	// Tp      *types.FieldType
	ft := new(FieldType)
	cd.Type = ft.ParseFieldType(acd.Tp)

	// Options []*ColumnOption
	for _, aco := range acd.Options {
		co := new(ColumnOption)
		coitem := co.ParseColumnOption(aco)

		cd.Options = append(cd.Options, coitem)
	}

	return cd
}

// ColumnOption 字段选项
type ColumnOption struct {
	Type string // ColumnOptionType
	// Expr is used for ColumnOptionDefaultValue/ColumnOptionOnUpdateColumnOptionGenerated.
	// For ColumnOptionDefaultValue or ColumnOptionOnUpdate, it's the target value.
	// For ColumnOptionGenerated, it's the target expression.
	Expr interface{} // ExprNode
	// Stored is only for ColumnOptionGenerated, default is false.
	Stored bool
	// Refer is used for foreign key.
	Refer               *ReferenceDef
	StrValue            string
	AutoRandomBitLength int
	// Enforced is only for Check, default is true.
	Enforced bool
	// Name is only used for Check Constraint name.
	ConstraintName string
}

// ParseColumnOption 解析Column Option
func (co *ColumnOption) ParseColumnOption(aco *ast.ColumnOption) *ColumnOption {
	if aco == nil {
		return nil
	}

	// Tp  ColumnOptionType
	if val, bv := classifyColumnOptionType(aco); bv == true {
		co.Type = val
	}
	// Expr ExprNode
	if val, bv := classifyExprNode(aco.Expr); bv == true {
		co.Expr = val
	}

	// Stored bool
	co.Stored = aco.Stored

	// Refer  *ReferenceDef
	rd := new(ReferenceDef)
	co.Refer = rd.ParseReferenceDef(aco.Refer)

	// StrValue string
	co.StrValue = aco.StrValue

	// AutoRandomBitLength int
	co.AutoRandomBitLength = aco.AutoRandomBitLength

	// Enforced bool
	co.Enforced = aco.Enforced

	// ConstraintName string
	//co.ConstraintName = aco.ConstraintName

	return co
}

// ColumnPosition like *ast.ColumnPosition
type ColumnPosition struct {
	// Tp is either ColumnPositionNone, ColumnPositionFirst or ColumnPositionAfter.
	Type string // Tp ColumnPositionType
	// RelativeColumn is the column the newly added column after if type is ColumnPositionAfter
	RelativeColumn *ColumnName
	// contains filtered or unexported fields
}

// ParseColumnPosition 解析Column Position
func (cp *ColumnPosition) ParseColumnPosition(acp *ast.ColumnPosition) *ColumnPosition {
	if acp == nil {
		return nil
	}

	// Tp  ColumnPositionType
	if val, bv := classifyColumnPositionType(acp); bv == true {
		cp.Type = val
	}

	// RelativeColumn *ColumnName
	cn := new(ColumnName)
	cp.RelativeColumn = cn.ParseColumnName(acp.RelativeColumn)

	return cp
}
