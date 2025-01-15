package user

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"time"
	"wnadmin/app/common/response"
	"wnadmin/app/common/utils"
	"wnadmin/app/global"
	"wnadmin/app/model/bo"
	"wnadmin/app/model/do/QPGameUserDB"
	"wnadmin/app/model/dto"
	"wnadmin/app/provider/register"
)

func Register(c *gin.Context, params []byte) {
	var registerDto dto.RegisterDto
	json.Unmarshal(params, &registerDto)
	global.Log.Info("register:" + registerDto.Uid)
	if registerDto.Uid == "" {
		response.FailWithMessage(2001, "The user parameter is empty", c)
		return
	}
	if registerDto.AppID == "" || registerDto.AppSecret == "" {
		response.FailWithMessage(8003, "The appID or appSecret parameter is empty", c)
		return
	}
	if registerDto.Name != "" {
		gbkEncoder := simplifiedchinese.GBK.NewEncoder()
		registerDto.Name, _, _ = transform.String(gbkEncoder, registerDto.Name)
	} else {
		registerDto.Name = registerDto.Uid
	}
	var agent QPGameUserDB.AccountsPlatformDo
	var err error
	var checkAgent bo.CheckAgentBo
	json.Unmarshal(params, &checkAgent)
	agent, err = CheckAgent(checkAgent)
	if err != nil {
		response.FailWithMessage(8004, "The appID or appSecret parameter is error", c)
		return
	}
	AccountsInfo := &QPGameUserDB.AccountsInfoDo{}
	AccountsInfo.Table().Where("Accounts", registerDto.Uid).Where("Platform_Id", agent.ID).First(&AccountsInfo)
	if AccountsInfo.UserID != nil {
		response.FailWithMessage(8015, fmt.Sprintf("Register failed,Accounts exits:%s", registerDto.Uid), c)
		return
	}
	AccountsInfoDo := &QPGameUserDB.AccountsInfoDo{
		Accounts:     registerDto.Uid,
		NickName:     &registerDto.Name,
		LogonPass:    utils.Ptr("httpaccount"),
		OrAccount:    utils.Ptr(registerDto.Uid),
		MoneyType:    agent.MoneyType,
		ChannelID:    utils.Ptr(agent.ID),
		PlatformID:   agent.ID,
		RegisterIP:   utils.Ptr(utils.GetClientIP(c.Request)),
		RegisterDate: time.Now(),
		MachineNo:    utils.Ptr("00000000000000000000000000000000"),
	}
	registerSignal := make(chan bool)
	register.RegisterHandler.Enqueue(c, *AccountsInfoDo, registerSignal)
	<-registerSignal
	response.OkWithData(AccountsInfoDo, c)
}
