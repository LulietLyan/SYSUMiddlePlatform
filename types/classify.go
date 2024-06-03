package types

import (
	"strconv"

	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/model"
	"github.com/pingcap/parser/mysql"
	"github.com/pingcap/parser/opcode"
)

// classifyResultSetNode 对ResultSetNode进行节点分类&解析
// ResultSetNode接口具有ResultFields属性，实现SelectStmt，SubqueryExpr，TableSource，TableName和Join
func classifyResultSetNode(ResNode ast.ResultSetNode) (interface{}, bool) {
	if ResNode == nil {
		return nil, false
	}

	switch node := ResNode.(type) {
	// SubSelect
	case *ast.SelectStmt:
		se := new(SelectStmt)
		return se.ParseSelectStmt(node), true
	// Join
	case *ast.Join:
		j := new(Join)
		return j.ParseJoin(node), true
	// TableSource
	case *ast.TableSource:
		ts := new(TableSource)
		return ts.ParseTableSource(node), true
	// Union
	case *ast.UnionStmt:
		// u := new(Union)
		// return u.ParseUnion(node)
		return node, true
	// TableName
	case *ast.TableName:
		tn := new(TableName)
		return tn.ParseTableName(node), true
	// SubqueryExpr
	case *ast.SubqueryExpr:
		return node, true
	default:
		return nil, false
	}
}

// classifyExprNode 解析表达式节点
func classifyExprNode(ExprNode ast.ExprNode) (interface{}, bool) {
	if ExprNode == nil {
		return nil, false
	}

	switch node := ExprNode.(type) {
	case *ast.AggregateFuncExpr:
		afe := new(AggregateFuncExpr)
		afe.ParseAggregateFuncExpr(node)
		return afe, true
	case *ast.BetweenExpr:
		bwe := new(BetweenExpr)
		bwe.ParseBetweenExpr(node)
		return bwe, true
	case *ast.BinaryOperationExpr: // where条件
		boe := new(BinaryOperationExpr)
		boe.ParseBinaryOperationExpr(node)
		return boe, true
	case *ast.CaseExpr:
		ce := new(CaseExpr)
		ce.ParseCaseExpr(node)
		return ce, true
	case *ast.ColumnNameExpr:
		cne := new(ColumnNameExpr)
		cne.ParseColumnNameExpr(node)
		return cne, true
	case *ast.CompareSubqueryExpr:
		csqe := new(CompareSubqueryExpr)
		csqe.ParseCompareSubqueryExpr(node)
		return csqe, true
	case *ast.DefaultExpr:
		de := new(DefaultExpr)
		de.ParseDefaultExpr(node)
		return de, true
	case *ast.ExistsSubqueryExpr:
		esqe := new(ExistsSubqueryExpr)
		esqe.ParseExistsSubqueryExpr(node)
		return esqe, true
	case *ast.FuncCallExpr:
		fce := new(FuncCallExpr)
		fce.ParseFuncCallExpr(node)
		return fce, true
	case *ast.FuncCastExpr:
		fcae := new(FuncCastExpr)
		fcae.ParseFuncCastExpr(node)
		return fcae, true
	case *ast.IsNullExpr:
		ine := new(IsNullExpr)
		ine.ParseIsNullExpr(node)
		return ine, true
	case *ast.IsTruthExpr:
		ite := new(IsTruthExpr)
		ite.ParseIsTruthExpr(node)
		return ite, true
	case *ast.MaxValueExpr:
		return "MAXVALUE", true
	case ast.ParamMarkerExpr:
		return node.GetValue(), true
	case *ast.ParenthesesExpr:
		pte := new(ParenthesesExpr)
		pte.ParseParenthesesExpr(node)
		return pte, true
	case *ast.PatternInExpr:
		pie := new(PatternInExpr)
		pie.ParsePatternInExpr(node)
		return pie, true
	case *ast.PatternLikeExpr:
		ple := new(PatternLikeExpr)
		ple.ParsePatternLikeExpr(node)
		return ple, true
	case *ast.PatternRegexpExpr:
		pre := new(PatternRegexpExpr)
		pre.ParsePatternRegexpExpr(node)
		return pre, true
	case *ast.PositionExpr:
		pe := new(PositionExpr)
		pe.ParsePositionExpr(node)
		return nil, false
	case *ast.RowExpr:
		re := new(RowExpr)
		re.ParseRowExpr(node)
		return re, true
	case *ast.SetCollationExpr:
		sce := new(SetCollationExpr)
		sce.ParseSetCollationExpr(node)
		return sce, true
	case *ast.SubqueryExpr:
		sqe := new(SubqueryExpr)
		sqe.ParseSubqueryExpr(node)
		return sqe, true
	case *ast.TableNameExpr:
		tne := new(TableNameExpr)
		tne.ParseTableNameExpr(node)
		return tne, true
	case *ast.TimeUnitExpr:
		tue := new(TimeUnitExpr)
		tue.ParseTimeUnitExpr(node)
		return tue, true
	case *ast.TrimDirectionExpr:
		tde := new(TrimDirectionExpr)
		tde.ParseTrimDirectionExpr(node)
		return tde, true
	case *ast.UnaryOperationExpr:
		uoe := new(UnaryOperationExpr)
		uoe.ParseUnaryOperationExpr(node)
		return uoe, true
	case ast.ValueExpr:
		return node.GetValue(), true
	case *ast.ValuesExpr:
		vse := new(ValuesExpr)
		vse.ParseValuesExpr(node)
		return vse, true
	case *ast.VariableExpr:
		ve := new(VariableExpr)
		ve.ParseVariableExpr(node)
		return ve, true
	case *ast.WindowFuncExpr:
		wfe := new(WindowFuncExpr)
		wfe.ParseWindowFuncExpr(node)
		return wfe, true
	default:
		return nil, false
	}
}

