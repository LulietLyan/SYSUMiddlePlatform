package types

import (
	"strconv"

	"github.com/pingcap/parser/ast"
)

// SelectStmt like *ast.SelectStmt
type SelectStmt struct {
	// SelectStmtOpts wraps around select hints and switches.
	*SelectStmtOpts
	// Distinct represents whether the select has distinct option.
	Distinct bool
	// From is the from clause of the query.
	From *TableRefsClause
	// Where is the where clause in select statement.
	Where interface{} // ExprNode
	// Fields is the select expression list.
	Fields *FieldList
	// GroupBy is the group by expression list.
	GroupBy *GroupByClause
	// Having is the having condition.
	Having *HavingClause
	// WindowSpecs is the window specification list.
	WindowSpecs []WindowSpec
	// OrderBy is the ordering expression list.
	OrderBy *OrderByClause
	// Limit is the limit clause.
	Limit *Limit
	// LockInfo is the lock type
	//LockInfo *ast.SelectLockInfo
	// TableHints represents the table level Optimizer Hint for join type
	TableHints []*TableOptimizerHint
	// IsInBraces indicates whether it's a stmt in brace.
	IsInBraces bool
	// QueryBlockOffset indicates the order of this SelectStmt if counted from left to right in the sql text.
	QueryBlockOffset int
	// SelectIntoOpt is the select-into option.
	SelectIntoOpt *SelectIntoOption
	// AfterSetOperator indicates the SelectStmt after which type of set operator
	//AfterSetOperator *ast.SetOprType
	// Kind refer to three kind of statement: SelectStmt, TableStmt and ValuesStmt
	//Kind             ast.SelectStmtKind
	// Lists is filled only when Kind == SelectStmtKindValues
	//Lists []*ast.RowExpr
}

// ParseSelectStmt 解析SELECT
func (sel *SelectStmt) ParseSelectStmt(node *ast.SelectStmt) *SelectStmt {
	// *SelectStmtOpts
	sso := new(SelectStmtOpts)
	sel.SelectStmtOpts = sso.ParseSelectStmtOpts(node.SelectStmtOpts)

	// Distinct bool
	sel.Distinct = node.Distinct

	// From *TableRefsClause
	trc := new(TableRefsClause)
	sel.From = trc.ParseTableRefsClause(node.From)

	// Where ExprNode
	if val, bv := classifyExprNode(node.Where); bv == true {
		sel.Where = val
	}

	// Fields *FieldList
	fl := new(FieldList)
	sel.Fields = fl.ParseFieldList(node.Fields)

	// GroupBy *GroupByClause
	gc := new(GroupByClause)
	sel.GroupBy = gc.ParseGroupByClause(node.GroupBy)

	// Having *HavingClause
	hc := new(HavingClause)
	sel.Having = hc.ParseHavingClause(node.Having)

	// WindowSpecs []WindowSpec
	for _, aws := range node.WindowSpecs {
		ws := new(WindowSpec)
		wsitem := ws.ParseWindowSpec(&aws)

		sel.WindowSpecs = append(sel.WindowSpecs, *wsitem)
	}

	// OrderBy *OrderByClause
	oc := new(OrderByClause)
	sel.OrderBy = oc.ParseOrderByClause(node.OrderBy)

	// Limit *Limit
	l := new(Limit)
	sel.Limit = l.ParseLimit(node.Limit)

	// LockInfo *SelectLockInfo

	// TableHints []*TableOptimizerHint
	for _, atoh := range node.TableHints {
		toh := new(TableOptimizerHint)
		tohitem := toh.ParseTableOptimizerHint(atoh)

		sel.TableHints = append(sel.TableHints, tohitem)
	}
	// IsInBraces bool
	sel.IsInBraces = node.IsInBraces

	// QueryBlockOffset int
	sel.QueryBlockOffset = node.QueryBlockOffset

	// SelectIntoOpt *SelectIntoOption
	sio := new(SelectIntoOption)
	sel.SelectIntoOpt = sio.ParseSelectIntoOption(node.SelectIntoOpt)

	// AfterSetOperator *SetOprType
	// Kind SelectStmtKind
	// Lists []*RowExpr

	return sel
}

// InsertStmt like *ast.InsertStmt
type InsertStmt struct {
	IsReplace   bool
	IgnoreErr   bool
	Table       *TableRefsClause
	Columns     []*ColumnName
	Lists       [][]interface{} // [][]ExprNode
	Setlist     []*Assignment
	Priority    string // mysql.PriorityEnum
	OnDuplicate []*Assignment
	Select      interface{} // ResultSetNode
	// TableHints represents the table level Optimizer Hint for join type.
	TableHints     []*TableOptimizerHint
	PartitionNames []map[string]string // []model.CIStr
}

// ParseInsertStmt 解析INSERT
func (i *InsertStmt) ParseInsertStmt(node *ast.InsertStmt) *InsertStmt {
	// IsReplace bool
	i.IsReplace = node.IsReplace

	// IgnoreErr bool
	i.IgnoreErr = node.IgnoreErr

	// Table *TableRefsClause
	trc := new(TableRefsClause)
	i.Table = trc.ParseTableRefsClause(node.Table)

	// Columns []*ColumnName
	for _, ac := range node.Columns {
		c := new(ColumnName)

		citem := c.ParseColumnName(ac)
		i.Columns = append(i.Columns, citem)
	}

	// Lists [][]interface{}
	for _, al := range node.Lists {
		mid := []interface{}{}
		for _, alitem := range al {
			if val, bv := classifyExprNode(alitem); bv == true {
				en := val
				mid = append(mid, en)
			}
		}
		i.Lists = append(i.Lists, mid)
	}

	// Setlist []*Assignment
	for _, as := range node.Setlist {
		a := new(Assignment)

		aitem := a.ParserAssignment(as)
		i.Setlist = append(i.Setlist, aitem)
	}

	// Priority    string
	if val, bv := classifyPriorityEnum(node.Priority); bv == true {
		i.Priority = val
	}

	// OnDuplicate []*Assignment
	for _, aod := range node.OnDuplicate {
		a := new(Assignment)

		aitem := a.ParserAssignment(aod)
		i.OnDuplicate = append(i.OnDuplicate, aitem)
	}

	// Select interface{}
	if val, bv := classifyResultSetNode(node.Select); bv == true {
		i.Select = val
	}

	// TableHints []*TableOptimizerHint
	for _, ath := range node.TableHints {
		toh := new(TableOptimizerHint)
		tohitem := toh.ParseTableOptimizerHint(ath)

		i.TableHints = append(i.TableHints, tohitem)
	}

	// PartitionNames []map[string]string
	for _, aparn := range node.PartitionNames {
		parn := beautifyModelCIStr(&aparn)
		i.PartitionNames = append(i.PartitionNames, parn)
	}

	return i
}

// UpdateStmt like *ast.UpdateStmt
type UpdateStmt struct {
	TableRefs     *TableRefsClause
	List          []*Assignment
	Where         interface{} // ExprNode
	Order         *OrderByClause
	Limit         *Limit
	Priority      string // mysql.PriorityEnum
	IgnoreErr     bool
	MultipleTable bool
	TableHints    []*TableOptimizerHint
}

// ParseUpdateStmt 解析UPDATE
func (u *UpdateStmt) ParseUpdateStmt(node *ast.UpdateStmt) *UpdateStmt {
	// TableRefs *TableRefsClause
	trc := new(TableRefsClause)
	u.TableRefs = trc.ParseTableRefsClause(node.TableRefs)

	// List []*Assignment
	for _, as := range node.List {
		a := new(Assignment)

		aitem := a.ParserAssignment(as)
		u.List = append(u.List, aitem)
	}

	// Where ExprNode
	if val, bv := classifyExprNode(node.Where); bv == true {
		u.Where = val
	}

	// Order *OrderByClause
	obc := new(OrderByClause)
	u.Order = obc.ParseOrderByClause(node.Order)

	// Limit *Limit
	l := new(Limit)
	u.Limit = l.ParseLimit(node.Limit)

	// Priority string
	if val, bv := classifyPriorityEnum(node.Priority); bv == true {
		u.Priority = val
	}

	// IgnoreErr bool
	u.IgnoreErr = node.IgnoreErr

	// MultipleTable bool
	u.MultipleTable = node.MultipleTable

	// TableHints []*TableOptimizerHint
	for _, ath := range node.TableHints {
		toh := new(TableOptimizerHint)
		th := toh.ParseTableOptimizerHint(ath)

		u.TableHints = append(u.TableHints, th)
	}

	return u
}

// DeleteStmt like *ast.DeleteStmt
type DeleteStmt struct {
	// TableRefs is used in both single table and multiple table delete statement.
	TableRefs *TableRefsClause
	// Tables is only used in multiple table delete statement.
	Tables       *DeleteTableList
	Where        interface{} // ExprNode
	Order        *OrderByClause
	Limit        *Limit
	Priority     string // mysql.PriorityEnum
	IgnoreErr    bool
	Quick        bool
	IsMultiTable bool
	BeforeFrom   bool
	TableHints   []*TableOptimizerHint
}

// ParseDeleteStmt 解析DELETE
func (d *DeleteStmt) ParseDeleteStmt(node *ast.DeleteStmt) *DeleteStmt {
	// TableRefs *TableRefsClause
	trc := new(TableRefsClause)
	d.TableRefs = trc.ParseTableRefsClause(node.TableRefs)

	// Tables       *DeleteTableList
	dtl := new(DeleteTableList)
	d.Tables = dtl.ParserDeleteTableList(node.Tables)

	// Where ExprNode
	if val, bv := classifyExprNode(node.Where); bv == true {
		d.Where = val
	}

	//Order *OrderByClause
	obc := new(OrderByClause)
	d.Order = obc.ParseOrderByClause(node.Order)

	// Limit  *Limit
	l := new(Limit)
	d.Limit = l.ParseLimit(node.Limit)

	// Priority string
	if val, bv := classifyPriorityEnum(node.Priority); bv == true {
		d.Priority = val
	}

	d.IgnoreErr = node.IgnoreErr
	d.Quick = node.Quick
	d.IsMultiTable = node.IsMultiTable
	d.BeforeFrom = node.BeforeFrom

	// TableHints []*TableOptimizerHint
	for _, ath := range node.TableHints {
		toh := new(TableOptimizerHint)
		th := toh.ParseTableOptimizerHint(ath)

		d.TableHints = append(d.TableHints, th)
	}

	return d
}

