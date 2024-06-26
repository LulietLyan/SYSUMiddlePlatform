package main

import (
	"backend/mysql"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"unsafe"
)

var (
	tokenOfUser_1 = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJMb2dpblRpbWUiOiIyMDI0LTA2LTI1VDIxOjQwOjM1LjMyMzM3NDMrMDg6MDAiLCJVc2VySWQiOjYsIklkZW50aXR5IjoiRGV2ZWxvcGVyIiwiUFVfdWlkIjozfQ.iGCAMDQil6OkM8Z1dZr-6PBgyGDa800WbezQ7ZHF90U"
	writeURL      = "https://127.0.0.1:8087/api/rNw/request/write"
)

func testWriting(projectUser string, tableName string, sqlCommand string) {
	mysql.DB.Select()

	client := &http.Client{}

	data := make(map[string]interface{})
	data["projectName"] = "1"
	data["tableName"] = "Student"
	data["sqlCommand"] = ""

	respdata, _ := json.Marshal(data)

	request, err := http.NewRequest("POST", writeURL, bytes.NewReader(respdata))
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
	}

	request.Header.Set("Authorization", tokenOfUser_1)

	response, err := client.Do(request)
	defer response.Body.Close()
	content, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
	}
	//fmt.Println(string(content))       // 直接打印
	str := (*string)(unsafe.Pointer(&content)) //转化为string,优化内存
	fmt.Println(*str)
}