// classifyMySQLColDataType 解析Column Data类型
func classifyColDataType(cdt byte) (string, bool) {
	// MySQL type information
	const (
		TypeDecimal   byte = 0
		TypeTiny      byte = 1
		TypeShort     byte = 2
		TypeLong      byte = 3
		TypeFloat     byte = 4
		TypeDouble    byte = 5
		TypeNull      byte = 6
		TypeTimestamp byte = 7
		TypeLonglong  byte = 8
		TypeInt24     byte = 9
		TypeDate      byte = 10
		/* TypeDuration original name was TypeTime, renamed to TypeDuration to resolve the conflict with Go type Time.*/
		TypeDuration byte = 11
		TypeDatetime byte = 12
		TypeYear     byte = 13
		TypeNewDate  byte = 14
		TypeVarchar  byte = 15
		TypeBit      byte = 16

		TypeJSON       byte = 0xf5
		TypeNewDecimal byte = 0xf6
		TypeEnum       byte = 0xf7
		TypeSet        byte = 0xf8
		TypeTinyBlob   byte = 0xf9
		TypeMediumBlob byte = 0xfa
		TypeLongBlob   byte = 0xfb
		TypeBlob       byte = 0xfc
		TypeVarString  byte = 0xfd
		TypeString     byte = 0xfe
		TypeGeometry   byte = 0xff
	)

	switch {
	case cdt == TypeVarchar:
		return "varchar", true
	case cdt == TypeDecimal:
		return "decimal", true
	case cdt == TypeTiny:
		return "tinyint", true
	case cdt == TypeShort:
		return "smallint", true
	case cdt == TypeInt24:
		return "mediumint", true
	case cdt == TypeLong:
		return "int", true
	case cdt == TypeLonglong:
		return "bigint", true
	case cdt == TypeFloat:
		return "float", true
	case cdt == TypeDouble:
		return "double", true
	case cdt == TypeDate:
		return "date", true
	case cdt == TypeDuration:
		return "time", true
	case cdt == TypeDatetime:
		return "datetime", true
	case cdt == TypeTimestamp:
		return "timestamp", true
	case cdt == TypeYear:
		return "year", true
	case cdt == TypeBit:
		return "bit", true
	case cdt == TypeJSON:
		return "json", true
	case cdt == TypeTinyBlob:
		return "tinytext/tinyblob", true
	case cdt == TypeMediumBlob:
		return "mediumtext/mediumblob", true
	case cdt == TypeLongBlob:
		return "longtext/longblob", true
	case cdt == TypeBlob:
		return "text/blob", true
	case cdt == TypeEnum:
		return "enum", true
	case cdt == TypeSet:
		return "set", true
	default:
		return strconv.Itoa(int(cdt)), true
	}
}

// classifyColumnOptionType 解析Column Option类型
func classifyColumnOptionType(aco *ast.ColumnOption) (string, bool) {
	switch aco.Tp {
	case ast.ColumnOptionAutoIncrement:
		return "AutoIncrement", true
	case ast.ColumnOptionPrimaryKey:
		return "PrimaryKey", true
	case ast.ColumnOptionUniqKey:
		return "UniqueKey", true
	case ast.ColumnOptionNotNull:
		return "NotNull", true
	case ast.ColumnOptionNull:
		return "NotNull", false
	case ast.ColumnOptionDefaultValue:
		return "DefaultValue", true
	case ast.ColumnOptionComment:
		return "Comment", true
	case ast.ColumnOptionOnUpdate:
		return "OnUpdate", true
	case ast.ColumnOptionFulltext:
		return "Fulltext", true
	case ast.ColumnOptionGenerated:
		return "Generated", true
	case ast.ColumnOptionReference:
		return "Reference", true
	case ast.ColumnOptionCollate:
		return "Collate", true
	case ast.ColumnOptionCheck:
		return "Check", true
	case ast.ColumnOptionColumnFormat:
		return "ColumnFormat", true
	case ast.ColumnOptionStorage:
		return "Storage", true
	case ast.ColumnOptionAutoRandom:
		return "AutoRandom", true
	default:
		return "", false
	}
}

// classifyColumnPositionType 解析Column Position类型
func classifyColumnPositionType(acp *ast.ColumnPosition) (string, bool) {
	switch acp.Tp {
	case ast.ColumnPositionNone:
		return "None", false
	case ast.ColumnPositionFirst:
		return "First", false
	case ast.ColumnPositionAfter:
		return "After", false
	default:
		return "", false
	}
}