// CreateDatabaseStmt like *ast.CreateDatabaseStmt
type CreateDatabaseStmt struct {
	IfNotExists bool
	Name        string
	Options     []*DatabaseOption
}

// ParseCreateDatabaseStmt 解析CREATE DATABASE
func (cd *CreateDatabaseStmt) ParseCreateDatabaseStmt(node *ast.CreateDatabaseStmt) *CreateDatabaseStmt {
	cd.IfNotExists = node.IfNotExists
	cd.Name = node.Name

	for _, ado := range node.Options {
		do := new(DatabaseOption)
		doitem := do.ParseDatabaseOption(ado)
		cd.Options = append(cd.Options, doitem)
	}
	return cd
}

// CreateUserStmt like *ast.CreateUserStmt
type CreateUserStmt struct {
	IsCreateRole          bool
	IfNotExists           bool
	Specs                 []*UserSpec
	TLSOptions            []*TLSOption
	ResourceOptions       []*ResourceOption
	PasswordOrLockOptions []*PasswordOrLockOption
}

// ParseCreateUserStmt 解析CREATE USER
func (cu *CreateUserStmt) ParseCreateUserStmt(node *ast.CreateUserStmt) *CreateUserStmt {
	cu.IsCreateRole = node.IsCreateRole
	cu.IfNotExists = node.IfNotExists

	// Specs []*UserSpec
	for _, aus := range node.Specs {
		us := new(UserSpec)
		usitem := us.ParseUserSpec(aus)

		cu.Specs = append(cu.Specs, usitem)
	}

	// TLSOptions  []*TLSOption
	for _, ato := range node.TLSOptions {
		to := new(TLSOption)
		toitem := to.ParseTLSOption(ato)

		cu.TLSOptions = append(cu.TLSOptions, toitem)
	}
	// if len(node.TLSOptions) != 0 {
	// 	for _, ato := range node.TLSOptions {
	// 		to := new(TLSOption)
	// 		toitem := to.ParseTLSOption(ato)

	// 		cu.TLSOptions = append(cu.TLSOptions, toitem)
	// 	}
	// } else {
	// 	cu.TLSOptions = []*TLSOption{}
	// }

	// ResourceOptions []*ResourceOption
	for _, aro := range node.ResourceOptions {
		ro := new(ResourceOption)
		roitem := ro.ParseResourceOption(aro)

		cu.ResourceOptions = append(cu.ResourceOptions, roitem)
	}

	// PasswordOrLockOptions []*PasswordOrLockOption
	for _, aplo := range node.PasswordOrLockOptions {
		plo := new(PasswordOrLockOption)
		ploitem := plo.ParsePasswordOrLockOption(aplo)

		cu.PasswordOrLockOptions = append(cu.PasswordOrLockOptions, ploitem)
	}

	return cu
}

// CreateTableStmt like *ast.CreateTableStmt
type CreateTableStmt struct {
	IfNotExists bool
	IsTemporary bool
	Table       *TableName
	ReferTable  *TableName
	Columns     []*ColumnDef
	Constraints []*Constraint
	Options     []*TableOption
	Partition   *PartitionOptions
	OnDuplicate string      // OnDuplicateKeyHandlingType
	Select      interface{} // ResultSetNode
}

// ParseCreateTableStmt 解析CREATE TABLE
func (ct *CreateTableStmt) ParseCreateTableStmt(node *ast.CreateTableStmt) *CreateTableStmt {
	// IfNotExists bool
	ct.IfNotExists = node.IfNotExists

	// IsTemporary bool
	ct.IsTemporary = node.IsTemporary

	// Table       *TableName
	tn1 := new(TableName)
	ct.Table = tn1.ParseTableName(node.Table)

	// ReferTable  *TableName
	tn2 := new(TableName)
	ct.ReferTable = tn2.ParseTableName(node.ReferTable)

	// Cols        []*ColumnDef
	for _, acd := range node.Cols {
		cd := new(ColumnDef)

		cditem := cd.ParseColumnDef(acd)
		ct.Columns = append(ct.Columns, cditem)
	}

	// Constraints []*Constraint
	for _, acst := range node.Constraints {
		cst := new(Constraint)

		cstitem := cst.ParseConstraint(acst)
		ct.Constraints = append(ct.Constraints, cstitem)
	}

	// Options     []*TableOption
	for _, ato := range node.Options {
		to := new(TableOption)

		toitem := to.ParseTableOption(ato)
		ct.Options = append(ct.Options, toitem)
	}

	// Partition   *PartitionOptions
	paro := new(PartitionOptions)
	ct.Partition = paro.ParsePartitionOptions(node.Partition)

	// OnDuplicate OnDuplicateKeyHandlingType
	if val, bv := classifyOnDuplicateKeyHandlingType(node.OnDuplicate); bv == true {
		ct.OnDuplicate = val
	}

	// Select      ResultSetNode
	if val, bv := classifyResultSetNode(node.Select); bv == true {
		ct.Select = val
	}

	return ct
}

// CreateIndexStmt like *ast.CreateIndexStmt
type CreateIndexStmt struct {
	// only supported by MariaDB 10.0.2+, see https://mariadb.com/kb/en/library/create-index/
	IfNotExists             bool
	IndexName               string
	Table                   *TableName
	IndexPartSpecifications []*IndexPartSpecification
	IndexOption             *IndexOption
	KeyType                 string // IndexKeyType
	LockAlg                 *IndexLockAndAlgorithm
}

// ParseCreateIndexStmt 解析CREATE INDEX
func (ci *CreateIndexStmt) ParseCreateIndexStmt(node *ast.CreateIndexStmt) *CreateIndexStmt {
	// IfNotExists bool
	ci.IfNotExists = node.IfNotExists

	// IndexName               string
	ci.IndexName = node.IndexName

	// Table  *TableName
	t := new(TableName)
	ci.Table = t.ParseTableName(node.Table)

	// IndexPartSpecifications []*IndexPartSpecification
	for _, aips := range node.IndexPartSpecifications {
		ips := new(IndexPartSpecification)
		ipsitem := ips.ParseIndexPartSpecification(aips)

		ci.IndexPartSpecifications = append(ci.IndexPartSpecifications, ipsitem)
	}

	// IndexOption *IndexOption
	io := new(IndexOption)
	ci.IndexOption = io.ParseIndexOption(node.IndexOption)

	// KeyType IndexKeyType
	if val, bv := classifyIndexKeyType(node.KeyType); bv == true {
		ci.KeyType = val
	}

	// LockAlg *IndexLockAndAlgorithm
	ilaa := new(IndexLockAndAlgorithm)
	ci.LockAlg = ilaa.ParseIndexLockAndAlgorithm(node.LockAlg)

	return ci
}

// CreateViewStmt like *ast.CreateViewStmt
type CreateViewStmt struct {
	OrReplace   bool
	ViewName    *TableName
	Cols        []map[string]string // []model.CIStr
	Select      interface{}         // StmtNode
	SchemaCols  []map[string]string // []model.CIStr
	Algorithm   string              // model.ViewAlgorithm
	Definer     *UserIdentity
	Security    string // model.ViewSecurity
	CheckOption string // model.ViewCheckOption
}

// ParseCreateViewStmt 解析CREATE VIEW
func (cv *CreateViewStmt) ParseCreateViewStmt(node *ast.CreateViewStmt) *CreateViewStmt {
	// OrReplace   bool
	cv.OrReplace = node.OrReplace

	// ViewName    *TableName
	tn := new(TableName)
	cv.ViewName = tn.ParseTableName(node.ViewName)

	// Cols        []model.CIStr
	for _, ac := range node.Cols {
		clitem := beautifyModelCIStr(&ac)

		cv.Cols = append(cv.Cols, clitem)
	}

	// Select      StmtNode
	st := new(SQLTree)
	st.ClassifyStmtNode(node.Select)
	cv.Select = st

	// SchemaCols  []model.CIStr
	for _, asc := range node.SchemaCols {
		sclitem := beautifyModelCIStr(&asc)

		cv.SchemaCols = append(cv.SchemaCols, sclitem)
	}

	// Algorithm   model.ViewAlgorithm
	if val, bv := classifyViewAlgorithm(node.Algorithm); bv == true {
		cv.Algorithm = val
	}

	// Definer     *auth.UserIdentity
	ui := new(UserIdentity)
	cv.Definer = ui.ParseUserIdentity(node.Definer)

	// Security    model.ViewSecurity
	if val, bv := classifyViewSecurity(node.Security); bv == true {
		cv.Security = val
	}

	// CheckOption model.ViewCheckOption
	if val, bv := classifyViewCheckOption(node.CheckOption); bv == true {
		cv.CheckOption = val
	}

	return cv
}

// CreateBindingStmt like *ast.CreateBindingStmt
type CreateBindingStmt struct {
	GlobalScope bool
	OriginNode  interface{} // StmtNode
	HintedNode  interface{} // StmtNode
}

// ParseCreateBindingStmt 解析CREATE BINDING
func (cb *CreateBindingStmt) ParseCreateBindingStmt(node *ast.CreateBindingStmt) *CreateBindingStmt {
	// GlobalScope bool
	cb.GlobalScope = node.GlobalScope

	// OriginNode  StmtNode
	st1 := new(SQLTree)
	st1.ClassifyStmtNode(node.OriginSel)
	cb.OriginNode = st1

	// HintedNode  StmtNode
	st2 := new(SQLTree)
	st2.ClassifyStmtNode(node.HintedSel)
	cb.HintedNode = st2

	return cb
}

