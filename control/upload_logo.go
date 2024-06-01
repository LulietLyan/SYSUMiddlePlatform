package control

import (
	"backend/models"
	"backend/mysql"
	"backend/response"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func UploadLogo(c *gin.Context) {
	//从上下文获取用户信息
	var userId uint
	if data, ok := c.Get("userId"); !ok {
		response.Fail(c, nil, "没有从token解析出所需信息")
		return
	} else {
		userId = data.(uint)
	}
	c.Request.ParseMultipartForm(32 << 20)
	//获取所有上传文件信息
	files := c.Request.MultipartForm.File
	uploadDir := "./image/logo_" + fmt.Sprintf("%d", userId)
	for _, fileHeaders := range files {
		for _, fileHeader := range fileHeaders {
			// 创建文件名，可以使用原始文件名或自定义文件名
			originalFilename := fileHeader.Filename
			ext := filepath.Ext(originalFilename)
			newFilename := filepath.Join(uploadDir, ext) // 假设randomString是生成随机字符串的函数

			// 打开文件
			file, err := fileHeader.Open()
			if err != nil {
				response.Fail(c, nil, err.Error())
				return
			}
			defer file.Close()

			// 创建目标文件
			out, err := os.Create(newFilename)
			if err != nil {
				response.Fail(c, nil, err.Error())
				return
			}
			defer out.Close()

			// 将文件内容写入到目标文件
			_, err = io.Copy(out, file)
			if err != nil {
				response.Fail(c, nil, err.Error())
				return
			}

			//保存url到数据库
			var puRecord models.ProjectUser
			if e := mysql.DB.Raw("SELECT ProjectUser.* FROM User,ProjectUser WHERE User.U_uid=ProjectUser.U_uid and User.U_uid = ?", userId).First(&puRecord).Error; e != nil {
				response.Fail(c, nil, "查找项目用户记录时出错")
				return
			}
			puRecord.PU_logo_url = "/logo/logo_" + fmt.Sprintf("%d", userId) + ext
			if e := mysql.DB.Save(&puRecord).Error; e != nil {
				response.Fail(c, nil, "更新项目logo的url路径时出错")
				return
			}
			// 告诉客户端文件上传成功
			response.Success(c, nil, "")
		}
	}
}