// classifyReferOptionType 解析Reference Option类型
func classifyReferOptionType(rot ast.ReferOptionType) (string, bool) {
	switch rot {
	case ast.ReferOptionNoOption:
		return "NoOption", true
	case ast.ReferOptionRestrict:
		return "Restrict", true
	case ast.ReferOptionCascade:
		return "Cascade", true
	case ast.ReferOptionSetNull:
		return "SetNull", true
	case ast.ReferOptionNoAction:
		return "NoAction", true
	case ast.ReferOptionSetDefault:
		return "SetDefault", true
	default:
		return "", false
	}
}

// classifyMatchType 解析Match的类型
func classifyMatchType(mt ast.MatchType) (string, bool) {
	switch mt {
	case ast.MatchNone:
		return "None", true
	case ast.MatchFull:
		return "Full", true
	case ast.MatchPartial:
		return "Partial", true
	case ast.MatchSimple:
		return "Simple", true
	default:
		return "", false
	}
}

// classifyConstraintType 解析Constraint类型
func classifyConstraintType(acst *ast.Constraint) (string, bool) {
	switch acst.Tp {
	case ast.ConstraintNoConstraint:
		return "hasNoConstraint", true
	case ast.ConstraintPrimaryKey:
		return "hasPrimaryKey", true
	case ast.ConstraintKey:
		return "hasKey", true
	case ast.ConstraintIndex:
		return "hasIndex", true
	case ast.ConstraintUniq:
		return "hasUniq", true
	case ast.ConstraintUniqKey:
		return "hasUniqKey", true
	case ast.ConstraintUniqIndex:
		return "hasUniqIndex", true
	case ast.ConstraintForeignKey:
		return "hasForeignKey", true
	case ast.ConstraintFulltext:
		return "hasFulltext", true
	case ast.ConstraintCheck:
		return "hasCheck", true
	default:
		return "", false
	}
}

// classifySequenceOptionType 解析Sequence Option的类型
func classifySequenceOptionType(asot ast.SequenceOptionType) (string, bool) {
	switch asot {
	case ast.SequenceOptionNone:
		return "None", true
	case ast.SequenceOptionIncrementBy:
		return "IncrementBy", true
	case ast.SequenceStartWith:
		return "StartWith", true
	case ast.SequenceNoMinValue:
		return "NoMinValue", true
	case ast.SequenceMinValue:
		return "MinValue", true
	case ast.SequenceNoMaxValue:
		return "NoMaxValue", true
	case ast.SequenceMaxValue:
		return "MaxValue", true
	case ast.SequenceNoCache:
		return "NoCache", true
	case ast.SequenceCache:
		return "Cache", true
	case ast.SequenceNoCycle:
		return "NoCycle", true
	case ast.SequenceCycle:
		return "Cycle", true
	// SequenceRestart is only used in alter sequence statement.
	// case ast.SequenceRestart:
	// 	return "Restart", true
	// case ast.SequenceRestartWith:
	// 	return "RestartWith", true
	default:
		return "", false
	}
}

// classifyDatabaseOptionType 解析Database Option类型
func classifyDatabaseOptionType(ado *ast.DatabaseOption) (string, bool) {
	switch ado.Tp {
	case ast.DatabaseOptionNone:
		return "None", true
	case ast.DatabaseOptionCharset:
		return "Charset", true
	case ast.DatabaseOptionCollate:
		return "Collate", true
	case ast.DatabaseOptionEncryption:
		return "Encryption", true
	default:
		return "", false
	}
}

// classifyTableOption 解析MySQL Table Option类型
func classifyTableOptionType(ato *ast.TableOption) (string, bool) {
	switch ato.Tp {
	case ast.TableOptionEngine:
		return "Engine", true
	case ast.TableOptionCharset:
		return "Charset", true
	case ast.TableOptionCollate:
		return "Collate", true
	case ast.TableOptionAutoIdCache:
		return "AutoIdCache", true
	case ast.TableOptionAutoIncrement:
		return "AutoIncrement", true
	case ast.TableOptionAutoRandomBase:
		return "AutoRandomBase", true
	case ast.TableOptionComment:
		return "Comment", true
	case ast.TableOptionAvgRowLength:
		return "AvgRowLength", true
	case ast.TableOptionCheckSum:
		return "CheckSum", true
	case ast.TableOptionCompression:
		return "Compression", true
	case ast.TableOptionConnection:
		return "Connection", true
	case ast.TableOptionPassword:
		return "Password", true
	case ast.TableOptionKeyBlockSize:
		return "KeyBlockSize", true
	case ast.TableOptionMaxRows:
		return "MaxRows", true
	case ast.TableOptionMinRows:
		return "MinRows", true
	case ast.TableOptionDelayKeyWrite:
		return "DelayKeyWrite", true
	case ast.TableOptionRowFormat:
		return "RowFormat", true
	case ast.TableOptionStatsPersistent:
		return "StatsPersistent", true
	case ast.TableOptionStatsAutoRecalc:
		return "StatsAutoRecalc", true
	case ast.TableOptionShardRowID:
		return "ShardRowID", true
	case ast.TableOptionPreSplitRegion:
		return "PreSplitRegion", true
	case ast.TableOptionPackKeys:
		return "PackKeys", true
	case ast.TableOptionTablespace:
		return "Tablespace", true
	case ast.TableOptionNodegroup:
		return "Nodegroup", true
	case ast.TableOptionDataDirectory:
		return "DataDirectory", true
	case ast.TableOptionIndexDirectory:
		return "IndexDirectory", true
	case ast.TableOptionStorageMedia:
		return "StorageMedia", true
	case ast.TableOptionStatsSamplePages:
		return "StatsSamplePages", true
	case ast.TableOptionSecondaryEngine:
		return "SecondaryEngine", true
	case ast.TableOptionSecondaryEngineNull:
		return "SecondaryEngineNull", true
	case ast.TableOptionInsertMethod:
		return "InsertMethod", true
	case ast.TableOptionTableCheckSum:
		return "TableCheckSum", true
	case ast.TableOptionUnion:
		return "Union", true
	case ast.TableOptionEncryption:
		return "Encryption", true
	default:
		return "", false
	}
}

