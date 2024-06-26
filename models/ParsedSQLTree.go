package models

// 通过解析语法树获得的输出
type ParsedResult struct {
	RequestId string `json:"request_id"`
	Code      string `json:"code"`
	Data      []struct {
		Tables []string `json:"tables"`
		Type   string   `json:"type"`
		Query  string   `json:"query"`
	} `json:"data"`
	Message string `json:"message"`
}

type SQLTreeJSON struct {
	IsReplace bool `json:"IsReplace"`
	IgnoreErr bool `json:"IgnoreErr"`
	Table     struct {
		TableRefs struct {
			Left struct {
				Source struct {
					Schema struct {
						O string `json:"O"`
						L string `json:"L"`
					} `json:"Schema"`
					Name struct {
						O string `json:"O"`
						L string `json:"L"`
					} `json:"Name"`
					DBInfo         interface{} `json:"DBInfo"`
					TableInfo      interface{} `json:"TableInfo"`
					IndexHints     interface{} `json:"IndexHints"`
					PartitionNames interface{} `json:"PartitionNames"`
				} `json:"Source"`
				AsName struct {
					O string `json:"O"`
					L string `json:"L"`
				} `json:"AsName"`
			} `json:"Left"`
			Right        interface{} `json:"Right"`
			Tp           int         `json:"Tp"`
			On           interface{} `json:"On"`
			Using        interface{} `json:"Using"`
			NaturalJoin  bool        `json:"NaturalJoin"`
			StraightJoin bool        `json:"StraightJoin"`
		} `json:"TableRefs"`
	} `json:"Table"`
	Columns []struct {
		Schema struct {
			O string `json:"O"`
			L string `json:"L"`
		} `json:"Schema"`
		Table struct {
			O string `json:"O"`
			L string `json:"L"`
		} `json:"Table"`
		Name struct {
			O string `json:"O"`
			L string `json:"L"`
		} `json:"Name"`
	} `json:"Columns"`
	Lists [][]struct {
		Type struct {
			Tp      int         `json:"Tp"`
			Flag    int         `json:"Flag"`
			Flen    int         `json:"Flen"`
			Decimal int         `json:"Decimal"`
			Charset string      `json:"Charset"`
			Collate string      `json:"Collate"`
			Elems   interface{} `json:"Elems"`
		} `json:"Type"`
	} `json:"Lists"`
	Setlist        interface{}   `json:"Setlist"`
	Priority       int           `json:"Priority"`
	OnDuplicate    interface{}   `json:"OnDuplicate"`
	Select         interface{}   `json:"Select"`
	PartitionNames []interface{} `json:"PartitionNames"`
	TableHints     interface{}   `json:"TableHints"`
}
