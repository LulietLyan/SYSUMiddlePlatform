package SQLParser

import (
	"backend/models"
	"backend/session"
	"backend/utils"
	"encoding/json"
	_ "github.com/pingcap/parser/test_driver"
)

func SQLTreeGenerator(sqlCommand string) models.SQLTreeJSON {
	var stj models.SQLTreeJSON
	s := new(session.Session) // 等价于var s *Session = new(Session)

	// 解析SQL
	stmtNodes, err := utils.ParseSQL(sqlCommand, "", "")
	if err != nil {
		return stj
	}

	for _, stmtNode := range stmtNodes {
		s.GetResult(stmtNode)
	}

	json.Unmarshal([]byte(utils.PrintResult(s.SQLTree, true)), &stj)

	return stj
}