// classifyTableLockType 解析MySQL Table Lock类型
func classifyTableLockType(mtlt model.TableLockType) (string, bool) {
	switch mtlt {
	case model.TableLockNone:
		return "None", true
	case model.TableLockRead:
		return "Read", true
	case model.TableLockReadLocal:
		return "ReadLocal", true
	//case model.TableLockReadOnly:
	case model.TableLockWrite:
		return "Write", true
	case model.TableLockWriteLocal:
		return "WriteLocal", true
	default:
		return "", false
	}
}

// classifyOnDuplicateKeyHandlingType 解析MySQL重复值的处理方式，适用于'CREATE TABLE ... SELECT' 或 'LOAD DATA'
func classifyOnDuplicateKeyHandlingType(aodkht ast.OnDuplicateKeyHandlingType) (string, bool) {
	switch aodkht {
	case ast.OnDuplicateKeyHandlingError:
		return "ERROR", true
	case ast.OnDuplicateKeyHandlingIgnore:
		return "IGNORE", true
	case ast.OnDuplicateKeyHandlingReplace:
		return "REPLACE", true
	default:
		return "", false
	}
}

// classifyJoinType 解析Join类型
func classifyJoinType(j *ast.Join) (string, bool) {
	switch j.Tp {
	case ast.CrossJoin:
		return "CrossJoin", true
	case ast.LeftJoin:
		return "LeftJoin", true
	case ast.RightJoin:
		return "RightJoin", true
	default:
		return "Single", true
	}
}

// classifyAlterTableType 解析alter table的类型
func classifyAlterTableType(ats *ast.AlterTableSpec) (string, bool) {
	switch ats.Tp {
	case ast.AlterTableOption:
		return "Option", true
	case ast.AlterTableAddColumns:
		return "AddColumns", true
	case ast.AlterTableAddConstraint:
		return "AddConstraint", true
	case ast.AlterTableDropColumn:
		return "DropColumn", true
	case ast.AlterTableDropPrimaryKey:
		return "DropPrimaryKey", true
	case ast.AlterTableDropIndex:
		return "DropIndex", true
	case ast.AlterTableDropForeignKey:
		return "DropForeignKey", true
	case ast.AlterTableModifyColumn:
		return "ModifyColumn", true
	case ast.AlterTableChangeColumn:
		return "ChangeColumn", true
	case ast.AlterTableRenameColumn:
		return "RenameColumn", true
	case ast.AlterTableRenameTable:
		return "RenameTable", true
	case ast.AlterTableAlterColumn:
		return "AlterColumn", true
	case ast.AlterTableLock:
		return "Lock", true
	//case ast.AlterTableWriteable:
	case ast.AlterTableAlgorithm:
		return "Algorithm", true
	case ast.AlterTableRenameIndex:
		return "RenameIndex", true
	case ast.AlterTableForce:
		return "Force", true
	case ast.AlterTableAddPartitions:
		return "AddPartitions", true
	//case ast.AlterTableAlterPartition:
	case ast.AlterTableCoalescePartitions:
		return "CoalescePartitions", true
	case ast.AlterTableDropPartition:
		return "DropPartition", true
	case ast.AlterTableTruncatePartition:
		return "TruncatePartition", true
	case ast.AlterTablePartition:
		return "Partition", true
	case ast.AlterTableEnableKeys:
		return "EnableKeys", true
	case ast.AlterTableDisableKeys:
		return "DisableKeys", true
	case ast.AlterTableRemovePartitioning:
		return "RemovePartitioning", true
	case ast.AlterTableWithValidation:
		return "Validation", true
	case ast.AlterTableWithoutValidation:
		return "WithoutValidation", true
	case ast.AlterTableSecondaryLoad:
		return "SecondaryLoad", true
	case ast.AlterTableSecondaryUnload:
		return "SecondaryUnload", true
	case ast.AlterTableRebuildPartition:
		return "RebuildPartition", true
	case ast.AlterTableReorganizePartition:
		return "ReorganizePartition", true
	case ast.AlterTableCheckPartitions:
		return "CheckPartitions", true
	case ast.AlterTableExchangePartition:
		return "ExchangePartition", true
	case ast.AlterTableOptimizePartition:
		return "OptimizePartition", true
	case ast.AlterTableRepairPartition:
		return "RepairPartition", true
	case ast.AlterTableImportPartitionTablespace:
		return "ImportPartitionTablespace", true
	case ast.AlterTableDiscardPartitionTablespace:
		return "DiscardPartitionTablespace", true
	case ast.AlterTableAlterCheck:
		return "AlterCheck", true
	case ast.AlterTableDropCheck:
		return "DropCheck", true
	case ast.AlterTableImportTablespace:
		return "ImportTablespace", true
	case ast.AlterTableDiscardTablespace:
		return "DiscardTablespace", true
	case ast.AlterTableIndexInvisible:
		return "IndexInvisible", true
	// TODO: Add more actions
	case ast.AlterTableOrderByColumns:
		return "OrderByColumns", true
	// AlterTableSetTiFlashReplica uses to set the table TiFlash replica.
	case ast.AlterTableSetTiFlashReplica:
		return "SetTiFlashReplica", true
	//case ast.AlterTablePlacement:
	default:
		return "", false
	}
}

