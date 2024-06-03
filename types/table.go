package types

import (
	"reflect"
	"strconv"

	"gosqltree/utils"

	"github.com/pingcap/parser/ast"
)

// TableName like *ast.TableName
type TableName struct {
	Schema map[string]string // model.CIStr
	Name   map[string]string // model.CIStr
	//DBInfo      *model.DBInfo
	//TableInfo   *model.TableInfo
	IndexHints     []*IndexHint
	PartitionNames []map[string]string // []model.CIStr
	//TableSample    *ast.TableSample
}

// ParseTableName 解析表名
func (t *TableName) ParseTableName(at *ast.TableName) *TableName {
	if at == nil {
		return nil
	}

	// Schema model.CIStr
	t.Schema = beautifyModelCIStr(&at.Schema)

	// Name   model.CIStr
	t.Name = beautifyModelCIStr(&at.Name)

	//DBInfo      *model.DBInfo
	//TableInfo   *model.TableInfo

	// IndexHints     []*IndexHint
	for _, aih := range at.IndexHints {
		ih := new(IndexHint)
		ihitem := ih.ParseIndexHint(aih)

		t.IndexHints = append(t.IndexHints, ihitem)
	}

	// PartitionNames []model.CIStr
	if at.PartitionNames != nil {
		for _, par := range at.PartitionNames {
			pitem := beautifyModelCIStr(&par)
			t.PartitionNames = append(t.PartitionNames, pitem)
		}
	} else {
		t.PartitionNames = []map[string]string{}
	}

	return t
}

// TableOption like *ast.TableOption 表选项
type TableOption struct {
	Type       string // TableOptionType
	Default    bool
	Value      string // StrValue string | UintValue uint64
	TableNames []*TableName
}

// ParseTableOption 解析表选项
func (to *TableOption) ParseTableOption(ato *ast.TableOption) *TableOption {
	if ato == nil {
		return nil
	}

	// Tp         TableOptionType
	if val, bv := classifyTableOptionType(ato); bv == true {
		to.Type = val
	}

	// Default    bool
	to.Default = ato.Default

	// StrValue   string
	// UintValue  uint64
	if utils.IsBlank(reflect.ValueOf(ato.StrValue)) == true {
		to.Value = strconv.FormatUint(ato.UintValue, 10)
	} else {
		to.Value = ato.StrValue
	}

	// TableNames []*TableName
	for _, atn := range ato.TableNames {
		tn := new(TableName)
		tnitem := tn.ParseTableName(atn)

		to.TableNames = append(to.TableNames, tnitem)
	}

	return to
}

// TableRefsClause like *ast.TableRefsClause
type TableRefsClause struct {
	TableRefs *Join
}

// ParseTableRefsClause 解析表关联条件
func (trc *TableRefsClause) ParseTableRefsClause(atrc *ast.TableRefsClause) *TableRefsClause {
	if atrc == nil {
		return nil
	}

	j := new(Join)
	trc.TableRefs = j.ParseJoin(atrc.TableRefs)

	return trc
}

// TableSource like *ast.TableSource
type TableSource struct {
	// Source is the source of the data, can be a TableName,
	// a SelectStmt, a SetOprStmt, or a JoinNode.
	Source interface{} // ResultSetNode

	// AsName is the alias name of the table source.
	AsName map[string]string // model.CIStr
}

// ParseTableSource 解析ResultSetNode的TableSource
func (ts *TableSource) ParseTableSource(ats *ast.TableSource) *TableSource {
	if ats == nil {
		return nil
	}

	if val, bv := classifyResultSetNode(ats.Source); bv == true {
		ts.Source = val
	}

	ts.AsName = beautifyModelCIStr(&ats.AsName)

	return ts
}

