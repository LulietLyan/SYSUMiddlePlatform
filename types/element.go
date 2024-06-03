package types

import (
	"gosqltree/utils"

	"github.com/pingcap/parser/ast"
)

// ast.Node接口的Accept方法，接受Visitor接口（Enter和Leave方法），对AST的处理可以依赖该Accept方法，遍历所有的ast节点并按照不同AST节点进行解析
// type Visitor interface {
// 	   Enter(n Node) (node Node, skipChildren bool)
// 	   Leave(n Node) (node Node, ok bool)
// }

// SQLElement 存储SQL语句中的对象信息
type SQLElement struct {
	SchemaList []string // 数据库列表
	TableList  []string // 表列表
	ColumnList []string // 字段列表
	AsNameList []string // 别名列表
}

// Enter Visitor接口的Enter方法
func (se *SQLElement) Enter(in ast.Node) (out ast.Node, skipChildren bool) {
	// fmt.Printf("%T\n", in)
	// fmt.Println(in)
	se.ParseSQLElement(in)

	return in, false
}

// Leave Visitor接口的Leave方法
func (se *SQLElement) Leave(in ast.Node) (out ast.Node, ok bool) {
	return in, true
}

// ParseSQLElement 解析SQL语句中的元素（表名、字段名等。。）
func (se *SQLElement) ParseSQLElement(node ast.Node) {
	switch n := node.(type) {
	// 表名
	case *ast.TableName:
		sn := n.Schema.O
		if sn != "" {
			se.SchemaList = append(se.SchemaList, sn)
		}

		tn := n.Name.O
		if tn != "" {
			se.TableList = append(se.TableList, tn)
		}
	// 字段名
	case *ast.ColumnName:
		sn := n.Schema.O
		if sn != "" {
			se.SchemaList = append(se.SchemaList, sn)
		}

		// 统计ast.ColumnName.Table会出现别名被记录到TableList的情况
		// tn := n.Table.O
		// if tn != "" {
		// 	se.TableList = append(se.TableList, tn)
		// }

		cn := n.Name.O
		if cn != "" {
			se.ColumnList = append(se.ColumnList, cn)
		}
	// SelectField.AsName | TableSource.AsName 作为别名记录到AsNameList
	case *ast.SelectField:
		an := n.AsName.O
		if an != "" {
			se.AsNameList = append(se.AsNameList, an)
		}
	case *ast.TableSource:
		an := n.AsName.O
		if an != "" {
			se.AsNameList = append(se.AsNameList, an)
		}
	default:
	}

	// 去重
	se.SchemaList = utils.RemoveDupOnSlice(se.SchemaList)
	se.TableList = utils.RemoveDupOnSlice(se.TableList)
	se.ColumnList = utils.RemoveDupOnSlice(se.ColumnList)
	se.AsNameList = utils.RemoveDupOnSlice(se.AsNameList)
}