// CreateSequenceStmt like *ast.CreateSequenceStmt
type CreateSequenceStmt struct {
	// TODO : support or replace if need : care for it will conflict on temporaryOpt.
	IfNotExists bool
	Name        *TableName
	SeqOptions  []*SequenceOption
	TblOptions  []*TableOption
}

// ParseCreateSequenceStmt 解析CREATE SEQUENCE
func (cs *CreateSequenceStmt) ParseCreateSequenceStmt(node *ast.CreateSequenceStmt) *CreateSequenceStmt {
	// IfNotExists bool
	cs.IfNotExists = node.IfNotExists

	// Name        *TableName
	tn := new(TableName)
	cs.Name = tn.ParseTableName(node.Name)

	// SeqOptions  []*SequenceOption
	for _, aso := range node.SeqOptions {
		seqo := new(SequenceOption)
		seqoitem := seqo.ParseSequenceOption(aso)

		cs.SeqOptions = append(cs.SeqOptions, seqoitem)
	}

	// TblOptions  []*TableOption
	for _, ato := range node.TblOptions {
		to := new(TableOption)
		toitem := to.ParseTableOption(ato)

		cs.TblOptions = append(cs.TblOptions, toitem)
	}

	return cs
}

// AlterDatabaseStmt like *ast.AlterDatabaseStmt
type AlterDatabaseStmt struct {
	Name                 string
	AlterDefaultDatabase bool
	Options              []*DatabaseOption
}

// ParseAlterDatabaseStmt 解析Alter Database
func (ad *AlterDatabaseStmt) ParseAlterDatabaseStmt(node *ast.AlterDatabaseStmt) *AlterDatabaseStmt {
	// Name string
	ad.Name = node.Name

	// AlterDefaultDatabase bool
	ad.AlterDefaultDatabase = node.AlterDefaultDatabase

	// Options []*DatabaseOption
	for _, ado := range node.Options {
		do := new(DatabaseOption)
		doitem := do.ParseDatabaseOption(ado)

		ad.Options = append(ad.Options, doitem)
	}

	return ad
}

// AlterInstanceStmt like *ast.AlterInstanceStmt
type AlterInstanceStmt struct {
	ReloadTLS         bool
	NoRollbackOnError bool
}

// ParseAlterInstanceStmt 解析Alter Instance
func (ai *AlterInstanceStmt) ParseAlterInstanceStmt(node *ast.AlterInstanceStmt) *AlterInstanceStmt {
	// ReloadTLS         bool
	ai.ReloadTLS = node.ReloadTLS

	// NoRollbackOnError bool
	ai.NoRollbackOnError = node.NoRollbackOnError

	return ai
}

// AlterUserStmt like *ast.AlterUserStmt
type AlterUserStmt struct {
	IfExists              bool
	CurrentAuth           *AuthOption
	Specs                 []*UserSpec
	TLSOptions            []*TLSOption
	ResourceOptions       []*ResourceOption
	PasswordOrLockOptions []*PasswordOrLockOption
}

// ParseAlterUserStmt 解析Alter User
func (au *AlterUserStmt) ParseAlterUserStmt(node *ast.AlterUserStmt) *AlterUserStmt {
	au.IfExists = node.IfExists

	// CurrentAuth *AuthOption
	ao := new(AuthOption)
	au.CurrentAuth = ao.ParseAuthOption(node.CurrentAuth)

	// Specs []*UserSpec
	for _, aus := range node.Specs {
		us := new(UserSpec)
		usitem := us.ParseUserSpec(aus)

		au.Specs = append(au.Specs, usitem)
	}

	// TLSOptions  []*TLSOption
	for _, ato := range node.TLSOptions {
		to := new(TLSOption)
		toitem := to.ParseTLSOption(ato)

		au.TLSOptions = append(au.TLSOptions, toitem)
	}

	// ResourceOptions []*ResourceOption
	for _, aro := range node.ResourceOptions {
		ro := new(ResourceOption)
		roitem := ro.ParseResourceOption(aro)

		au.ResourceOptions = append(au.ResourceOptions, roitem)
	}

	// PasswordOrLockOptions []*PasswordOrLockOption
	for _, aplo := range node.PasswordOrLockOptions {
		plo := new(PasswordOrLockOption)
		ploitem := plo.ParsePasswordOrLockOption(aplo)

		au.PasswordOrLockOptions = append(au.PasswordOrLockOptions, ploitem)
	}
	return au
}

// AlterTableStmt like *ast.AlterTableStmt
type AlterTableStmt struct {
	Table *TableName
	Specs []*AlterTableSpec
}

// ParseAlterTableStmt 解析ALTER TABLE
func (at *AlterTableStmt) ParseAlterTableStmt(node *ast.AlterTableStmt) *AlterTableStmt {
	// Table *TableName
	t := new(TableName)
	at.Table = t.ParseTableName(node.Table)

	// Specs []*AlterTableSpec
	for _, aats := range node.Specs {
		aat := new(AlterTableSpec)
		atitem := aat.ParseAlterTableSpec(aats)

		at.Specs = append(at.Specs, atitem)
	}

	return at
}

// DropDatabaseStmt like *ast.DropDatabaseStmt
type DropDatabaseStmt struct {
	IfExists bool
	Name     string
}

// ParseDropDatabaseStmt 解析Drop Database
func (dd *DropDatabaseStmt) ParseDropDatabaseStmt(node *ast.DropDatabaseStmt) *DropDatabaseStmt {
	// IfExists bool
	dd.IfExists = node.IfExists

	// Name string
	dd.Name = node.Name

	return dd
}

// DropUserStmt like *ast.DropUserStmt
type DropUserStmt struct {
	IfExists   bool
	IsDropRole bool
	UserList   []*UserIdentity
}

// ParseDropUserStmt 解析Drop User
func (du *DropUserStmt) ParseDropUserStmt(node *ast.DropUserStmt) *DropUserStmt {
	// IfExists   bool
	du.IfExists = node.IfExists

	// IsDropRole bool
	du.IsDropRole = node.IsDropRole

	// UserList   []*auth.UserIdentity
	for _, aui := range node.UserList {
		ui := new(UserIdentity)
		uiitem := ui.ParseUserIdentity(aui)

		du.UserList = append(du.UserList, uiitem)
	}

	return du
}

// DropTableStmt like *ast.DropTableStmt
type DropTableStmt struct {
	IfExists    bool
	Tables      []*TableName
	IsView      bool
	IsTemporary bool // make sense ONLY if/when IsView == false
}

// ParseDropTableStmt 解析Drop Table
func (dt *DropTableStmt) ParseDropTableStmt(node *ast.DropTableStmt) *DropTableStmt {
	// IfExists    bool
	dt.IfExists = node.IfExists

	// Tables      []*TableName
	for _, atn := range node.Tables {
		tn := new(TableName)
		tnitem := tn.ParseTableName(atn)

		dt.Tables = append(dt.Tables, tnitem)
	}

	// IsView      bool
	dt.IsView = node.IsView

	// IsTemporary bool
	dt.IsTemporary = node.IsTemporary

	return dt
}

// DropIndexStmt like *ast.DropIndexStmt
type DropIndexStmt struct {
	IfExists  bool
	IndexName string
	Table     *TableName
	LockAlg   *IndexLockAndAlgorithm
}

// ParseDropIndexStmt 解析Drop Index
func (di *DropIndexStmt) ParseDropIndexStmt(node *ast.DropIndexStmt) *DropIndexStmt {
	// IfExists  bool
	di.IfExists = node.IfExists

	// IndexName string
	di.IndexName = node.IndexName

	// Table     *TableName
	tn := new(TableName)
	di.Table = tn.ParseTableName(node.Table)

	// LockAlg   *IndexLockAndAlgorithm
	ilaa := new(IndexLockAndAlgorithm)
	di.LockAlg = ilaa.ParseIndexLockAndAlgorithm(node.LockAlg)

	return di
}

// DropBindingStmt like *ast.DropBindingStmt
type DropBindingStmt struct {
	GlobalScope bool
	OriginNode  interface{} // StmtNode
	HintedNode  interface{} // StmtNode
}

// ParseDropBindingStmt 解析CREATE BINDING
func (db *DropBindingStmt) ParseDropBindingStmt(node *ast.DropBindingStmt) *DropBindingStmt {
	// GlobalScope bool
	db.GlobalScope = node.GlobalScope

	// OriginNode  StmtNode
	st1 := new(SQLTree)
	st1.ClassifyStmtNode(node.OriginSel)
	db.OriginNode = st1

	// HintedNode  StmtNode
	st2 := new(SQLTree)
	st2.ClassifyStmtNode(node.HintedSel)
	db.HintedNode = st2

	return db
}

// DropSequenceStmt like *ast.DropSequenceStmt
type DropSequenceStmt struct {
	IfExists  bool
	Sequences []*TableName
}

// ParseDropSequenceStmt 解析Drop Sequence
func (ds *DropSequenceStmt) ParseDropSequenceStmt(node *ast.DropSequenceStmt) *DropSequenceStmt {
	// IfExists  bool
	ds.IfExists = node.IfExists

	// Sequences []*TableName
	for _, atn := range node.Sequences {
		tn := new(TableName)
		tnitem := tn.ParseTableName(atn)

		ds.Sequences = append(ds.Sequences, tnitem)
	}

	return ds
}

// DropStatsStmt like *ast.DropStatsStmt
type DropStatsStmt struct {
	Table *TableName
}

