package dto

type LoginDto struct {
	Action    string `form:"action"`
	AppID     string `form:"appID"`
	AppSecret string `form:"appSecret"`
	Lang      string `form:"lang"`
	Uid       string `form:"uid"`
	Name      string `form:"name"`
	Debug     string `form:"debug"`
	IsHttps   string `form:"is_https"`
	CloseType string `form:"closeType"`
	PhoneType string `form:"PhoneType"`
	GameCode  string `form:"gameCode"`
	Optional  string `form:"optional"`
}
