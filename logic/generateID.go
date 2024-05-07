package logic

import (
	"github.com/bwmarrin/snowflake"
	"strconv"
)

var node *snowflake.Node

func init() {
	var err error
	node, err = snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}
}

// 生成唯一ID
func GenerateID() string {
	return strconv.FormatInt(node.Generate().Int64(), 10)
}