// ParseDropStatsStmt 解析Drop Stats
func (ds *DropStatsStmt) ParseDropStatsStmt(node *ast.DropStatsStmt) *DropStatsStmt {
	// Table *TableName
	tn := new(TableName)
	ds.Table = tn.ParseTableName(node.Table)

	return ds
}

// TruncateTableStmt like *ast.TruncateTableStmt
type TruncateTableStmt struct {
	Table *TableName
}

// ParseTruncateTableStmt 解析Truncate Table
func (tt *TruncateTableStmt) ParseTruncateTableStmt(node *ast.TruncateTableStmt) *TruncateTableStmt {
	//Table *TableName
	tn := new(TableName)
	tt.Table = tn.ParseTableName(node.Table)

	return tt
}

// RenameTableStmt like *ast.RenameTableStmt
type RenameTableStmt struct {
	TableToTables []*TableToTable
}

// ParseRenameTableStmt 解析Rename Table
func (rt *RenameTableStmt) ParseRenameTableStmt(node *ast.RenameTableStmt) *RenameTableStmt {
	//TableToTables []*TableToTable
	for _, attt := range node.TableToTables {
		ttt := new(TableToTable)
		tttitem := ttt.ParseTableToTable(attt)

		rt.TableToTables = append(rt.TableToTables, tttitem)
	}

	return rt
}

// GrantStmt like *ast.GrantStmt
type GrantStmt struct {
	Privs      []*PrivElem
	ObjectType string // ObjectTypeType
	Level      *GrantLevel
	Users      []*UserSpec
	TLSOptions []*TLSOption
	WithGrant  bool
}

// ParseGrantStmt 解析Grant
func (g *GrantStmt) ParseGrantStmt(node *ast.GrantStmt) *GrantStmt {
	// Privs      []*PrivElem
	for _, ape := range node.Privs {
		pe := new(PrivElem)
		peitem := pe.ParsePrivElem(ape)

		g.Privs = append(g.Privs, peitem)
	}

	// ObjectType ObjectTypeType
	if val, bv := classifyObjectTypeType(node.ObjectType); bv == true {
		g.ObjectType = val
	}

	// Level      *GrantLevel
	gl := new(GrantLevel)
	g.Level = gl.ParseGrantLevel(node.Level)

	// Users      []*UserSpec
	for _, aus := range node.Users {
		us := new(UserSpec)
		usitem := us.ParseUserSpec(aus)

		g.Users = append(g.Users, usitem)
	}

	// TLSOptions []*TLSOption
	for _, ato := range node.TLSOptions {
		to := new(TLSOption)
		toitem := to.ParseTLSOption(ato)

		g.TLSOptions = append(g.TLSOptions, toitem)
	}

	// WithGrant  bool
	g.WithGrant = node.WithGrant

	return g
}

// RevokeStmt like *ast.RevokeStmt
type RevokeStmt struct {
	Privs      []*PrivElem
	ObjectType string // ObjectTypeType
	Level      *GrantLevel
	Users      []*UserSpec
}

// ParseRevokeStmt 解析Revoke
func (r *RevokeStmt) ParseRevokeStmt(node *ast.RevokeStmt) *RevokeStmt {
	// Privs      []*PrivElem
	for _, ape := range node.Privs {
		pe := new(PrivElem)
		peitem := pe.ParsePrivElem(ape)

		r.Privs = append(r.Privs, peitem)
	}

	// ObjectType ObjectTypeType
	if val, bv := classifyObjectTypeType(node.ObjectType); bv == true {
		r.ObjectType = val
	}

	// Level      *GrantLevel
	gl := new(GrantLevel)
	r.Level = gl.ParseGrantLevel(node.Level)

	// Users      []*UserSpec
	for _, aus := range node.Users {
		us := new(UserSpec)
		usitem := us.ParseUserSpec(aus)

		r.Users = append(r.Users, usitem)
	}

	return r
}

// GrantRoleStmt like *ast.GrantRoleStmt
type GrantRoleStmt struct {
	Roles []*RoleIdentity // []*auth.RoleIdentity
	Users []*UserIdentity // []*auth.UserIdentity
}

// ParseGrantRoleStmt 解析Grant Role
func (gr *GrantRoleStmt) ParseGrantRoleStmt(node *ast.GrantRoleStmt) *GrantRoleStmt {
	// Roles []*auth.RoleIdentity
	for _, ari := range node.Roles {
		ri := new(RoleIdentity)
		riitem := ri.ParseRoleIdentity(ari)

		gr.Roles = append(gr.Roles, riitem)
	}

	// Users []*auth.UserIdentity
	for _, aui := range node.Users {
		ui := new(UserIdentity)
		uiitem := ui.ParseUserIdentity(aui)

		gr.Users = append(gr.Users, uiitem)
	}

	return gr
}

// RevokeRoleStmt like *ast.RevokeRoleStmt
type RevokeRoleStmt struct {
	Roles []*RoleIdentity // []*auth.RoleIdentity
	Users []*UserIdentity // []*auth.UserIdentity
}

// ParseRevokeRoleStmt 解析Revoke Role
func (rr *RevokeRoleStmt) ParseRevokeRoleStmt(node *ast.RevokeRoleStmt) *RevokeRoleStmt {
	// Roles []*auth.RoleIdentity
	for _, ari := range node.Roles {
		ri := new(RoleIdentity)
		riitem := ri.ParseRoleIdentity(ari)

		rr.Roles = append(rr.Roles, riitem)
	}

	// Users []*auth.UserIdentity
	for _, aui := range node.Users {
		ui := new(UserIdentity)
		uiitem := ui.ParseUserIdentity(aui)

		rr.Users = append(rr.Users, uiitem)
	}

	return rr
}

// UseStmt like *ast.UseStmt
type UseStmt struct {
	DBName string
}

// ParseUseStmt 解析Use
func (u *UseStmt) ParseUseStmt(node *ast.UseStmt) *UseStmt {
	// DBName string
	u.DBName = node.DBName

	return u
}

// ShowStmt like *ast.ShowStmt
type ShowStmt struct {
	Type        ShowStmtType // Databases/Tables/Columns/....
	DBName      string
	Table       *TableName        // Used for showing columns.
	Column      *ColumnName       // Used for `desc table column`.
	IndexName   map[string]string // model.CIStr
	Flag        int               // Some flag parsed from sql, such as FULL.
	Full        bool
	User        *UserIdentity   // Used for show grants/create user.
	Roles       []*RoleIdentity // Used for show grants .. using
	IfNotExists bool            // Used for `show create database if not exists`
	Extended    bool            // Used for `show extended columns from ...`
	// GlobalScope is used by `show variables` and `show bindings`
	GlobalScope      bool
	Pattern          *PatternLikeExpr
	Where            interface{} // ExprNode
	ShowProfileTypes []int       // Used for `SHOW PROFILE` syntax
	ShowProfileArgs  *int64      // Used for `SHOW PROFILE` syntax
	ShowProfileLimit *Limit      // Used for `SHOW PROFILE` syntax
}

// ParseShowStmt 解析Show
func (show *ShowStmt) ParseShowStmt(node *ast.ShowStmt) *ShowStmt {
	// Type        ShowStmtType
	show.Type = ShowStmtType(node.Tp)

	// DBName      string
	show.DBName = node.DBName

	// Table       *TableName
	tn := new(TableName)
	show.Table = tn.ParseTableName(node.Table)

	// Column      *ColumnName
	cn := new(ColumnName)
	show.Column = cn.ParseColumnName(node.Column)

	// IndexName   model.CIStr
	show.IndexName = beautifyModelCIStr(&node.IndexName)

	// Flag        int
	show.Flag = node.Flag

	// Full        bool
	show.Full = node.Full

	// User        *UserIdentity
	ui := new(UserIdentity)
	show.User = ui.ParseUserIdentity(node.User)

	// Roles       []*RoleIdentity
	for _, ari := range node.Roles {
		ri := new(RoleIdentity)
		riitem := ri.ParseRoleIdentity(ari)

		show.Roles = append(show.Roles, riitem)
	}

	// IfNotExists bool
	show.IfNotExists = node.IfNotExists

	// Extended    bool
	show.Extended = node.Extended

	// GlobalScope      bool
	show.GlobalScope = node.GlobalScope

	// Pattern          *PatternLikeExpr
	ple := new(PatternLikeExpr)
	show.Pattern = ple.ParsePatternLikeExpr(node.Pattern)

	// Where            ExprNode
	if val, bv := classifyExprNode(node.Where); bv == true {
		show.Where = val
	}

	// ShowProfileTypes []int
	show.ShowProfileTypes = node.ShowProfileTypes

	// ShowProfileArgs  *int64
	show.ShowProfileArgs = node.ShowProfileArgs

	// ShowProfileLimit *Limit
	l := new(Limit)
	show.ShowProfileLimit = l.ParseLimit(node.ShowProfileLimit)

	return show
}

// BeginStmt like *ast.BeginStmt
type BeginStmt struct {
	Mode     string
	ReadOnly bool
	Bound    *TimestampBound
}

// ParseBeginStmt 解析Begin
func (beg *BeginStmt) ParseBeginStmt(node *ast.BeginStmt) *BeginStmt {
	// Mode     string
	beg.Mode = node.Mode

	// ReadOnly bool
	beg.ReadOnly = node.ReadOnly

	// Bound    *TimestampBound
	tsb := new(TimestampBound)
	beg.Bound = tsb.ParseTimestampBound(node.Bound)

	return beg
}

// CommitStmt like *ast.CommitStmt
type CommitStmt struct {
	// CompletionType overwrites system variable `completion_type` within transaction
	CompletionType string // CompletionType
}

// ParseCommitStmt like *ast.ParseCommitStmt
func (cm *CommitStmt) ParseCommitStmt(node *ast.CommitStmt) *CommitStmt {
	// CompletionType CompletionType
	if val, bv := classifyCompletionType(node.CompletionType); bv == true {
		cm.CompletionType = val
	}

	return cm
}