// classifyPriorityEnum 解析Priority
func classifyPriorityEnum(mpe mysql.PriorityEnum) (string, bool) {
	switch mpe {
	case mysql.NoPriority:
		return "NoPriority", true
	case mysql.LowPriority:
		return "LowPriority", true
	case mysql.HighPriority:
		return "HighPriority", true
	case mysql.DelayedPriority:
		return "DelayedPriority", true
	default:
		return "", false
	}
}

// classifyLockType 解析Lock类型
func classifyLockType(lt ast.LockType) (string, bool) {
	switch lt {
	case ast.LockTypeNone:
		return "None", true
	case ast.LockTypeDefault:
		return "Default", true
	case ast.LockTypeShared:
		return "Shared", true
	case ast.LockTypeExclusive:
		return "Exclusive", true
	default:
		return "", false
	}
}

// classifySelectLockType 解析Select Lock类型
func classifySelectLockType(slt ast.SelectLockType) (string, bool) {
	switch slt {
	case ast.SelectLockNone:
		return "None", true
	case ast.SelectLockForUpdate:
		return "ForUpdate", true
	case ast.SelectLockInShareMode:
		return "Share", true
	case ast.SelectLockForUpdateNoWait:
		return "ForUpdateNoWait", true
	//case ast.SelectLockForUpdateWaitN:
	default:
		return "", false
	}
}

// classifyAlgorithmType 解析Algorithm类型
func classifyAlgorithmType(at ast.AlgorithmType) (string, bool) {
	switch at {
	case ast.AlgorithmTypeDefault:
		return "Default", true
	case ast.AlgorithmTypeCopy:
		return "Copy", true
	case ast.AlgorithmTypeInplace:
		return "Inplace", true
	case ast.AlgorithmTypeInstant:
		return "Instant", true
	default:
		return "", false
	}
}

// classifyTimeUnitType 解析时间单位类型
func classifyTimeUnitType(tut ast.TimeUnitType) (string, bool) {
	switch tut {
	// TimeUnitInvalid is a placeholder for an invalid time or timestamp unit
	case ast.TimeUnitInvalid:
		return "INVALID", true
	// TimeUnitMicrosecond is the time or timestamp unit MICROSECOND.
	case ast.TimeUnitMicrosecond:
		return "MICROSECOND", true
	// TimeUnitSecond is the time or timestamp unit SECOND.
	case ast.TimeUnitSecond:
		return "SECOND", true
	// TimeUnitMinute is the time or timestamp unit MINUTE.
	case ast.TimeUnitMinute:
		return "MINUTE", true
	// TimeUnitHour is the time or timestamp unit HOUR.
	case ast.TimeUnitHour:
		return "HOUR", true
	// TimeUnitDay is the time or timestamp unit DAY.
	case ast.TimeUnitDay:
		return "DAY", true
	// TimeUnitWeek is the time or timestamp unit WEEK.
	case ast.TimeUnitWeek:
		return "WEEK", true
	// TimeUnitMonth is the time or timestamp unit MONTH.
	case ast.TimeUnitMonth:
		return "MONTH", true
	// TimeUnitQuarter is the time or timestamp unit QUARTER.
	case ast.TimeUnitQuarter:
		return "QUARTER", true
	// TimeUnitYear is the time or timestamp unit YEAR.
	case ast.TimeUnitYear:
		return "YEAR", true
	// TimeUnitSecondMicrosecond is the time unit SECOND_MICROSECOND.
	case ast.TimeUnitSecondMicrosecond:
		return "SECOND_MICROSECOND", true
	// TimeUnitMinuteMicrosecond is the time unit MINUTE_MICROSECOND.
	case ast.TimeUnitMinuteMicrosecond:
		return "MINUTE_MICROSECOND", true
	// TimeUnitMinuteSecond is the time unit MINUTE_SECOND.
	case ast.TimeUnitMinuteSecond:
		return "MINUTE_SECOND", true
	// TimeUnitHourMicrosecond is the time unit HOUR_MICROSECOND.
	case ast.TimeUnitHourMicrosecond:
		return "HOUR_MICROSECOND", true
	// TimeUnitHourSecond is the time unit HOUR_SECOND.
	case ast.TimeUnitHourSecond:
		return "HOUR_SECOND", true
	// TimeUnitHourMinute is the time unit HOUR_MINUTE.
	case ast.TimeUnitHourMinute:
		return "HOUR_MINUTE", true
	// TimeUnitDayMicrosecond is the time unit DAY_MICROSECOND.
	case ast.TimeUnitDayMicrosecond:
		return "DAY_MICROSECOND", true
	// TimeUnitDaySecond is the time unit DAY_SECOND.
	case ast.TimeUnitDaySecond:
		return "DAY_SECOND", true
	// TimeUnitDayMinute is the time unit DAY_MINUTE.
	case ast.TimeUnitDayMinute:
		return "DAY_MINUTE", true
	// TimeUnitDayHour is the time unit DAY_HOUR.
	case ast.TimeUnitDayHour:
		return "DAY_HOUR", true
	// TimeUnitYearMonth is the time unit YEAR_MONTH.
	case ast.TimeUnitYearMonth:
		return "YEAR_MONTH", true
	default:
		return "", false
	}
}