// TableOptimizerHint like *ast.TableOptimizerHint
type TableOptimizerHint struct {
	// HintName is the name or alias of the table(s) which the hint will affect.
	// Table hints has no schema info
	// It allows only table name or alias (if table has an alias)
	HintName map[string]string // model.CIStr
	// HintData is the payload of the hint. The actual type of this field
	// is defined differently as according `HintName`. Define as following:
	//
	// Statement Execution Time Optimizer Hints
	// See https://dev.mysql.com/doc/refman/5.7/en/optimizer-hints.html#optimizer-hints-execution-time
	// - MAX_EXECUTION_TIME  => uint64
	// - MEMORY_QUOTA        => int64
	// - QUERY_TYPE          => model.CIStr
	//
	// Time Range is used to hint the time range of inspection tables
	// e.g: select /*+ time_range('','') */ * from information_schema.inspection_result.
	// - TIME_RANGE          => ast.HintTimeRange
	// - READ_FROM_STORAGE   => model.CIStr
	// - USE_TOJA            => bool
	// - NTH_PLAN            => int64
	HintData interface{}
	// QBName is the default effective query block of this hint.
	QBName  map[string]string // model.CIStr
	Tables  []*HintTable
	Indexes []map[string]string // []model.CIStr
}

// ParseTableOptimizerHint 解析表HINT
func (toh *TableOptimizerHint) ParseTableOptimizerHint(atoh *ast.TableOptimizerHint) *TableOptimizerHint {
	if toh == nil {
		return nil
	}

	// HintName model.CIStr
	toh.HintName = beautifyModelCIStr(&atoh.HintName)

	// HintData interface{}
	toh.HintData = atoh.HintData

	// QBName  model.CIStr
	toh.QBName = beautifyModelCIStr(&atoh.QBName)

	// Tables  []HintTable
	for _, aht := range atoh.Tables {
		ht := new(HintTable)
		htitem := ht.ParseHintTable(&aht)

		toh.Tables = append(toh.Tables, htitem)
	}

	// Indexes []model.CIStr
	for _, aidx := range atoh.Indexes {
		idx := beautifyModelCIStr(&aidx)
		toh.Indexes = append(toh.Indexes, idx)
	}

	return toh
}

// HintTable like *ast.HintTable
type HintTable struct {
	DBName    map[string]string // model.CIStr
	TableName map[string]string // model.CIStr
	QBName    map[string]string // model.CIStr
	// PartitionList []model.CIStr
}

// ParseHintTable 解析HINT表
func (ht *HintTable) ParseHintTable(aht *ast.HintTable) *HintTable {
	if aht == nil {
		return nil
	}

	// DBName        model.CIStr
	ht.DBName = beautifyModelCIStr(&aht.DBName)

	// TableName     model.CIStr
	ht.TableName = beautifyModelCIStr(&aht.TableName)

	// QBName        model.CIStr
	ht.QBName = beautifyModelCIStr(&aht.QBName)

	// PartitionList []model.CIStr
	//for _, ap := range ath.PartitionList {}

	return ht
}

// TableToTable like *ast.TableToTable
type TableToTable struct {
	OldTable *TableName
	NewTable *TableName
}

// ParseTableToTable 解析Table To Table
func (ttt *TableToTable) ParseTableToTable(attt *ast.TableToTable) *TableToTable {
	if attt == nil {
		return nil
	}

	// OldTable *TableName
	tn1 := new(TableName)
	ttt.OldTable = tn1.ParseTableName(attt.OldTable)

	//NewTable *TableName
	tn2 := new(TableName)
	ttt.NewTable = tn2.ParseTableName(attt.NewTable)

	return ttt
}

// TableLock like *ast.TableLock
type TableLock struct {
	Table *TableName
	Type  string // model.TableLockType
}

// ParseTableLock 解析Table Lock
func (tl *TableLock) ParseTableLock(atl *ast.TableLock) *TableLock {
	if atl == nil {
		return nil
	}

	// Table *TableName
	tn := new(TableName)
	tl.Table = tn.ParseTableName(atl.Table)

	// Type  model.TableLockType
	if val, bv := classifyTableLockType(atl.Type); bv == true {
		tl.Type = val
	}

	return tl
}
