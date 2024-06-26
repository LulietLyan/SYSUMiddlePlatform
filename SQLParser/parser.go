package SQLParser

import (
	"errors"
	"fmt"
	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
)

type TiStmt struct {
	Stmts []ast.StmtNode
}

// 解析一条SQL语句
func NewParseOneStmt(sqltext, charset, collation string) (ast.StmtNode, error) {
	// tidb parser 语法解析
	stmt, err := parser.New().ParseOneStmt(sqltext, charset, collation)
	if err != nil {
		return stmt, fmt.Errorf("SQL解析错误:%s", err.Error())
	}
	return stmt, nil
}

// 获取表名
func GetTableNameFromAlterStatement(sqltext string) (string, error) {
	stmt, err := NewParseOneStmt(sqltext, "", "")
	if err != nil {
		return "", err
	}
	switch s := stmt.(type) {
	case *ast.AlterTableStmt:
		return s.Table.Name.String(), nil
	}
	return "", errors.New("未提取到表名")
}