// classifyTrimDirectionType 解析TrimDirection类型
func classifyTrimDirectionType(tdt ast.TrimDirectionType) (string, bool) {
	switch tdt {
	// TrimBothDefault trims from both direction by default.
	case ast.TrimBothDefault:
		return "Default", true
	// TrimBoth trims from both direction with explicit notation.
	case ast.TrimBoth:
		return "Both", true
	// TrimLeading trims from left.
	case ast.TrimLeading:
		return "Leading", true
	// TrimTrailing trims from right.
	case ast.TrimTrailing:
		return "Trailing", true
	default:
		return "", false
	}
}

// classifyOpcode 解析Opcode类型
func classifyOpcode(opt opcode.Op) (string, bool) {
	switch opt {
	case opcode.LogicAnd:
		return "LogicAnd", true
	case opcode.LeftShift:
		return "LeftShift", true
	case opcode.RightShift:
		return "RightShift", true
	case opcode.LogicOr:
		return "LogicOr", true
	case opcode.GE:
		return "GE", true
	case opcode.LE:
		return "LE", true
	case opcode.EQ:
		return "EQ", true
	case opcode.NE:
		return "NE", true
	case opcode.LT:
		return "LT", true
	case opcode.GT:
		return "GT", true
	case opcode.Plus:
		return "Plus", true
	case opcode.Minus:
		return "Minus", true
	case opcode.And:
		return "And", true
	case opcode.Or:
		return "Or", true
	case opcode.Mod:
		return "Mod", true
	case opcode.Xor:
		return "Xor", true
	case opcode.Div:
		return "Div", true
	case opcode.Mul:
		return "Mul", true
	case opcode.Not:
		return "Not", true
	// case opcode.Not2:
	// 	return "Not2", true
	case opcode.BitNeg:
		return "BitNeg", true
	case opcode.IntDiv:
		return "IntDiv", true
	case opcode.LogicXor:
		return "LogicXor", true
	case opcode.NullEQ:
		return "NullEQ", true
	case opcode.In:
		return "In", true
	case opcode.Like:
		return "Like", true
	case opcode.Case:
		return "Case", true
	case opcode.Regexp:
		return "Regexp", true
	case opcode.IsNull:
		return "IsNull", true
	case opcode.IsTruth:
		return "IsTruth", true
	case opcode.IsFalsity:
		return "IsFalsity", true
	default:
		return "", false
	}
}

// classifyIndexKeyType 解析Index key的类型
func classifyIndexKeyType(ikt ast.IndexKeyType) (string, bool) {
	switch ikt {
	case ast.IndexKeyTypeNone:
		return "None", true
	case ast.IndexKeyTypeUnique:
		return "Unique", true
	case ast.IndexKeyTypeSpatial:
		return "Spatial", true
	case ast.IndexKeyTypeFullText:
		return "FullText", true
	default:
		return "", false
	}
}

// classifyIndexType 解析Index类型
func classifyIndexType(it model.IndexType) (string, bool) {
	switch it {
	case model.IndexTypeInvalid:
		return "Invalid", true
	case model.IndexTypeBtree:
		return "Btree", true
	case model.IndexTypeHash:
		return "Hash", true
	case model.IndexTypeRtree:
		return "Rtree", true
	default:
		return "", false
	}
}

// classifyIndexHintType 解析Index Hint类型
func classifyIndexHintType(iht ast.IndexHintType) (string, bool) {
	switch iht {
	case ast.HintUse:
		return "Use", true
	case ast.HintIgnore:
		return "Ignore", true
	case ast.HintForce:
		return "Force", true
	default:
		return "", false
	}
}

// classifyIndexHintScope 解析Index Hint Scope
func classifyIndexHintScope(ihs ast.IndexHintScope) (string, bool) {
	switch ihs {
	case ast.HintForScan:
		return "Scan", true
	case ast.HintForJoin:
		return "Join", true
	case ast.HintForOrderBy:
		return "OrderBy", true
	case ast.HintForGroupBy:
		return "GroupBy", true
	default:
		return "", false
	}
}

// classifyIndexVisibility 解析索引可见性
func classifyIndexVisibility(iv ast.IndexVisibility) (string, bool) {
	switch iv {
	case ast.IndexVisibilityDefault:
		return "Default", true
	case ast.IndexVisibilityVisible:
		return "Visible", true
	case ast.IndexVisibilityInvisible:
		return "Invisible", true
	default:
		return "", false
	}
}