// RollbackStmt like *ast.RollbackStmt
type RollbackStmt struct {
	// CompletionType overwrites system variable `completion_type` within transaction
	CompletionType string // CompletionType
}

// ParseRollbackStmt 解析Rollback
func (rb *RollbackStmt) ParseRollbackStmt(node *ast.RollbackStmt) *RollbackStmt {
	// CompletionType CompletionType
	if val, bv := classifyCompletionType(node.CompletionType); bv == true {
		rb.CompletionType = val
	}

	return rb
}

// PrepareStmt like *ast.PrepareStmt
type PrepareStmt struct {
	Name    string
	SQLText string
	SQLVar  *VariableExpr
}

// ParsePrepareStmt 解析Prepare
func (pp *PrepareStmt) ParsePrepareStmt(node *ast.PrepareStmt) *PrepareStmt {
	// Name    string
	pp.Name = node.Name

	// SQLText string
	pp.SQLText = node.SQLText

	// SQLVar  *VariableExpr
	ve := new(VariableExpr)
	pp.SQLVar = ve.ParseVariableExpr(node.SQLVar)

	return pp
}

// ExecuteStmt like *ast.ExecuteStmt
type ExecuteStmt struct {
	Name       string
	UsingVars  []interface{} // []ExprNode
	BinaryArgs interface{}
	ExecID     uint32
	IdxInMulti int
}

// ParseExecuteStmt 解析Execute
func (exe *ExecuteStmt) ParseExecuteStmt(node *ast.ExecuteStmt) *ExecuteStmt {
	// Name       string
	exe.Name = node.Name

	// UsingVars  []ExprNode
	for _, auv := range node.UsingVars {
		if val, bv := classifyExprNode(auv); bv == true {
			uvitem := val

			exe.UsingVars = append(exe.UsingVars, uvitem)
		}
	}

	// BinaryArgs interface{}
	exe.BinaryArgs = node.BinaryArgs

	// ExecID     uint32
	exe.ExecID = node.ExecID

	// IdxInMulti int
	exe.IdxInMulti = node.IdxInMulti

	return exe
}

// ExplainStmt like *ast.ExplainStmt
type ExplainStmt struct {
	Stmt    interface{} // StmtNode
	Format  string
	Analyze bool
}

// ParseExplainStmt 解析Explain
func (ep *ExplainStmt) ParseExplainStmt(node *ast.ExplainStmt) *ExplainStmt {
	// Stmt    StmtNode
	st := new(SQLTree)
	st.ClassifyStmtNode(node.Stmt)
	ep.Stmt = st

	// Format  string
	ep.Format = node.Format

	// Analyze bool
	ep.Analyze = node.Analyze

	return ep
}

// ExplainForStmt like *ast.ExplainForStmt
type ExplainForStmt struct {
	Format       string
	ConnectionID string // uint64
}

// ParseExplainForStmt 解析Explain For
func (epf *ExplainForStmt) ParseExplainForStmt(node *ast.ExplainForStmt) *ExplainForStmt {
	// Format       string
	epf.Format = node.Format

	// ConnectionID uint64
	epf.ConnectionID = strconv.FormatUint(node.ConnectionID, 10)

	return epf
}

// BinlogStmt like *ast.BinlogStmt
type BinlogStmt struct {
	StrValue string
}

// FlushStmt like *ast.FlushStmt
type FlushStmt struct {
	Type            string // FlushStmtType // Privileges/Tables/...
	NoWriteToBinLog bool
	LogType         string       // LogType
	Tables          []*TableName // For FlushTableStmt, if Tables is empty, it means flush all tables.
	ReadLock        bool
	Plugins         []string
}

// ParseFlushStmt 解析Flush
func (flu *FlushStmt) ParseFlushStmt(node *ast.FlushStmt) *FlushStmt {
	// Type            FlushStmtType
	if val, bv := classifyFlushStmtType(node.Tp); bv == true {
		flu.Type = val
	}

	// NoWriteToBinLog bool
	flu.NoWriteToBinLog = node.NoWriteToBinLog

	// LogType         LogType
	if val, bv := classifyLogType(node.LogType); bv == true {
		flu.LogType = val
	}

	// Tables          []*TableName
	for _, atn := range node.Tables {
		tn := new(TableName)
		tnitem := tn.ParseTableName(atn)

		flu.Tables = append(flu.Tables, tnitem)
	}

	// ReadLock        bool
	flu.ReadLock = node.ReadLock

	// Plugins         []string
	flu.Plugins = node.Plugins

	return flu
}

// ParseBinlogStmt 解析内部Binlog
func (bin *BinlogStmt) ParseBinlogStmt(node *ast.BinlogStmt) *BinlogStmt {
	// Str string
	bin.StrValue = node.Str

	return bin
}

// ChangeStmt like *ast.ChangeStmt
type ChangeStmt struct {
	NodeType string
	State    string
	NodeID   string
}

// ParseChangeStmt 解析Change
func (chg *ChangeStmt) ParseChangeStmt(node *ast.ChangeStmt) *ChangeStmt {
	// NodeType string
	chg.NodeType = node.NodeType

	// State    string
	chg.State = node.State

	// NodeID   string
	chg.NodeID = node.NodeID

	return chg
}

// CleanupTableLockStmt like *ast.CleanupTableLockStmt
type CleanupTableLockStmt struct {
	Tables []*TableName
}

// ParseCleanupTableLockStmt 解析Cleanup Table Lock
func (ctl *CleanupTableLockStmt) ParseCleanupTableLockStmt(node *ast.CleanupTableLockStmt) *CleanupTableLockStmt {
	// Tables []*TableName
	for _, atn := range node.Tables {
		tn := new(TableName)
		tnitem := tn.ParseTableName(atn)

		ctl.Tables = append(ctl.Tables, tnitem)
	}

	return ctl
}

// DeallocateStmt like *ast.DeallocateStmt
type DeallocateStmt struct {
	Name string
}

// ParseDeallocateStmt 解析Deallocate
func (dlc *DeallocateStmt) ParseDeallocateStmt(node *ast.DeallocateStmt) *DeallocateStmt {
	// Name string
	dlc.Name = node.Name

	return dlc
}

// AdminStmt like *ast.AdminStmt
type AdminStmt struct {
	Type         AdminStmtType
	Index        string
	Tables       []*TableName
	JobIDs       []int64
	JobNumber    int64
	HandleRanges []HandleRange
	ShowSlow     *ShowSlow
	Plugins      []string
	Where        interface{} // ExprNode
}

// ParseAdminStmt 解析Admin
func (adm *AdminStmt) ParseAdminStmt(node *ast.AdminStmt) *AdminStmt {
	// Type         AdminStmtType
	adm.Type = AdminStmtType(node.Tp)

	// Index        string
	adm.Index = node.Index

	// Tables       []*TableName
	for _, atn := range node.Tables {
		tn := new(TableName)
		tnitem := tn.ParseTableName(atn)

		adm.Tables = append(adm.Tables, tnitem)
	}

	// JobIDs       []int64
	adm.JobIDs = node.JobIDs

	// JobNumber    int64
	adm.JobNumber = node.JobNumber

	// HandleRanges []HandleRange
	for _, ahr := range node.HandleRanges {
		hr := new(HandleRange)
		hritem := hr.ParseHandleRange(&ahr)

		adm.HandleRanges = append(adm.HandleRanges, *hritem)
	}

	// ShowSlow     *ShowSlow
	ss := new(ShowSlow)
	adm.ShowSlow = ss.ParseShowSlow(node.ShowSlow)

	// Plugins      []string
	adm.Plugins = node.Plugins

	// Where        ExprNode
	if val, bv := classifyExprNode(node.Where); bv == true {
		adm.Where = val
	}

	return adm
}

// DoStmt like *ast.DoStmt
type DoStmt struct {
	Exprs []interface{} // []ExprNode
}

// ParseDoStmt 解析Do
func (do *DoStmt) ParseDoStmt(node *ast.DoStmt) *DoStmt {
	// Exprs []ExprNode
	for _, aen := range node.Exprs {
		if val, bv := classifyExprNode(aen); bv == true {
			enitem := val

			do.Exprs = append(do.Exprs, enitem)
		}
	}

	return do
}

// IndexAdviseStmt like *ast.IndexAdviseStmt
type IndexAdviseStmt struct {
	IsLocal     bool
	Path        string
	MaxMinutes  string // uint64
	MaxIndexNum *MaxIndexNumClause
	LinesInfo   *LinesClause
}

// ParseIndexAdviseStmt 解析IndexAdvise
func (ia *IndexAdviseStmt) ParseIndexAdviseStmt(node *ast.IndexAdviseStmt) *IndexAdviseStmt {
	// IsLocal     bool
	ia.IsLocal = node.IsLocal

	// Path        string
	ia.Path = node.Path

	// MaxMinutes  uint64
	ia.MaxMinutes = strconv.FormatUint(node.MaxMinutes, 10)

	// MaxIndexNum *MaxIndexNumClause
	minc := new(MaxIndexNumClause)
	ia.MaxIndexNum = minc.ParseMaxIndexNumClause(node.MaxIndexNum)

	// LinesInfo   *LinesClause
	lc := new(LinesClause)
	ia.LinesInfo = lc.ParseLinesClause(node.LinesInfo)

	return ia
}

