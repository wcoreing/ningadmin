package dto

type RegisterDto struct {
	Action    string `form:"action"`
	AppID     string `form:"appID"`
	AppSecret string `form:"appSecret"`
	Lang      string `form:"lang"`
	Uid       string `form:"uid"`
	Name      string `form:"name"`
	Debug     string `form:"debug"`
}
