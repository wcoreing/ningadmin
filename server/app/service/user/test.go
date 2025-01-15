package user

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/v2/util/gconv"
	"wnadmin/app/common/myhttp"
	"wnadmin/app/common/response"
	"wnadmin/app/common/utils"
	"wnadmin/app/global"
	"wnadmin/app/model/dto"
)

func LoginTest(c *gin.Context, params []byte) {
	if global.CONFIG.App.ENV == "prod" {
		response.FailWithMessage(202, "login test only support dev", c)
		return
	}
	loginDto := dto.LoginDto{}
	json.Unmarshal(params, &loginDto)
	if loginDto.AppID == "" || loginDto.AppSecret == "" {
		response.FailWithMessage(202, "appID or appSecret is empty", c)
		return
	}
	account := fmt.Sprintf("logintest%s", utils.SnowflakeID())
	// 第一次请求
	dataRegister := map[string]interface{}{
		"action":    "11",
		"appID":     loginDto.AppID,
		"appSecret": loginDto.AppSecret,
		"uid":       account,
		"name":      account,
	}
	dataRegister["x"] = utils.DesECBEncrypt(dataRegister)
	body := gconv.MapStrStr(dataRegister)
	host := fmt.Sprintf("http://127.0.0.1:%d/entrance", global.CONFIG.App.Port)
	registerResBytes, err := myhttp.Post(host, body)
	if err != nil {
		response.FailWithMessage(202, fmt.Sprintf("register error,%s", err.Error()), c)
		return
	}
	res := response.Response{}
	json.Unmarshal(registerResBytes, &res)
	res.Data = account
	if res.ReturnCode == 0000 {
		response.OkWithData(res, c)
		// 第二次请求
		dataLogin := map[string]interface{}{
			"action":    "12",
			"appID":     loginDto.AppID,
			"appSecret": loginDto.AppSecret,
			"uid":       account,
			"gameCode":  "slots-MajiongWays2",
		}
		dataLogin["x"] = utils.DesECBEncrypt(dataLogin)
		body := gconv.MapStrStr(dataLogin)
		loginResBytes, err := myhttp.Post(host, body)
		if err != nil {
			response.FailWithMessage(202, fmt.Sprintf("register error,%s,%s", account, string(loginResBytes)), c)
			return
		}
		json.Unmarshal(loginResBytes, &res)
		if res.ReturnCode == 0000 {
			response.OkWithData(res, c)
			return
		} else {
			response.FailWithMessage(203, fmt.Sprintf("login error,%s,%s", account, string(loginResBytes)), c)
			return
		}
	} else {
		response.FailWithMessage(204, fmt.Sprintf("register error,%s,%s", account, string(registerResBytes)), c)
		return
	}
}