// KillStmt like *ast.KillStmt
type KillStmt struct {
	// Query indicates whether terminate a single query on this connection or the whole connection.
	// If Query is true, terminates the statement the connection is currently executing, but leaves the connection itself intact.
	// If Query is false, terminates the connection associated with the given ConnectionID, after terminating any statement the connection is executing.
	Query        bool
	ConnectionID string // uint64
	// TiDBExtension is used to indicate whether the user knows he is sending kill statement to the right tidb-server.
	// When the SQL grammar is "KILL TIDB [CONNECTION | QUERY] connectionID", TiDBExtension will be set.
	// It's a special grammar extension in TiDB. This extension exists because, when the connection is:
	// client -> LVS proxy -> TiDB, and type Ctrl+C in client, the following action will be executed:
	// new a connection; kill xxx;
	// kill command may send to the wrong TiDB, because the exists of LVS proxy, and kill the wrong session.
	// So, "KILL TIDB" grammar is introduced, and it REQUIRES DIRECT client -> TiDB TOPOLOGY.
	// TODO: The standard KILL grammar will be supported once we have global connectionID.
	TiDBExtension bool
}

// ParseKillStmt 解析Kill
func (kl *KillStmt) ParseKillStmt(node *ast.KillStmt) *KillStmt {
	// Query        bool
	kl.Query = node.Query

	// ConnectionID uint64
	kl.ConnectionID = strconv.FormatUint(node.ConnectionID, 10)

	// TiDBExtension bool
	kl.TiDBExtension = node.TiDBExtension

	return kl
}

// LoadDataStmt like *ast.LoadDataStmt
type LoadDataStmt struct {
	IsLocal            bool
	Path               string
	OnDuplicate        string // OnDuplicateKeyHandlingType
	Table              *TableName
	Columns            []*ColumnName
	FieldsInfo         *FieldsClause
	LinesInfo          *LinesClause
	IgnoreLines        string // uint64
	ColumnAssignments  []*Assignment
	ColumnsAndUserVars []*ColumnNameOrUserVar
}

// ParseLoadDataStmt 解析Load Data
func (ld *LoadDataStmt) ParseLoadDataStmt(node *ast.LoadDataStmt) *LoadDataStmt {
	// IsLocal            bool
	ld.IsLocal = node.IsLocal

	// Path               string
	ld.Path = node.Path

	// OnDuplicate        OnDuplicateKeyHandlingType
	if val, bv := classifyOnDuplicateKeyHandlingType(node.OnDuplicate); bv == true {
		ld.OnDuplicate = val
	}

	// Table              *TableName
	tn := new(TableName)
	ld.Table = tn.ParseTableName(node.Table)

	// Columns            []*ColumnName
	// FieldsInfo         *FieldsClause
	fc := new(FieldsClause)
	ld.FieldsInfo = fc.ParseFieldsClause(node.FieldsInfo)

	// LinesInfo          *LinesClause
	lc := new(LinesClause)
	ld.LinesInfo = lc.ParseLinesClause(node.LinesInfo)

	// IgnoreLines        uint64
	ld.IgnoreLines = strconv.FormatUint(node.IgnoreLines, 10)

	// ColumnAssignments  []*Assignment
	// ColumnsAndUserVars []*ColumnNameOrUserVar

	return ld
}

// LoadStatsStmt like *ast.LoadStatsStmt
type LoadStatsStmt struct {
	Path string
}

// ParseLoadStatsStmt 解析Load Stats
func (ls *LoadStatsStmt) ParseLoadStatsStmt(node *ast.LoadStatsStmt) *LoadStatsStmt {
	//  Path string
	ls.Path = node.Path

	return ls
}

// LockTablesStmt like *ast.LockTablesStmt
type LockTablesStmt struct {
	TableLocks []*TableLock
}

// ParseLockTablesStmt 解析Lock Tables
func (lt *LockTablesStmt) ParseLockTablesStmt(node *ast.LockTablesStmt) *LockTablesStmt {
	// TableLocks []*TableLock
	for _, atl := range node.TableLocks {
		tl := new(TableLock)
		tlitem := tl.ParseTableLock(&atl)

		lt.TableLocks = append(lt.TableLocks, tlitem)
	}

	return lt
}

// RecoverTableStmt like *ast.RecoverTableStmt
type RecoverTableStmt struct {
	JobID  int64
	Table  *TableName
	JobNum int64
}

// ParseRecoverTableStmt 解析Recover Table
func (rct *RecoverTableStmt) ParseRecoverTableStmt(node *ast.RecoverTableStmt) *RecoverTableStmt {
	// JobID  int64
	rct.JobID = node.JobID

	// Table  *TableName
	tn := new(TableName)
	rct.Table = tn.ParseTableName(node.Table)

	// JobNum int64
	rct.JobNum = node.JobNum

	return rct
}

// RepairTableStmt like *ast.RepairTableStmt
type RepairTableStmt struct {
	Table      *TableName
	CreateStmt *CreateTableStmt
}

// ParseRepairTableStmt 解析Repair Table
func (rpt *RepairTableStmt) ParseRepairTableStmt(node *ast.RepairTableStmt) *RepairTableStmt {
	// Table      *TableName
	tn := new(TableName)
	rpt.Table = tn.ParseTableName(node.Table)

	// CreateStmt *CreateTableStmt
	ct := new(CreateTableStmt)
	rpt.CreateStmt = ct.ParseCreateTableStmt(node.CreateStmt)

	return rpt
}

// SetStmt like *ast.SetStmt
type SetStmt struct {
	// Variables is the list of variable assignment.
	Variables []*VariableAssignment
}

// ParseSetStmt 解析Set
func (set *SetStmt) ParseSetStmt(node *ast.SetStmt) *SetStmt {
	// Variables []*VariableAssignment
	for _, ava := range node.Variables {
		va := new(VariableAssignment)
		vaitem := va.ParserVariableAssignment(ava)

		set.Variables = append(set.Variables, vaitem)
	}

	return set
}

// SetConfigStmt like *ast.SetConfigStmt
type SetConfigStmt struct {
	Type     string      // TiDB, TiKV, PD
	Instance string      // '127.0.0.1:3306'
	Name     string      // the variable name
	Value    interface{} // ExprNode
}

// ParseSetConfigStmt 解析Set Config
func (scf *SetConfigStmt) ParseSetConfigStmt(node *ast.SetConfigStmt) *SetConfigStmt {
	// Type     string
	scf.Type = node.Type

	// Instance string
	scf.Instance = node.Instance

	// Name     string
	scf.Name = node.Name

	// Value    ExprNode
	if val, bv := classifyExprNode(node.Value); bv == true {
		scf.Value = val
	}

	return scf
}

// SetDefaultRoleStmt like *ast.SetDefaultRoleStmt
type SetDefaultRoleStmt struct {
	SetRoleOpt string // SetRoleStmtType
	RoleList   []*RoleIdentity
	UserList   []*UserIdentity
}

// ParseSetDefaultRoleStmt 解析Set Default Role
func (sdr *SetDefaultRoleStmt) ParseSetDefaultRoleStmt(node *ast.SetDefaultRoleStmt) *SetDefaultRoleStmt {
	// SetRoleOpt SetRoleStmtType
	if val, bv := classifySetRoleStmtType(node.SetRoleOpt); bv == true {
		sdr.SetRoleOpt = val
	}

	// RoleList   []*auth.RoleIdentity
	for _, ari := range node.RoleList {
		ri := new(RoleIdentity)
		riitem := ri.ParseRoleIdentity(ari)

		sdr.RoleList = append(sdr.RoleList, riitem)
	}

	// UserList   []*auth.UserIdentity
	for _, aui := range node.UserList {
		ui := new(UserIdentity)
		uiitem := ui.ParseUserIdentity(aui)

		sdr.UserList = append(sdr.UserList, uiitem)
	}

	return sdr
}

// SetPwdStmt like *ast.SetPwdStmt
type SetPwdStmt struct {
	User     *UserIdentity
	Password string
}

// ParseSetPwdStmt 解析Set Password
func (spw *SetPwdStmt) ParseSetPwdStmt(node *ast.SetPwdStmt) *SetPwdStmt {
	// User     *auth.UserIdentity
	ui := new(UserIdentity)
	spw.User = ui.ParseUserIdentity(node.User)

	// Password string
	spw.Password = node.Password

	return spw
}

// SetRoleStmt like *ast.SetRoleStmt
type SetRoleStmt struct {
	SetRoleOpt string // SetRoleStmtType
	RoleList   []*RoleIdentity
}

// ParseSetRoleStmt 解析Set Role
func (sro *SetRoleStmt) ParseSetRoleStmt(node *ast.SetRoleStmt) *SetRoleStmt {
	// SetRoleOpt SetRoleStmtType
	if val, bv := classifySetRoleStmtType(node.SetRoleOpt); bv == true {
		sro.SetRoleOpt = val
	}
	// RoleList   []*RoleIdentity
	for _, ari := range node.RoleList {
		ri := new(RoleIdentity)
		riitem := ri.ParseRoleIdentity(ari)

		sro.RoleList = append(sro.RoleList, riitem)
	}

	return sro
}

// SplitRegionStmt like *ast.SplitRegionStmt
type SplitRegionStmt struct {
	Table          *TableName
	IndexName      map[string]string   // model.CIStr
	PartitionNames []map[string]string // []model.CIStr
	SplitSyntaxOpt *SplitSyntaxOption
	SplitOpt       *SplitOption
}

// ParseSplitRegionStmt 解析Split Region
func (sre *SplitRegionStmt) ParseSplitRegionStmt(node *ast.SplitRegionStmt) *SplitRegionStmt {
	// Table          *TableName
	tn := new(TableName)
	sre.Table = tn.ParseTableName(node.Table)

	// IndexName      model.CIStr
	sre.IndexName = beautifyModelCIStr(&node.IndexName)

	// PartitionNames []model.CIStr
	for _, apn := range node.PartitionNames {
		pnitem := beautifyModelCIStr(&apn)

		sre.PartitionNames = append(sre.PartitionNames, pnitem)
	}

	// SplitSyntaxOpt *SplitSyntaxOption
	sso := new(SplitSyntaxOption)
	sre.SplitSyntaxOpt = sso.ParseSplitSyntaxOption(node.SplitSyntaxOpt)

	// SplitOpt       *SplitOption
	so := new(SplitOption)
	sre.SplitOpt = so.ParseSplitOption(node.SplitOpt)

	return sre
}