// classifyPrivilegeType 解析权限类型
func classifyPrivilegeType(pt mysql.PrivilegeType) (string, bool) {
	switch pt {
	// UsagePriv is a synonym for "no privileges"
	//case mysql.UsagePriv:
	// CreatePriv is the privilege to create schema/table.
	case mysql.CreatePriv:
		return "create", true
	// SelectPriv is the privilege to read from table.
	case mysql.SelectPriv:
		return "select", true
	// InsertPriv is the privilege to insert data into table.
	case mysql.InsertPriv:
		return "insert", true
	// UpdatePriv is the privilege to update data in table.
	case mysql.UpdatePriv:
		return "update", true
	// DeletePriv is the privilege to delete data from table.
	case mysql.DeletePriv:
		return "delete", true
	// ShowDBPriv is the privilege to run show databases statement.
	case mysql.ShowDBPriv:
		return "show database", true
	// SuperPriv enables many operations and server behaviors.
	case mysql.SuperPriv:
		return "super", true
	// CreateUserPriv is the privilege to create user.
	case mysql.CreateUserPriv:
		return "create user", true
	// TriggerPriv is not checked yet.
	case mysql.TriggerPriv:
		return "trigger", true
	// DropPriv is the privilege to drop schema/table.
	case mysql.DropPriv:
		return "drop", true
	// ProcessPriv pertains to display of information about the threads executing within the server.
	case mysql.ProcessPriv:
		return "process", true
	// GrantPriv is the privilege to grant privilege to user.
	case mysql.GrantPriv:
		return "grant", true
	// ReferencesPriv is not checked yet.
	case mysql.ReferencesPriv:
		return "references", true
	// AlterPriv is the privilege to run alter statement.
	case mysql.AlterPriv:
		return "alter", true
	// ExecutePriv is the privilege to run execute statement.
	case mysql.ExecutePriv:
		return "execute", true
	// IndexPriv is the privilege to create/drop index.
	case mysql.IndexPriv:
		return "index", true
	// CreateViewPriv is the privilege to create view.
	case mysql.CreateViewPriv:
		return "create view", true
	// ShowViewPriv is the privilege to show create view.
	case mysql.ShowViewPriv:
		return "show view", true
	// CreateRolePriv the privilege to create a role.
	case mysql.CreateRolePriv:
		return "create role", true
	// DropRolePriv is the privilege to drop a role.
	case mysql.DropRolePriv:
		return "drop role", true
	case mysql.CreateTMPTablePriv:
		return "create temporary tables", true
	case mysql.LockTablesPriv:
		return "lock tables", true
	case mysql.CreateRoutinePriv:
		return "create routine", true
	case mysql.AlterRoutinePriv:
		return "alter routine", true
	case mysql.EventPriv:
		return "event", true
	// ShutdownPriv the privilege to shutdown a server.
	case mysql.ShutdownPriv:
		return "shutdown", true
	// ReloadPriv is the privilege to enable the use of the FLUSH statement.
	case mysql.ReloadPriv:
		return "reload", true
	// FilePriv is the privilege to enable the use of LOAD DATA and SELECT ... INTO OUTFILE.
	case mysql.FilePriv:
		return "file", true
	// ConfigPriv is the privilege to enable the use SET CONFIG statements.
	case mysql.ConfigPriv:
		return "config", true
	// CreateTablespacePriv is the privilege to create tablespace.
	//case mysql.CreateTablespacePriv:
	// ReplicationClientPriv is used in MySQL replication
	//case mysql.ReplicationClientPriv:
	// ReplicationSlavePriv is used in MySQL replication
	//case mysql.ReplicationSlavePriv:
	// AllPriv is the privilege for all actions.
	case mysql.AllPriv:
		return "all", true
	default:
		return "", true
	}
}

// classifyObjectTypeType 解析对象类型
func classifyObjectTypeType(ott ast.ObjectTypeType) (string, bool) {
	switch ott {
	// ObjectTypeNone is for empty object type.
	case ast.ObjectTypeNone:
		return "None", true
	// ObjectTypeTable means the following object is a table.
	case ast.ObjectTypeTable:
		return "Table", true
	default:
		return "", false
	}
}

// classifySelectIntoType 解析Select Into类型
func classifySelectIntoType(sit ast.SelectIntoType) (string, bool) {
	switch sit {
	case ast.SelectIntoOutfile:
		return "Outfile", true
	case ast.SelectIntoDumpfile:
		return "Dumpfile", true
	case ast.SelectIntoVars:
		return "Vars", true
	default:
		return "", false
	}
}

// classifyGrantLevelType 解析授权级别类型
func classifyGrantLevelType(glt ast.GrantLevelType) (string, bool) {
	switch glt {
	// GrantLevelNone is the dummy const for default value.
	case ast.GrantLevelNone:
		return "None", true
	// GrantLevelGlobal means the privileges are administrative or apply to all databases on a given server.
	case ast.GrantLevelGlobal:
		return "Global", true
	// GrantLevelDB means the privileges apply to all objects in a given database.
	case ast.GrantLevelDB:
		return "Database", true
	// GrantLevelTable means the privileges apply to all columns in a given table.
	case ast.GrantLevelTable:
		return "Table", true
	default:
		return "", false
	}
}

// classifyFuncCallExprType 解析Function Call Expr类型
// func classifyFuncCallExprType(fcet ast.FuncCallExprType) (string, bool) {
// }

