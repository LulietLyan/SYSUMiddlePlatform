package control

import (
	"backend/models"
	"backend/mysql"
	"backend/response"
	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"strconv"
)

// 注册
func RegisterHandler(c *gin.Context) {
	var userRegister *models.RegisterForm
	node, err := snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}
	if err := c.ShouldBind(&userRegister); err != nil {
		response.Fail(c, nil, "注册失败!")
	} else {
		//判断用户名是否存在
		var findUser models.User
		result := mysql.DB.Where("user_name=?", userRegister.UserName).First(&findUser)
		if result.Error == nil {
			response.Fail(c, nil, "用户名已被使用!")
		} else if result.Error != nil {
			var newUser = models.User{
				Id:                    strconv.FormatUint(uint64(node.Generate().Int64()%10000000000), 10),
				Password:              userRegister.Password,
				UserName:              userRegister.UserName,
				Url:                   "https://baidu.com/",
				HeadUrl:               "http://localhost:2020/user/1",
				BackgroundURL:         "http://localhost:2020/user/1",
				ButtonColor:           "#ffffff",
				BackgroundColor:       "#ffffff",
				CustomButtonColor:     false,
				CustomBackgroundColor: false,
				FontColor:             "black",
				PresetColor:           0,
				BackgroundAlpha:       0.8,
				ButtonAlpha:           0.8,
				SearchEngine:          "Baidu",
				SuggestApi:            "Baidu",
				OpenNewPage:           "true",
				ShowEngineIcon:        "true",
				ShowEngineList:        "true",
				Language:              "Auto",
				SearchItemCount:       0,
				Log:                   "",
			}
			mysql.DB.Create(&newUser)
			response.Success(c, gin.H{"id": newUser.Id}, "注册成功!")
		}
	}
}