// TraceStmt like *ast.TraceStmt
type TraceStmt struct {
	Stmt   interface{} // StmtNode
	Format string
}

// ParseTraceStmt 解析Trace
func (trc *TraceStmt) ParseTraceStmt(node *ast.TraceStmt) *TraceStmt {
	// Stmt   StmtNode
	st := new(SQLTree)
	st.ClassifyStmtNode(node.Stmt)
	trc.Stmt = st

	// Format string
	trc.Format = node.Format

	return trc
}

// AlterTableSpec like *ast.AlterTableSpec
type AlterTableSpec struct {
	// only supported by MariaDB 10.0.2+ (DROP COLUMN, CHANGE COLUMN, MODIFY COLUMN, DROP INDEX, DROP FOREIGN KEY, DROP PARTITION). see https://mariadb.com/kb/en/library/alter-table/
	IfExists bool
	// only supported by MariaDB 10.0.2+ (ADD COLUMN, ADD PARTITION). see https://mariadb.com/kb/en/library/alter-table/
	IfNotExists     bool
	NoWriteToBinlog bool
	OnAllPartitions bool
	Type            string // Tp  AlterTableType
	Name            string
	IndexName       map[string]string // model.CIStr
	Constraint      *Constraint
	Options         []*TableOption
	OrderByList     []*AlterOrderItem
	NewTable        *TableName
	NewColumns      []*ColumnDef
	NewConstraints  []*Constraint
	OldColumnName   *ColumnName
	NewColumnName   *ColumnName
	Position        *ColumnPosition
	LockType        string // LockType
	Algorithm       string // AlgorithmType
	Comment         string
	FromKey         map[string]string // model.CIStr
	ToKey           map[string]string // model.CIStr
	Partition       *PartitionOptions
	PartitionNames  []map[string]string // []model.CIStr
	PartDefinitions []*PartitionDefinition
	WithValidation  bool
	Num             string // uint64
	Visibility      string // Visibility IndexVisibility
	TiFlashReplica  *TiFlashReplicaSpec
	//PlacementSpecs  []*ast.PlacementSpec
	Writeable bool
}

// ParseAlterTableSpec 解析Alter Table的内容项
func (ats *AlterTableSpec) ParseAlterTableSpec(aats *ast.AlterTableSpec) *AlterTableSpec {
	if aats == nil {
		return nil
	}

	// IfExists bool
	ats.IfExists = aats.IfExists

	// IfNotExists     bool
	ats.IfNotExists = aats.IfNotExists

	// NoWriteToBinlog bool
	ats.NoWriteToBinlog = aats.NoWriteToBinlog

	// OnAllPartitions bool
	ats.OnAllPartitions = aats.OnAllPartitions

	// Tp AlterTableType
	if val, bv := classifyAlterTableType(aats); bv == true {
		ats.Type = val
	}

	// Name string
	ats.Name = aats.Name

	// IndexName model.CIStr
	//ats.IndexName = beautifyModelCIStr(&aats.IndexName)

	// Constraint *Constraint
	cst := new(Constraint)
	ats.Constraint = cst.ParseConstraint(aats.Constraint)

	// Options []*TableOption
	for _, ato := range aats.Options {
		to := new(TableOption)
		toitem := to.ParseTableOption(ato)

		ats.Options = append(ats.Options, toitem)
	}

	// OrderByList     []*AlterOrderItem
	for _, aaoi := range aats.OrderByList {
		aoi := new(AlterOrderItem)
		aoiitem := aoi.ParseAlterOrderItem(aaoi)

		ats.OrderByList = append(ats.OrderByList, aoiitem)
	}

	// NewTable *TableName
	tn := new(TableName)
	ats.NewTable = tn.ParseTableName(aats.NewTable)

	// NewColumns []*ColumnDef
	for _, acd := range aats.NewColumns {
		cd := new(ColumnDef)
		cditem := cd.ParseColumnDef(acd)

		ats.NewColumns = append(ats.NewColumns, cditem)
	}

	// NewConstraints  []*Constraint
	for _, acst := range aats.NewConstraints {
		cst := new(Constraint)
		cstitem := cst.ParseConstraint(acst)

		ats.NewConstraints = append(ats.NewConstraints, cstitem)
	}

	// OldColumnName   *ColumnName
	ocn := new(ColumnName)
	ats.OldColumnName = ocn.ParseColumnName(aats.OldColumnName)

	// NewColumnName   *ColumnName
	ncn := new(ColumnName)
	ats.NewColumnName = ncn.ParseColumnName(aats.NewColumnName)

	// Position        *ColumnPosition
	cp := new(ColumnPosition)
	ats.Position = cp.ParseColumnPosition(aats.Position)

	// LockType LockType
	if val, bv := classifyLockType(aats.LockType); bv == true {
		ats.LockType = val
	}

	// Algorithm AlgorithmType
	if val, bv := classifyAlgorithmType(aats.Algorithm); bv == true {
		ats.Algorithm = val
	}

	// Comment string
	ats.Comment = aats.Comment

	// FromKey model.CIStr
	ats.FromKey = beautifyModelCIStr(&aats.FromKey)

	// ToKey model.CIStr
	ats.ToKey = beautifyModelCIStr(&aats.ToKey)

	// Partition *PartitionOptions
	paro := new(PartitionOptions)
	ats.Partition = paro.ParsePartitionOptions(aats.Partition)

	// PartitionNames  []model.CIStr
	for _, aparn := range aats.PartitionNames {
		parn := beautifyModelCIStr(&aparn)

		ats.PartitionNames = append(ats.PartitionNames, parn)
	}

	// PartDefinitions []*PartitionDefinition
	for _, apard := range aats.PartDefinitions {
		pard := new(PartitionDefinition)
		parditem := pard.ParsePartitionDefinition(apard)

		ats.PartDefinitions = append(ats.PartDefinitions, parditem)
	}

	// WithValidation  bool
	ats.WithValidation = aats.WithValidation

	// Num uint64
	ats.Num = strconv.FormatUint(aats.Num, 10)

	// Visibility      IndexVisibility
	if val, bv := classifyIndexVisibility(aats.Visibility); bv == true {
		ats.Visibility = val
	}

	// TiFlashReplica  *TiFlashReplicaSpec
	if aats.TiFlashReplica != nil {
		tfr := new(TiFlashReplicaSpec)
		ats.TiFlashReplica = tfr.ParseTiFlashReplicaSpec(aats.TiFlashReplica)
	}

	// PlacementSpecs  []*PlacementSpec

	// Writeable bool
	//ats.Writeable = aats.Writeable

	return ats
}

// WindowSpec like *ast.WindowSpec
type WindowSpec struct {
	Name map[string]string // model.CIStr
	// Ref is the reference window of this specification. For example, in `w2 as (w1 order by a)`,
	// the definition of `w2` references `w1`.
	Ref         map[string]string // model.CIStr
	PartitionBy *PartitionByClause
	OrderBy     *OrderByClause
	Frame       *FrameClause
	// OnlyAlias will set to true of the first following case.
	// To make compatible with MySQL, we need to distinguish `select func over w` from `select func over (w)`.
	OnlyAlias bool
}

// ParseWindowSpec 解析WindowSpec
func (ws *WindowSpec) ParseWindowSpec(aws *ast.WindowSpec) *WindowSpec {
	if aws == nil {
		return nil
	}

	// Name        model.CIStr
	ws.Name = beautifyModelCIStr(&aws.Name)

	// Ref         model.CIStr
	ws.Ref = beautifyModelCIStr(&aws.Ref)

	// PartitionBy *PartitionByClause
	parc := new(PartitionByClause)
	ws.PartitionBy = parc.ParsePartitionByClause(aws.PartitionBy)

	// OrderBy     *OrderByClause
	oc := new(OrderByClause)
	ws.OrderBy = oc.ParseOrderByClause(aws.OrderBy)

	// Frame       *FrameClause
	fc := new(FrameClause)
	ws.Frame = fc.ParseFrameClause(aws.Frame)

	// OnlyAlias   bool
	ws.OnlyAlias = aws.OnlyAlias

	return ws
}

// UserSpec like *ast.UserSpec
type UserSpec struct {
	User    *UserIdentity
	AuthOpt *AuthOption
	IsRole  bool
}

// ParseUserSpec 解析UserSpec
func (us *UserSpec) ParseUserSpec(aus *ast.UserSpec) *UserSpec {
	if aus == nil {
		return nil
	}

	// User *UserIdentity
	ui := new(UserIdentity)
	us.User = ui.ParseUserIdentity(aus.User)

	// AuthOpt *AuthOption
	ao := new(AuthOption)
	us.AuthOpt = ao.ParseAuthOption(aus.AuthOpt)

	// IsRole  bool
	us.IsRole = aus.IsRole

	return us
}

// TiFlashReplicaSpec like *ast.TiFlashReplicaSpec
type TiFlashReplicaSpec struct {
	Count  string // uint64
	Labels []string
}

// ParseTiFlashReplicaSpec 解析ParseTiFlashReplicaSpec
func (tfr *TiFlashReplicaSpec) ParseTiFlashReplicaSpec(atfr *ast.TiFlashReplicaSpec) *TiFlashReplicaSpec {
	if atfr == nil {
		return nil
	}

	// Count uint64
	tfr.Count = strconv.FormatUint(atfr.Count, 10)

	// Labels []string
	tfr.Labels = atfr.Labels

	return tfr
}

