package types

import (
	"fmt"

	"github.com/pingcap/parser/ast"
)

// SQLTree SQL语法树
type SQLTree struct {
	StmtType string
	StmtTree interface{}
}

// ClassifyStmtNode 对StmtNode进行SQL分类&解析
func (st *SQLTree) ClassifyStmtNode(StmtNode ast.StmtNode) {
	if StmtNode == nil {
		return
	}

	switch node := StmtNode.(type) {
	// Select
	case *ast.SelectStmt:
		st.StmtType = "Select"
		s := new(SelectStmt)
		st.StmtTree = s.ParseSelectStmt(node)
	// DML
	case *ast.InsertStmt:
		st.StmtType = "Insert"
		i := new(InsertStmt)
		st.StmtTree = i.ParseInsertStmt(node)
	case *ast.UpdateStmt:
		st.StmtType = "Update"
		u := new(UpdateStmt)
		st.StmtTree = u.ParseUpdateStmt(node)
	case *ast.DeleteStmt:
		st.StmtType = "Delete"
		d := new(DeleteStmt)
		st.StmtTree = d.ParseDeleteStmt(node)
	// DDL
	//  CREATE
	case *ast.CreateDatabaseStmt:
		st.StmtType = "CreateDatabase"
		cd := new(CreateDatabaseStmt)
		st.StmtTree = cd.ParseCreateDatabaseStmt(node)
	case *ast.CreateUserStmt:
		st.StmtType = "CreateUser"
		cu := new(CreateUserStmt)
		st.StmtTree = cu.ParseCreateUserStmt(node)
	case *ast.CreateTableStmt:
		st.StmtType = "CreateTable"
		ct := new(CreateTableStmt)
		st.StmtTree = ct.ParseCreateTableStmt(node)
	case *ast.CreateIndexStmt:
		st.StmtType = "CreateIndex"
		ci := new(CreateIndexStmt)
		st.StmtTree = ci.ParseCreateIndexStmt(node)
	case *ast.CreateViewStmt:
		st.StmtType = "CreateView"
		cv := new(CreateViewStmt)
		st.StmtTree = cv.ParseCreateViewStmt(node)
	case *ast.CreateBindingStmt:
		st.StmtType = "CreateBinding"
		cb := new(CreateBindingStmt)
		st.StmtTree = cb.ParseCreateBindingStmt(node)
	case *ast.CreateSequenceStmt:
		st.StmtType = "CreateSequence"
		cs := new(CreateSequenceStmt)
		st.StmtTree = cs.ParseCreateSequenceStmt(node)
	//case *ast.CreateStatisticsStmt:
	//  ALTER
	case *ast.AlterDatabaseStmt:
		st.StmtType = "AlterDatabase"
		ad := new(AlterDatabaseStmt)
		st.StmtTree = ad.ParseAlterDatabaseStmt(node)
	case *ast.AlterInstanceStmt:
		st.StmtType = "AlterInstance"
		ai := new(AlterInstanceStmt)
		st.StmtTree = ai.ParseAlterInstanceStmt(node)
	case *ast.AlterUserStmt:
		st.StmtType = "AlterUser"
		au := new(AlterUserStmt)
		st.StmtTree = au.ParseAlterUserStmt(node)
	case *ast.AlterTableStmt:
		st.StmtType = "AlterTable"
		at := new(AlterTableStmt)
		st.StmtTree = at.ParseAlterTableStmt(node)
	//case *ast.AlterSequenceStmt:
	//  DROP
	case *ast.DropDatabaseStmt:
		st.StmtType = "DropDatabase"
		dd := new(DropDatabaseStmt)
		st.StmtTree = dd.ParseDropDatabaseStmt(node)
	case *ast.DropUserStmt:
		st.StmtType = "DropUser"
		du := new(DropUserStmt)
		st.StmtTree = du.ParseDropUserStmt(node)
	case *ast.DropTableStmt:
		st.StmtType = "DropTable"
		dt := new(DropTableStmt)
		st.StmtTree = dt.ParseDropTableStmt(node)
	case *ast.DropIndexStmt:
		st.StmtType = "DropIndex"
		di := new(DropIndexStmt)
		st.StmtTree = di.ParseDropIndexStmt(node)
	case *ast.DropBindingStmt:
		st.StmtType = "DropBinding"
		db := new(DropBindingStmt)
		st.StmtTree = db.ParseDropBindingStmt(node)
	case *ast.DropSequenceStmt:
		st.StmtType = "DropSequence"
		ds := new(DropSequenceStmt)
		st.StmtTree = ds.ParseDropSequenceStmt(node)
	case *ast.DropStatsStmt:
		st.StmtType = "DropStats"
		ds := new(DropStatsStmt)
		st.StmtTree = ds.ParseDropStatsStmt(node)
	//case *ast.DropStatisticsStmt:
	//  TRUNCATE
	case *ast.TruncateTableStmt:
		st.StmtType = "TruncateTable"
		tt := new(TruncateTableStmt)
		st.StmtTree = tt.ParseTruncateTableStmt(node)
	//  RENAME
	case *ast.RenameTableStmt:
		st.StmtType = "RenameTable"
		rt := new(RenameTableStmt)
		st.StmtTree = rt.ParseRenameTableStmt(node)
	// DCL
	case *ast.GrantStmt:
		st.StmtType = "Grant"
		g := new(GrantStmt)
		st.StmtTree = g.ParseGrantStmt(node)
	case *ast.GrantRoleStmt:
		st.StmtType = "GrantRole"
		gr := new(GrantRoleStmt)
		st.StmtTree = gr.ParseGrantRoleStmt(node)
	case *ast.RevokeStmt:
		st.StmtType = "Revoke"
		r := new(RevokeStmt)
		st.StmtTree = r.ParseRevokeStmt(node)
	case *ast.RevokeRoleStmt:
		st.StmtType = "RevokeRole"
		rr := new(RevokeRoleStmt)
		st.StmtTree = rr.ParseRevokeRoleStmt(node)
	// Other
	case *ast.UseStmt:
		st.StmtType = "Use"
		u := new(UseStmt)
		st.StmtTree = u.ParseUseStmt(node)
	case *ast.ShowStmt:
		st.StmtType = "Show"
		show := new(ShowStmt)
		st.StmtTree = show.ParseShowStmt(node)
	case *ast.BeginStmt:
		st.StmtType = "Begin"
		beg := new(BeginStmt)
		st.StmtTree = beg.ParseBeginStmt(node)
	case *ast.CommitStmt:
		st.StmtType = "Commit"
		cm := new(CommitStmt)
		st.StmtTree = cm.ParseCommitStmt(node)
	case *ast.RollbackStmt:
		st.StmtType = "Rollback"
		rb := new(RollbackStmt)
		st.StmtTree = rb.ParseRollbackStmt(node)
	case *ast.PrepareStmt:
		st.StmtType = "Prepare"
		pp := new(PrepareStmt)
		st.StmtTree = pp.ParsePrepareStmt(node)
	case *ast.ExecuteStmt:
		st.StmtType = "Execute"
		exe := new(ExecuteStmt)
		st.StmtTree = exe.ParseExecuteStmt(node)
	case *ast.ExplainStmt:
		st.StmtType = "Explain"
		ep := new(ExplainStmt)
		st.StmtTree = ep.ParseExplainStmt(node)
	case *ast.ExplainForStmt:
		st.StmtType = "ExplainFor"
		epf := new(ExplainForStmt)
		st.StmtTree = epf.ParseExplainForStmt(node)
	case *ast.FlushStmt:
		st.StmtType = "Flush"
		flu := new(FlushStmt)
		st.StmtTree = flu.ParseFlushStmt(node)
	case *ast.BinlogStmt: // BinlogStmt is an internal-use statement.We just Parse and ignore it.
		st.StmtType = "Binlog"
		bin := new(BinlogStmt)
		st.StmtTree = bin.ParseBinlogStmt(node)
	//case *ast.CallStmt:
	case *ast.ChangeStmt:
		st.StmtType = "Change"
		chg := new(ChangeStmt)
		st.StmtTree = chg.ParseChangeStmt(node)
	case *ast.CleanupTableLockStmt:
		st.StmtType = "CleanupTableLock"
		ctl := new(CleanupTableLockStmt)
		st.StmtTree = ctl.ParseCleanupTableLockStmt(node)
	case *ast.DeallocateStmt:
		st.StmtType = "Deallocate"
		dlc := new(DeallocateStmt)
		st.StmtTree = dlc.ParseDeallocateStmt(node)
	case *ast.AdminStmt:
		st.StmtType = "Admin"
		adm := new(AdminStmt)
		st.StmtTree = adm.ParseAdminStmt(node)
	case *ast.DoStmt:
		st.StmtType = "Do"
		do := new(DoStmt)
		st.StmtTree = do.ParseDoStmt(node)
	case *ast.IndexAdviseStmt:
		st.StmtType = "IndexAdvise"
		ia := new(IndexAdviseStmt)
		st.StmtTree = ia.ParseIndexAdviseStmt(node)
	case *ast.KillStmt:
		st.StmtType = "Kill"
		kl := new(KillStmt)
		st.StmtTree = kl.ParseKillStmt(node)
	case *ast.LoadDataStmt:
		st.StmtType = "LoadData"
		ld := new(LoadDataStmt)
		st.StmtTree = ld.ParseLoadDataStmt(node)
	case *ast.LoadStatsStmt:
		st.StmtType = "LoadStats"
		ls := new(LoadStatsStmt)
		st.StmtTree = ls.ParseLoadStatsStmt(node)
	case *ast.LockTablesStmt:
		st.StmtType = "LockTables"
		lt := new(LockTablesStmt)
		st.StmtTree = lt.ParseLockTablesStmt(node)
	case *ast.UnlockTablesStmt:
		st.StmtType = "UnlockTables"
		st.StmtTree = ""
	//case *ast.PurgeImportStmt:
	case *ast.RecoverTableStmt:
		st.StmtType = "RecoverTable"
		rct := new(RecoverTableStmt)
		st.StmtTree = rct.ParseRecoverTableStmt(node)
	case *ast.RepairTableStmt:
		st.StmtType = "RepairTable"
		rpt := new(RepairTableStmt)
		st.StmtTree = rpt.ParseRepairTableStmt(node)
	//  SET
	case *ast.SetStmt:
		st.StmtType = "Set"
		set := new(SetStmt)
		st.StmtTree = set.ParseSetStmt(node)
	case *ast.SetConfigStmt: // TiDB, TiKV, PD
		st.StmtType = "SetConfig"
		scf := new(SetConfigStmt)
		st.StmtTree = scf.ParseSetConfigStmt(node)
	case *ast.SetDefaultRoleStmt:
		st.StmtType = "SetDefaultRole"
		sdr := new(SetDefaultRoleStmt)
		st.StmtTree = sdr.ParseSetDefaultRoleStmt(node)
	case *ast.SetPwdStmt:
		st.StmtType = "SetPassword"
		spw := new(SetPwdStmt)
		st.StmtTree = spw.ParseSetPwdStmt(node)
	case *ast.SetRoleStmt:
		st.StmtType = "SetRole"
		sro := new(SetRoleStmt)
		st.StmtTree = sro.ParseSetRoleStmt(node)
	//case *ast.SetOprStmt:
	case *ast.ShutdownStmt:
		st.StmtType = "Shutdown"
		st.StmtTree = ""
	case *ast.SplitRegionStmt:
		st.StmtType = "SplitRegion"
		sre := new(SplitRegionStmt)
		st.StmtTree = sre.ParseSplitRegionStmt(node)
	case *ast.TraceStmt:
		st.StmtType = "Trace"
		trc := new(TraceStmt)
		st.StmtTree = trc.ParseTraceStmt(node)
	default:
		fmt.Println("not match")
	}
}
