package models

import (
	"github.com/jinzhu/gorm"
)

// 用户基本属性
type User struct {
	gorm.Model
	Id                    string  `json:"id,string" db:"id"`
	Password              string  `json:"password" db:"password"`
	UserName              string  `json:"username" db:"user_name"`
	Url                   string  `json:"url" db:"url"`
	HeadUrl               string  `json:"avatar" db:"head_url"`
	BackgroundURL         string  `json:"backgroundURL" db:"background_url"`
	BackgroundColor       string  `json:"backgroundColor" db:"background_color"`
	ButtonColor           string  `json:"buttonColor" db:"button_color"`
	CustomBackgroundColor bool    `json:"customBackgroundColor" db:"custom_background_color"`
	CustomButtonColor     bool    `json:"customButtonColor" db:"custom_button_color"`
	PresetColor           int     `json:"presetColor" db:"preset_color"`
	BackgroundAlpha       float64 `json:"backgroundAlpha" db:"background_alpha"`
	ButtonAlpha           float64 `json:"buttonAlpha" db:"button_alpha"`
	SearchItemCount       int     `json:"searchItemCount" db:"search_item_count""`
	FontColor             string  `json:"fontColor" db:"font_color"`
	Title                 string  `json:"title" db:"title"`
	Icon                  string  `json:"icon?" db:"icon"`
	TextIcon              string  `json:"textIcon" db:"text_icon"`
	Custom                string  `json:"custom" db:"custom"`
	Page                  string  `json:"page" db:"page"`
	Record                string  `json:"record" db:"record"`
	SearchEngine          string  `json:"searchEngine" db:"search_engine"`
	SuggestApi            string  `json:"suggestAPI" db:"suggest_api"`
	OpenNewPage           string  `json:"openNewPage" db:"open_new_page"`
	ShowEngineIcon        string  `json:"showEngineIcon" db:"show_engine_icon"`
	ShowEngineList        string  `json:"showEngineList" db:"show_engine_list"`
	Language              string  `json:"language" db:"language"`
	Log                   string  `json:"log" db:"log"`
	Date                  string  `json:"date" db:"date"`
}

type PresetBackground struct {
	Bgid          int    `json:"id" db:"bgid"`
	BackgroundURL string `json:"backgroundURL" db:"background_url"`
	Description   string `json:"description" db:"description"`
}

