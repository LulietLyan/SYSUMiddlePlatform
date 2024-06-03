package models

// 通过解析语法树获得的输出
type SQLTreeJSON struct {
	StmtType string `json:"StmtType"`
	StmtTree struct {
		IfNotExists bool `json:"IfNotExists"`
		IsTemporary bool `json:"IsTemporary"`
		Table       struct {
			Schema struct {
				LowerCase string `json:"lower-case"`
				Original  string `json:"original"`
			} `json:"Schema"`
			Name struct {
				LowerCase string `json:"lower-case"`
				Original  string `json:"original"`
			} `json:"Name"`
			IndexHints     interface{}   `json:"IndexHints"`
			PartitionNames []interface{} `json:"PartitionNames"`
		} `json:"Table"`
		ReferTable interface{} `json:"ReferTable"`
		Columns    []struct {
			Name struct {
				Schema struct {
					LowerCase string `json:"lower-case"`
					Original  string `json:"original"`
				} `json:"Schema"`
				Table struct {
					LowerCase string `json:"lower-case"`
					Original  string `json:"original"`
				} `json:"Table"`
				Name struct {
					LowerCase string `json:"lower-case"`
					Original  string `json:"original"`
				} `json:"Name"`
			} `json:"Name"`
			Type struct {
				Datatype string      `json:"Datatype"`
				Flag     int         `json:"Flag"`
				Fieldlen int         `json:"Fieldlen"`
				Decimal  int         `json:"Decimal"`
				Charset  string      `json:"Charset"`
				Collate  string      `json:"Collate"`
				Elems    interface{} `json:"Elems"`
			} `json:"Type"`
			Options []struct {
				Type                string      `json:"Type"`
				Expr                *string     `json:"Expr"`
				Stored              bool        `json:"Stored"`
				Refer               interface{} `json:"Refer"`
				StrValue            string      `json:"StrValue"`
				AutoRandomBitLength int         `json:"AutoRandomBitLength"`
				Enforced            bool        `json:"Enforced"`
				ConstraintName      string      `json:"ConstraintName"`
			} `json:"Options"`
		} `json:"Columns"`
		Constraints interface{} `json:"Constraints"`
		Options     []struct {
			Type       string      `json:"Type"`
			Default    bool        `json:"Default"`
			Value      string      `json:"Value"`
			TableNames interface{} `json:"TableNames"`
		} `json:"Options"`
		Partition   interface{} `json:"Partition"`
		OnDuplicate string      `json:"OnDuplicate"`
		Select      interface{} `json:"Select"`
	} `json:"StmtTree"`
}
