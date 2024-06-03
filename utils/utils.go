package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"reflect"
	"time"

	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	//_ "github.com/pingcap/parser/test_driver"
)

// JSONMarshal 不转义HTML
func JSONMarshal(t interface{}) ([]byte, error) {
	buf := &bytes.Buffer{}
	encoder := json.NewEncoder(buf)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)

	return buf.Bytes(), err
}

// JSONMarshalIndent 格式化JSONMarshal的结果
func JSONMarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	b, err := JSONMarshal(v)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	err = json.Indent(&buf, b, prefix, indent)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// IsBlank 判断是否为零值
func IsBlank(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String:
		return value.Len() == 0
	case reflect.Bool:
		return !value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return value.IsNil()
	}
	return reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface())
}

// RemoveDupOnSlice 对[]string进行去重
func RemoveDupOnSlice(slc []string) []string {
	res := []string{}
	tempMap := map[string]byte{}
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = 0
		if len(tempMap) != l {
			res = append(res, e)
		}
	}
	return res
}

// Crc32Uint32 生成CRC32
func Crc32Uint32(s string) uint32 {
	return crc32.ChecksumIEEE([]byte(s))
}

// Md5String 生成MD5字符串
func Md5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))

	return hex.EncodeToString(h.Sum(nil))
}

// TimeNow 获取当前时间
func TimeNow(time.Time) (string, error) {
	timeStr := time.Now().Format("2006-01-02 15:04:05")

	loc, _ := time.LoadLocation("Local") //获取时区

	formatTime, err := time.ParseInLocation("2006-01-02 15:04:05", timeStr, loc)
	if err != nil {
		return "", err
	}
	return formatTime.Format("2006-01-02 15:04:05"), err
}

// PrintResult 打印结果，根据pretty标志位确定是否美化json
func PrintResult(r interface{}, pretty bool) string {
	var (
		res     []byte
		jsonerr error
	)

	switch {
	case pretty == true:
		res, jsonerr = JSONMarshalIndent(r, "", "    ")
	case pretty == false:
		res, jsonerr = JSONMarshal(r)
	}
	if jsonerr != nil {
		fmt.Println(jsonerr.Error())
	}

	return string(res)
}

// ParseSQL 解析SQL为[]ast.StmtNode
func ParseSQL(sql string, charset string, collation string) ([]ast.StmtNode, error) {
	p := parser.New()

	// SQL MODE设置
	// p.SetSQLMode(mode mysql.SQLMode)
	// p.SetStrictDoubleTypeCheck(val bool)

	stmtNodes, _, err := p.Parse(sql, charset, collation)

	return stmtNodes, err
}