// RegisterForm 注册请求参数
type RegisterForm struct {
	UserName        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

// LoginForm 登录请求参数
type LoginForm struct {
	Id       string `json:"id"`
	Password string `json:"password"`
}

// 设置请求参数
type SettingForm struct {
	Id                    string  `json:"Id" db:"id"`
	BackgroundAlpha       float64 `json:"backgroundAlpha" db:"background_alpha"`
	BackgroundColor       string  `json:"backgroundColor" db:"background_color"`
	BackgroundURL         string  `json:"backgroundURL" db:"background_url"`
	ButtonAlpha           float64 `json:"buttonAlpha" db:"button_alpha"`
	ButtonColor           string  `json:"buttonColor" db:"button_color"`
	CustomBackgroundColor bool    `json:"customBackgroundColor" db:"custom_background_color"`
	CustomButtonColor     bool    `json:"customButtonColor" db:"custom_button_color"`
	FontColor             string  `json:"fontColor" db:"font_color"`
	Language              string  `json:"language" db:"language"`
	OpenNewPage           string  `json:"openNewPage" db:"open_new_page"`
	PresetColor           int     `json:"presetColor" db:"preset_color"`
	SearchEngine          string  `json:"searchEngine" db:"search_engine"`
	ShowEngineIcon        string  `json:"showEngineIcon" db:"show_engine_icon"`
	ShowEngineList        string  `json:"showEngineList" db:"show_engine_list"`
	SuggestApi            string  `json:"suggestAPI" db:"suggest_api"`
	SearchItemCount       int     `json:"searchItemCount" db:"search_item_count""`
}

// 销号请求参数
type DeleteForm struct {
	Id string `json:"Id" db:"id"`
}

// 上传图片请求参数
type UploadImageForm struct {
	Id          string `json:"id" db:"id"`
	ImageData   string `json:"imageData"`
	Description string `json:"description" db:"description"`
}

// 修改密码请求参数
type ChangePasswdForm struct {
	Id     string `json:"Id" db:"id"`
	NewPwd string `json:"newPassword" db:"password"`
}

// 修改昵称请求参数
type ChangeUserNameForm struct {
	Id      string `json:"Id" db:"id"`
	NewName string `json:"userName" db:"user_name"`
}

// 上传URL请求参数
type NewURLForm struct {
	Title    string `json:"title" db:"title"`
	Url      string `json:"url" db:"url"`
	Icon     string `json:"icon" db:"icon"`
	TextIcon string `json:"textIcon" db:"text_icon"`
	Custom   string `json:"custom" db:"custom"`
	Page     string `json:"page" db:"page"`
}

// 上传搜索记录请求参数
type NewRecordForm struct {
	Record string `json:"record" db:"record"`
}

// 百度搜索建议参数
type BaiduSuggestionResponse struct {
	Query string   `json:"q"`
	P     bool     `json:"p"`
	Data  []string `json:"s"`
}

// 必应搜索建议参数
type BingSuggestionResponse struct {
	AS struct {
		Query       string `json:"Query"`
		FullResults int    `json:"FullResults"`
		Results     []struct {
			Type     string `json:"Type"`
			Suggests []struct {
				Txt  string `json:"Txt"`
				Type string `json:"Type"`
				Sk   string `json:"Sk"`
				HCS  int    `json:"HCS"`
			} `json:"Suggests"`
		} `json:"Results"`
	} `json:"AS"`
}
type SortedBingSuggestionResponse struct {
	Txt []string `json:"Txt"`
}

// 谷歌搜索建议参数
type Suggestion struct {
	Data string `xml:"suggestion>data"`
}

type GoogleSuggestionResponse struct {
	Suggestions []Suggestion `xml:"CompleteSuggestion"`
}

// 整合信息
type DataGroup struct {
	Rank int    `json:"rank"`
	Name string `json:"name"`
	Url  string `json:"url"`
}

// IT热榜
type GetITAPI struct {
	Success    bool   `json:"success"`
	Title      string `json:"title"`
	Subtitle   string `json:"subtitle"`
	UpdateTime string `json:"update_time"`
	Data       []struct {
		Index    int    `json:"index"`
		Title    string `json:"title"`
		Url      string `json:"url"`
		MobilUrl string `json:"mobilUrl"`
	} `json:"data"`
}

// 微博热榜
type GetWbAPI struct {
	Success    bool   `json:"success"`
	Title      string `json:"title"`
	Subtitle   string `json:"subtitle"`
	UpdateTime string `json:"update_time"`
	Data       []struct {
		Index    int    `json:"index"`
		Title    string `json:"title"`
		Hot      string `json:"hot"`
		Url      string `json:"url"`
		MobilUrl string `json:"mobilUrl"`
	} `json:"data"`
}

// 知乎热榜
type GetZHAPI struct {
	Success    bool   `json:"success"`
	Title      string `json:"title"`
	Subtitle   string `json:"subtitle"`
	UpdateTime string `json:"update_time"`
	Data       []struct {
		Index    int    `json:"index"`
		Title    string `json:"title"`
		Desc     string `json:"desc"`
		Pic      string `json:"pic"`
		Hot      string `json:"hot"`
		Url      string `json:"url"`
		MobilUrl string `json:"mobilUrl"`
	} `json:"data"`
}

// 哔哩哔哩热榜
type GetBiliAPI struct {
	Success    bool   `json:"success"`
	Title      string `json:"title"`
	Subtitle   string `json:"subtitle"`
	UpdateTime string `json:"update_time"`
	Data       []struct {
		Index    int    `json:"index"`
		Title    string `json:"title"`
		Desc     string `json:"desc"`
		Pic      string `json:"pic"`
		Hot      string `json:"hot"`
		Url      string `json:"url"`
		MobilUrl string `json:"mobilUrl"`
	} `json:"data"`
}

type LogForm struct {
	Id   string `json:"id" db:"id"`
	Log  string `json:"log" db:"log"`
	Date string `json:"date" db:"date"`
}

type GetLogForm struct {
	Log  []string `json:"log" db:"log"`
	Date []string `json:"date" db:"date"`
}

// 数据表名字
func (User) TableName() string {
	return "user"
}

func (PresetBackground) TableName() string {
	return "presetBackground"
}
