package types

import "github.com/pingcap/parser/ast"

// DatabaseOption like *ast.DatabaseOption
type DatabaseOption struct {
	Type  string // DatabaseOptionType
	Value string
}

// ParseDatabaseOption 解析Database Option
func (do *DatabaseOption) ParseDatabaseOption(ado *ast.DatabaseOption) *DatabaseOption {
	if ado == nil {
		return nil
	}

	// Type  // DatabaseOptionType
	if val, bv := classifyDatabaseOptionType(ado); bv == true {
		do.Type = val
	}

	// Value string
	do.Value = ado.Value

	return do
}
