package types

import "github.com/pingcap/parser/ast"

// Constraint 约束
type Constraint struct {
	// only supported by MariaDB 10.0.2+ (ADD {INDEX|KEY}, ADD FOREIGN KEY),
	// see https://mariadb.com/kb/en/library/alter-table/
	IfNotExists  bool
	Type         string // ConstraintType
	Name         string
	Keys         []*IndexPartSpecification // Used for PRIMARY KEY, UNIQUE, ......
	Refer        *ReferenceDef             // Used for foreign key.
	Option       *IndexOption              // Index Options
	Expr         interface{}               // ExprNode // Used for Check
	Enforced     bool                      // Used for Check
	InColumn     bool                      // Used for Check
	InColumnName string                    // Used for Check
	IsEmptyIndex bool                      // Used for Check
}

// ParseConstraint 解析约束
func (cst *Constraint) ParseConstraint(acst *ast.Constraint) *Constraint {
	if acst == nil {
		return nil
	}
	// IfNotExists  bool
	cst.IfNotExists = acst.IfNotExists

	// Type         ConstraintType
	if val, bv := classifyConstraintType(acst); bv == true {
		cst.Type = val
	}

	// Name         string
	cst.Name = acst.Name

	// Keys         []*IndexPartSpecification
	for _, aips := range acst.Keys {
		ips := new(IndexPartSpecification)
		ipsitem := ips.ParseIndexPartSpecification(aips)

		cst.Keys = append(cst.Keys, ipsitem)
	}

	// Refer        *ReferenceDef
	rd := new(ReferenceDef)
	cst.Refer = rd.ParseReferenceDef(acst.Refer)

	// Option       *IndexOption
	io := new(IndexOption)
	cst.Option = io.ParseIndexOption(acst.Option)

	// Expr         ExprNode
	cst.Expr = acst.Expr

	// Enforced     bool
	cst.Enforced = acst.Enforced

	// InColumn     bool
	//cst.InColumn = acst.InColumn

	// InColumnName string
	//cst.InColumnName = acst.InColumnName

	// IsEmptyIndex bool
	//cst.IsEmptyIndex = acst.IsEmptyIndex

	return cst
}