// Join like *ast.Join
type Join struct {
	// Left table can be TableSource or JoinNode.
	Left interface{} // ResultSetNode
	// Right table can be TableSource or JoinNode or nil.
	Right interface{} // ResultSetNode
	// Tp represents join type.
	Type string // JoinType
	// On represents join on condition.
	On *OnCondition
	// Using represents join using clause.
	Using []*ColumnName
	// NaturalJoin represents join is natural join.
	NaturalJoin bool
	// StraightJoin represents a straight join.
	StraightJoin bool
	//ExplicitParens bool
}

// ParseJoin 解析Join
func (j *Join) ParseJoin(aj *ast.Join) *Join {
	if aj == nil {
		return nil
	}

	// Left ResultSetNode
	if val, bv := classifyResultSetNode(aj.Left); bv == true {
		j.Left = val
	}

	// Right ResultSetNode
	if val, bv := classifyResultSetNode(aj.Right); bv == true {
		j.Right = val
	}

	// Tp  JoinType
	if val, bv := classifyJoinType(aj); bv == true {
		j.Type = val
	}

	// On  *OnCondition
	oc := new(OnCondition)
	j.On = oc.ParseOnCondition(aj.On)

	// Using []*ColumnName
	for _, acn := range aj.Using {
		cn := new(ColumnName)

		cnitem := cn.ParseColumnName(acn)
		j.Using = append(j.Using, cnitem)
	}

	// NaturalJoin bool
	j.NaturalJoin = aj.NaturalJoin

	// StraightJoin bool
	j.StraightJoin = aj.StraightJoin

	//ExplicitParens bool

	return j
}

// OnCondition like *ast.OnCondition
type OnCondition struct {
	Expr interface{} // ExprNode
}

// ParseOnCondition 解析OnCondition
func (oc *OnCondition) ParseOnCondition(aoc *ast.OnCondition) *OnCondition {
	if aoc == nil {
		return nil
	}

	// Expr ExprNode
	if val, bv := classifyExprNode(aoc.Expr); bv == true {
		oc.Expr = val
	}

	return oc
}

// Limit like *ast.Limit
type Limit struct {
	Count  interface{} // ExprNode
	Offset interface{} // ExprNode
}

// ParseLimit 解析limit的count和offset
func (l *Limit) ParseLimit(al *ast.Limit) *Limit {
	if al == nil {
		return nil
	}

	// Count  ExprNode
	if val, bv := classifyExprNode(al.Count); bv == true {
		l.Count = val
	}

	// Offset ExprNode
	if val, bv := classifyExprNode(al.Offset); bv == true {
		l.Offset = val
	}

	return l
}

// AlterOrderItem like *ast.AlterOrderItem
type AlterOrderItem struct {
	Column *ColumnName
	Desc   bool
}

// ParseAlterOrderItem 解析alter语句中的orderitem
func (aoi *AlterOrderItem) ParseAlterOrderItem(aaoi *ast.AlterOrderItem) *AlterOrderItem {
	if aaoi == nil {
		return nil
	}

	// Column *ColumnName
	cn := new(ColumnName)
	aoi.Column = cn.ParseColumnName(aaoi.Column)

	// Desc   bool
	aoi.Desc = aaoi.Desc

	return aoi
}

// ByItem like *ast.ByItem
type ByItem struct {
	Expr      interface{} // ExprNode
	Desc      bool
	NullOrder bool
}

// ParseByItem 解析ByItem
func (bi *ByItem) ParseByItem(abi *ast.ByItem) *ByItem {
	if abi == nil {
		return nil
	}

	// Expr ExprNode
	if val, bv := classifyExprNode(abi.Expr); bv == true {
		bi.Expr = val
	}

	// Desc bool
	bi.Desc = abi.Desc

	// NullOrder bool
	// bi.NullOrder = abi.NullOrder

	return bi
}

// Assignment like a = 1.
type Assignment struct {
	// Column is the column name to be assigned.
	Column *ColumnName
	// Expr is the expression assigning to ColName.
	Expr interface{} // ExprNode
}

// ParserAssignment 解析Assignment
func (a *Assignment) ParserAssignment(aa *ast.Assignment) *Assignment {
	if aa == nil {
		return nil
	}
	// Column *ColumnName
	cn := new(ColumnName)
	a.Column = cn.ParseColumnName(aa.Column)

	// Expr ExprNode
	if val, bv := classifyExprNode(aa.Expr); bv == true {
		a.Expr = val
	}

	return a
}

// VariableAssignment like *ast.VariableAssignment
type VariableAssignment struct {
	Name     string
	Value    interface{} // ExprNode
	IsGlobal bool
	IsSystem bool
	// ExtendValue is a way to store extended info.
	// VariableAssignment should be able to store information for SetCharset/SetPWD Stmt.
	// For SetCharsetStmt, Value is charset, ExtendValue is collation.
	// TODO: Use SetStmt to implement set password statement.
	ExtendValue interface{} // ValueExpr
}

// ParserVariableAssignment 解析Variable Assignment
func (va *VariableAssignment) ParserVariableAssignment(ava *ast.VariableAssignment) *VariableAssignment {
	if ava == nil {
		return nil
	}

	// Name     string
	va.Name = ava.Name

	// Value    ExprNode
	if val, bv := classifyExprNode(ava.Value); bv == true {
		va.Value = val
	}

	// IsGlobal bool
	va.IsGlobal = ava.IsGlobal

	// IsSystem bool
	va.IsSystem = ava.IsSystem

	// ExtendValue ValueExpr
	va.ExtendValue = ava.ExtendValue.GetValue()

	return va
}

// DeleteTableList like *ast.DeleteTableList
type DeleteTableList struct {
	Tables []*TableName
}

// ParserDeleteTableList 解析DeleteTable List
func (dtl *DeleteTableList) ParserDeleteTableList(adtl *ast.DeleteTableList) *DeleteTableList {
	if adtl == nil {
		return nil
	}

	// Tables []*TableName
	for _, atn := range adtl.Tables {
		tn := new(TableName)
		tnitem := tn.ParseTableName(atn)

		dtl.Tables = append(dtl.Tables, tnitem)
	}

	return dtl
}

// ReferenceDef like *ast.ReferenceDef
type ReferenceDef struct {
	Table                   *TableName
	IndexPartSpecifications []*IndexPartSpecification
	OnDelete                *OnDeleteOpt
	OnUpdate                *OnUpdateOpt
	Match                   string // MatchType
}

// ParseReferenceDef 解析Reference Definition
func (rd *ReferenceDef) ParseReferenceDef(ard *ast.ReferenceDef) *ReferenceDef {
	if ard == nil {
		return nil
	}

	// Table                   *TableName
	tn := new(TableName)
	rd.Table = tn.ParseTableName(ard.Table)

	// IndexPartSpecifications []*IndexPartSpecification

	// OnDelete                *OnDeleteOpt
	odo := new(OnDeleteOpt)
	rd.OnDelete = odo.ParseOnDeleteOpt(ard.OnDelete)

	// OnUpdate                *OnUpdateOpt
	ouo := new(OnUpdateOpt)
	rd.OnUpdate = ouo.ParseOnUpdateOpt(ard.OnUpdate)

	// Match                   MatchType
	if val, bv := classifyMatchType(ard.Match); bv == true {
		rd.Match = val
	}

	return rd
}

// TimestampBound like *ast.TimestampBound
type TimestampBound struct {
	Mode      string      // TimestampBoundMode
	Timestamp interface{} // ExprNode
}

// ParseTimestampBound 解析TimestampBound
func (tsb *TimestampBound) ParseTimestampBound(atsb *ast.TimestampBound) *TimestampBound {
	if atsb == nil {
		return nil
	}

	// Mode      TimestampBoundMode
	if val, bv := classifyTimestampBoundMode(atsb.Mode); bv == true {
		tsb.Mode = val
	}

	// Timestamp ExprNode
	if val, bv := classifyExprNode(atsb.Timestamp); bv == true {
		tsb.Timestamp = val
	}

	return tsb
}

// ColumnNameOrUserVar like *ast.ColumnNameOrUserVar
type ColumnNameOrUserVar struct {
	ColumnName *ColumnName
	UserVar    *VariableExpr
}

// ParseColumnNameOrUserVar 解析Column Name Or User Variable
func (cou *ColumnNameOrUserVar) ParseColumnNameOrUserVar(acou *ast.ColumnNameOrUserVar) *ColumnNameOrUserVar {
	if acou == nil {
		return nil
	}

	// ColumnName *ColumnName
	cn := new(ColumnName)
	cou.ColumnName = cn.ParseColumnName(acou.ColumnName)

	// UserVar    *VariableExpr
	ve := new(VariableExpr)
	cou.UserVar = ve.ParseVariableExpr(acou.UserVar)

	return cou
}

// AdminStmtType like *ast.AdminStmtType
type AdminStmtType int

// HandleRange like *ast.HandleRange
type HandleRange struct {
	Begin int64
	End   int64
}

// ParseHandleRange 解析Handle Range
func (hr *HandleRange) ParseHandleRange(ahr *ast.HandleRange) *HandleRange {
	if ahr == nil {
		return nil
	}

	// Begin int64
	hr.Begin = ahr.Begin

	// End   int64
	hr.End = ahr.End

	return hr
}

// ShowSlow like *ast.ShowSlow
type ShowSlow struct {
	Type  string // ShowSlowType
	Count string // uint64
	Kind  string // ShowSlowKind
}

// ParseShowSlow 解析Show Slow
func (ss *ShowSlow) ParseShowSlow(ass *ast.ShowSlow) *ShowSlow {
	if ass == nil {
		return nil
	}

	// Type  ShowSlowType
	if val, bv := classifyShowSlowType(ass.Tp); bv == true {
		ss.Type = val
	}
	// Count uint64
	ss.Count = strconv.FormatUint(ass.Count, 10)

	// Kind  ShowSlowKind
	if val, bv := classifyShowSlowKind(ass.Kind); bv == true {
		ss.Kind = val
	}

	return ss
}

// ShowStmtType like *ast.ShowStmtType
type ShowStmtType int
