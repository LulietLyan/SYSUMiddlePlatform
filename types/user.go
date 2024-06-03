package types

import (
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/auth"
)

// FrameType like *ast.FrameType
type FrameType int

// FrameExtent like *ast.FrameExtent
type FrameExtent struct {
	Start *FrameBound
	End   *FrameBound
}

// ParseFrameExtent 解析FrameExtent
func (fe *FrameExtent) ParseFrameExtent(afe *ast.FrameExtent) *FrameExtent {
	if afe == nil {
		return nil
	}

	// Start FrameBound
	fb1 := new(FrameBound)
	fe.Start = fb1.ParseFrameBound(&afe.Start)

	// End   FrameBound
	fb2 := new(FrameBound)
	fe.End = fb2.ParseFrameBound(&afe.End)

	return fe
}

// BoundType like *ast.BoundType
type BoundType int

// FrameBound like *ast.FrameBound
type FrameBound struct {
	Type      BoundType
	UnBounded bool
	Expr      interface{} // ExprNode
	// `Unit` is used to indicate the units in which the `Expr` should be interpreted.
	// For example: '2:30' MINUTE_SECOND.
	Unit string // TimeUnitType
}

// ParseFrameBound 解析ParseFrameBound
func (fb *FrameBound) ParseFrameBound(afb *ast.FrameBound) *FrameBound {
	if afb == nil {
		return nil
	}

	// Type      BoundType
	fb.Type = BoundType(afb.Type)

	// UnBounded bool
	fb.UnBounded = afb.UnBounded

	// Expr      ExprNode
	if val, bv := classifyExprNode(afb.Expr); bv == true {
		fb.Expr = val
	}

	// Unit      TimeUnitType
	if val, bv := classifyTimeUnitType(afb.Unit); bv == true {
		fb.Unit = val
	}

	return fb
}

// UserIdentity like *auth.UserIdentity
type UserIdentity struct {
	Username     string
	Hostname     string
	CurrentUser  bool
	AuthUsername string // Username matched in privileges system
	AuthHostname string // Match in privs system (i.e. could be a wildcard)
}

// ParseUserIdentity 解析User Identity
func (ui *UserIdentity) ParseUserIdentity(aui *auth.UserIdentity) *UserIdentity {
	if aui == nil {
		return nil
	}

	ui.Username = aui.Username
	ui.Hostname = aui.Hostname
	ui.CurrentUser = aui.CurrentUser
	ui.AuthUsername = aui.AuthUsername
	ui.AuthHostname = aui.AuthHostname

	return ui
}

// RoleIdentity like *auth.RoleIdentity
type RoleIdentity struct {
	Username string
	Hostname string
}

// ParseRoleIdentity 解析Role Identity
func (ri *RoleIdentity) ParseRoleIdentity(ari *auth.RoleIdentity) *RoleIdentity {
	if ari == nil {
		return nil
	}

	// Username string
	ri.Username = ari.Username

	// Hostname string
	ri.Hostname = ari.Hostname

	return ri
}

// PrivElem like *ast.PrivElem
type PrivElem struct {
	Priv string // mysql.PrivilegeType
	Cols []*ColumnName
}

// ParsePrivElem 解析Privilege Element
func (pe *PrivElem) ParsePrivElem(ape *ast.PrivElem) *PrivElem {
	if ape == nil {
		return nil
	}

	// Priv mysql.PrivilegeType
	if val, bv := classifyPrivilegeType(ape.Priv); bv == true {
		pe.Priv = val
	}

	// Cols []*ColumnName
	for _, acn := range ape.Cols {
		cn := new(ColumnName)
		cnitem := cn.ParseColumnName(acn)

		pe.Cols = append(pe.Cols, cnitem)
	}

	return pe
}

// GrantLevel like *ast.GrantLevel
type GrantLevel struct {
	Level     string // GrantLevelType
	DBName    string
	TableName string
}

// ParseGrantLevel 解析授权级别
func (gl *GrantLevel) ParseGrantLevel(agl *ast.GrantLevel) *GrantLevel {
	if agl == nil {
		return nil
	}

	// Level     string // GrantLevelType
	if val, bv := classifyGrantLevelType(agl.Level); bv == true {
		gl.Level = val
	}

	// DBName    string
	gl.DBName = agl.DBName

	// TableName string
	gl.TableName = agl.TableName

	return gl
}
