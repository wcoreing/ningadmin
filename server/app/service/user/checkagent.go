package user

import (
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/util/gconv"
	"gorm.io/gorm"
	"wnadmin/app/common/utils"
	"wnadmin/app/global"
	"wnadmin/app/model/bo"
	"wnadmin/app/model/do/LoginDB"
	"wnadmin/app/model/do/QPGameUserDB"
	"wnadmin/app/model/do/SMAdminDB"

	"time"
)

func CheckAgent(checkAgentBo bo.CheckAgentBo) (agent QPGameUserDB.AccountsPlatformDo, err error) {
	agent = QPGameUserDB.AccountsPlatformDo{}
	if err := global.DB.QPGameUserDB.Table("Accounts_platform").Session(&gorm.Session{PrepareStmt: true}).
		Where("appID", checkAgentBo.AppID).
		Where("appSecret", checkAgentBo.AppSecret).
		Where("isDel", 1).First(&agent).Error; err != nil {
		return agent, err
	}
	global.DB.RedisGame.Set(fmt.Sprintf("agentInfo_%s", agent.ID), agent, time.Second*60)
	//todo 设置代理访问数
	offline, _ := global.DB.RedisGame.Get("OpenApi:SystemOffline").Int()
	if offline == 1 {
		SystemConfigDo := SMAdminDB.SystemConfigDo{}
		if err := SystemConfigDo.Table().Where("name", "testAgentIds").First(&SystemConfigDo).Error; err == nil {
			var ids []string
			json.Unmarshal([]byte(SystemConfigDo.Value), &ids)
			if garray.NewStrArrayFrom(ids).Contains(gconv.String(agent.ID)) == false {
				return agent, fmt.Errorf("system maintaining")
			}
		}

	}
	needServerUrlActions := garray.NewStrArrayFrom([]string{"login", "demoLogin", "kickUser", "transfer", "12", "14", "37", "15", "specifyDesktop"})
	if needServerUrlActions.Contains(checkAgentBo.Action) {
		if checkAgentBo.Uid == "" {
			RegionInfoDo := QPGameUserDB.RegionInfoDo{}
			if err := RegionInfoDo.Table().Where("isOpen", 1).
				Where("MoneyType", agent.MoneyType).
				Order("ConnectCount").
				First(&RegionInfoDo).Error; err != nil {
				return agent, fmt.Errorf("Server is not open at all")
			}
			agent.RegionID = utils.Ptr(RegionInfoDo.RegionID)

		} else {
			accountsInfoDo := QPGameUserDB.AccountsInfoDo{}
			if err := accountsInfoDo.Table().Where("Accounts", checkAgentBo.Uid).
				Where("Platform_Id", agent.ID).
				First(&accountsInfoDo).Error; err != nil {
				return agent, fmt.Errorf("user does not exist,when get by agent.")
			}
			if regionID, err := GetUserReginID(gconv.Int(accountsInfoDo.UserID), agent.MoneyType); err != nil {
				return agent, err
			} else {
				agent.RegionID = &regionID
			}
		}
		relayServerInfo := LoginDB.RelayServerInfoDo{}
		if err := relayServerInfo.Table().Where("Using", 1).
			Where("MoneyType", agent.MoneyType).
			Where("RegionID", agent.RegionID).First(&relayServerInfo).Error; err != nil {
			return agent, fmt.Errorf("Interface service not open.")
		}
		agent.ServerUrl = utils.Ptr(fmt.Sprintf("http://%s:%d", relayServerInfo.IP, relayServerInfo.Port))
	}
	return agent, nil
}
func GetUserReginID(uid, moneyType int) (int, error) {
	regionIDkey := fmt.Sprintf("AccountRegion:%d", uid)
	regionIDLockkey := fmt.Sprintf("LockAccountRegion:%d", uid)
	RegionID, _ := global.DB.RedisGame.Get(regionIDkey).Int()
	if RegionID != 0 {
		RegionInfoDo := QPGameUserDB.RegionInfoDo{}
		if err := RegionInfoDo.Table().Where("isOpen", 1).
			Where("MoneyType", moneyType).
			Where("RegionID", RegionID).
			First(&RegionInfoDo).Error; err == nil {
			if global.DB.RedisGame.TTL(regionIDkey).Val() != -1 {
				global.DB.RedisGame.Set(regionIDkey, RegionID, time.Second*300)
			}
			return RegionID, nil
		}

	}
	lock := global.DB.RedisGame.SetNX(regionIDLockkey, 1, time.Second*3)
	defer func() {
		global.DB.RedisGame.Del(regionIDLockkey)
	}()
	if lock != nil {
		RegionInfoDo := QPGameUserDB.RegionInfoDo{}
		if err := RegionInfoDo.Table().Where("isOpen", 1).
			Where("MoneyType", moneyType).
			Order("ConnectCount").
			First(&RegionInfoDo).Error; err != nil {
			return 0, fmt.Errorf("Server is not open at all")
		}
		global.DB.RedisGame.Set(regionIDkey, RegionInfoDo.RegionID, time.Second*300)
		return RegionInfoDo.RegionID, nil
	} else {
		regionID := 0
		for {
			if err := global.DB.RedisGame.Get(regionIDkey).Err(); err == nil {
				// 从缓存中获取数据
				regionID, err = global.DB.RedisGame.Get(regionIDkey).Int()
				if err == nil {
					break
				}
			}

			// 等待一段时间后再次尝试获取缓存数据
			time.Sleep(time.Second)
		}
		return regionID, nil
	}

}