// classifyCastFunctionType 解析Cast Function类型
func classifyCastFunctionType(cft ast.CastFunctionType) (string, bool) {
	// Cast, Convert or Binary
	switch cft {
	case ast.CastFunction:
		return "Cast", true
	case ast.CastConvertFunction:
		return "Convert", true
	case ast.CastBinaryOperator:
		return "Binary", true
	default:
		return "", false
	}
}

// classifyCompletionType 解析Completion类型
func classifyCompletionType(act ast.CompletionType) (string, bool) {
	// Cast, Convert or Binary
	switch act {
	case ast.CompletionTypeDefault:
		return "Default", true
	case ast.CompletionTypeChain:
		return "Chain", true
	case ast.CompletionTypeRelease:
		return "Release", true
	default:
		return "", false
	}
}

// classifySetRoleStmtType 解析SetRoleStmt类型
func classifySetRoleStmtType(asrst ast.SetRoleStmtType) (string, bool) {
	switch asrst {
	case ast.SetRoleDefault:
		return "Default", true
	case ast.SetRoleNone:
		return "None", true
	case ast.SetRoleAll:
		return "All", true
	case ast.SetRoleAllExcept:
		return "AllExcept", true
	case ast.SetRoleRegular:
		return "Regular", true
	default:
		return "", false
	}
}

// classifyFlushStmtType 解析FlushStmt类型
func classifyFlushStmtType(afst ast.FlushStmtType) (string, bool) {
	switch afst {
	case ast.FlushNone:
		return "None", true
	case ast.FlushTables:
		return "Tables", true
	case ast.FlushPrivileges:
		return "Privileges", true
	case ast.FlushStatus:
		return "Status", true
	case ast.FlushTiDBPlugin:
		return "TiDBPlugin", true
	case ast.FlushHosts:
		return "Hosts", true
	case ast.FlushLogs:
		return "Logs", true
	default:
		return "", false
	}
}

// classifyLogType 解析Log类型
func classifyLogType(alt ast.LogType) (string, bool) {
	switch alt {
	case ast.LogTypeDefault:
		return "Default", true
	case ast.LogTypeBinary:
		return "Binary", true
	case ast.LogTypeEngine:
		return "Engine", true
	case ast.LogTypeError:
		return "Error", true
	case ast.LogTypeGeneral:
		return "General", true
	case ast.LogTypeSlow:
		return "Slow", true
	default:
		return "", false
	}
}

// classifyShowSlowType 解析Show Slow类型
func classifyShowSlowType(asst ast.ShowSlowType) (string, bool) {
	switch asst {
	case ast.ShowSlowTop:
		return "Top", true
	case ast.ShowSlowRecent:
		return "Recent", true
	default:
		return "", false
	}
}

// classifyShowSlowKind 解析Show Slow类型
func classifyShowSlowKind(assk ast.ShowSlowKind) (string, bool) {
	switch assk {
	// ShowSlowKindDefault is a ShowSlowKind constant.
	case ast.ShowSlowKindDefault:
		return "Default", true
	// ShowSlowKindInternal is a ShowSlowKind constant.
	case ast.ShowSlowKindInternal:
		return "Internal", true
	// ShowSlowKindAll is a ShowSlowKind constant.
	case ast.ShowSlowKindAll:
		return "All", true
	default:
		return "", false
	}
}

// classifyTimestampBoundMode 解析Timestamp Bound模式
func classifyTimestampBoundMode(atsbm ast.TimestampBoundMode) (string, bool) {
	switch atsbm {
	case ast.TimestampBoundStrong:
		return "Strong", true
	case ast.TimestampBoundMaxStaleness:
		return "MaxStaleness", true
	case ast.TimestampBoundExactStaleness:
		return "ExactStaleness", true
	case ast.TimestampBoundReadTimestamp:
		return "ReadTimestamp", true
	case ast.TimestampBoundMinReadTimestamp:
		return "MinReadTimestamp", true
	default:
		return "", false
	}
}

// classifyViewAlgorithm 解析View Algorithm
func classifyViewAlgorithm(mva model.ViewAlgorithm) (string, bool) {
	switch mva {
	case model.AlgorithmUndefined:
		return "Undefined", true
	case model.AlgorithmMerge:
		return "Merge", true
	case model.AlgorithmTemptable:
		return "Temptable", true
	default:
		return "", false
	}
}

// classifyViewSecurity 解析View Security
func classifyViewSecurity(mvs model.ViewSecurity) (string, bool) {
	switch mvs {
	case model.SecurityDefiner:
		return "Definer", true
	case model.SecurityInvoker:
		return "Invoker", true
	default:
		return "", false
	}
}

// classifyViewCheckOption 解析View CheckOption
func classifyViewCheckOption(mvco model.ViewCheckOption) (string, bool) {
	switch mvco {
	case model.CheckOptionLocal:
		return "Local", true
	case model.CheckOptionCascaded:
		return "Cascaded", true
	default:
		return "", false
	}
}

// beautifyModelCIStr 美化model.CIStr的可读性
func beautifyModelCIStr(mcs *model.CIStr) map[string]string {
	if mcs == nil {
		return nil
	}

	oname := mcs.O
	lname := mcs.L
	res := map[string]string{
		"original":   oname,
		"lower-case": lname,
	}

	return res
}
