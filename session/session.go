package session

import (
	"gosqltree/types"
	"gosqltree/utils"

	"github.com/pingcap/parser/ast"
)

// Session 存放结果
type Session struct {
	SQLID             string `json:"SQL_ID"`
	SQLText           string `json:"SQL_Text"`
	*types.SQLElement `json:"SQL_Element"`
	*types.SQLTree    `json:"SQL_Tree"`
}

// GetResult 解析ast.StmtNode
func (s *Session) GetResult(stmtNode ast.StmtNode) *Session {
	// SQLID & SQLText
	if stmtNode.Text() != "" {
		// s.SQLID = CRC32Uint32(stmtNode.Text())
		s.SQLID = utils.Md5String(stmtNode.Text())
		s.SQLText = stmtNode.Text()
	} else {
		// s.SQLID = 000000000
		s.SQLID = ""
		s.SQLText = ""
	}

	// SQLElement
	se := new(types.SQLElement)
	stmtNode.Accept(se)
	s.SQLElement = se

	// SQLTree
	st := new(types.SQLTree)
	st.ClassifyStmtNode(stmtNode)
	s.SQLTree = st

	return s
}
