package bo

type CheckAgentBo struct {
	Action    string `form:"action"`
	AppID     string `form:"appID"`
	AppSecret string `form:"appSecret"`
	Uid       string `form:"uid"`
}
