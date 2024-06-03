package types

import (
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/tidb/types"
)

// FieldType 字段属性
type FieldType struct {
	Datatype string
	Flag     int
	Fieldlen int
	Decimal  int
	Charset  string
	Collate  string
	// Elems is the element list for enum and set type.
	Elems []string
}

// ParseFieldType 解析字段属性
func (ft *FieldType) ParseFieldType(tft *types.FieldType) *FieldType {
	if tft == nil {
		return nil
	}

	// Tp      byte 数据类型
	if val, bv := classifyColDataType(tft.Tp); bv == true {
		ft.Datatype = val
	}

	// Flag    uint
	ft.Flag = int(tft.Flag)

	// Flen    int 字段长度
	ft.Fieldlen = tft.Flen

	// Decimal int 字段精度，-1代表精度为0
	ft.Decimal = tft.Decimal

	// Charset string 字符集
	ft.Charset = tft.Charset

	// Collate string 校验字符集
	ft.Collate = tft.Collate

	// Elems []string
	ft.Elems = tft.Elems

	return ft
}

// FieldList like *ast.FieldList
type FieldList struct {
	Fields []*SelectField
}

// ParseFieldList 解析Field列表
func (fl *FieldList) ParseFieldList(afl *ast.FieldList) *FieldList {
	if afl == nil {
		return nil
	}

	// Fields []*SelectField
	for _, asf := range afl.Fields {
		sf := new(SelectField)
		sfitem := sf.ParseSelectField(asf)

		fl.Fields = append(fl.Fields, sfitem)
	}

	return fl
}

// SelectField like *ast.SelectField
type SelectField struct {
	// Offset is used to get original text.
	Offset int
	// WildCard is not nil, Expr will be nil.
	WildCard *WildCardField
	// Expr is not nil, WildCard will be nil.
	Expr interface{} // ExprNode
	// AsName is alias name for Expr.
	AsName map[string]string // model.CIStr
	// Auxiliary stands for if this field is auxiliary.
	// When we add a Field into SelectField list which is used for having/orderby clause but the field is not in select clause,
	// we should set its Auxiliary to true. Then the TrimExec will trim the field.
	Auxiliary bool
}

// ParseSelectField 解析Select Field
func (sf *SelectField) ParseSelectField(asf *ast.SelectField) *SelectField {
	if asf == nil {
		return nil
	}

	// Offset int
	sf.Offset = asf.Offset

	// WildCard *WildCardField
	wcf := new(WildCardField)
	sf.WildCard = wcf.ParseWildCardField(asf.WildCard)

	// Expr ExprNode
	if val, bv := classifyExprNode(asf.Expr); bv == true {
		sf.Expr = val
	}

	// AsName model.CIStr
	sf.AsName = beautifyModelCIStr(&asf.AsName)

	// Auxiliary bool
	sf.Auxiliary = asf.Auxiliary

	return sf
}

// WildCardField like *ast.WildCardField
type WildCardField struct {
	Schema map[string]string //model.CIStr
	Table  map[string]string // model.CIStr
}

// ParseWildCardField 解析WildCard Field
func (wcf *WildCardField) ParseWildCardField(awcf *ast.WildCardField) *WildCardField {
	if awcf == nil {
		return nil
	}

	// Schema model.CIStr
	wcf.Schema = beautifyModelCIStr(&awcf.Schema)

	// Table model.CIStr
	wcf.Table = beautifyModelCIStr(&awcf.Table)

	return wcf
}

// ResultField like *ast.ResultField
type ResultField struct {
	//Column       *model.ColumnInfo
	ColumnAsName map[string]string // model.CIStr
	//Table        *model.TableInfo
	TableAsName map[string]string // model.CIStr
	DBName      map[string]string // model.CIStr
	// Expr represents the expression for the result field. If it is generated from a select field, it would
	// be the expression of that select field, otherwise the type would be ValueExpr and value
	// will be set for every retrieved row.
	Expr      interface{} // ExprNode
	TableName *TableName
	// Referenced indicates the result field has been referenced or not.
	// If not, we don't need to get the values.
	Referenced bool
}

// ParseResultField 解析Result Field
func (rf *ResultField) ParseResultField(arf *ast.ResultField) *ResultField {
	if arf == nil {
		return nil
	}

	// ColumnAsName model.CIStr
	rf.ColumnAsName = beautifyModelCIStr(&arf.ColumnAsName)

	// TableAsName model.CIStr
	rf.TableAsName = beautifyModelCIStr(&arf.TableAsName)

	// DBName      model.CIStr
	rf.DBName = beautifyModelCIStr(&arf.DBName)

	// Expr      ExprNode
	if val, bv := classifyExprNode(arf.Expr); bv == true {
		rf.Expr = val
	}

	// TableName *TableName
	tn := new(TableName)
	rf.TableName = tn.ParseTableName(arf.TableName)

	// Referenced bool
	rf.Referenced = arf.Referenced

	return rf
}
