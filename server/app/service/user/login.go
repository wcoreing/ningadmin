package user

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/util/gconv"
	"strings"
	"wnadmin/app/common/response"
	"wnadmin/app/common/utils"
	"wnadmin/app/global"
	"wnadmin/app/model/bo"
	"wnadmin/app/model/do/QPGameUserDB"
	"wnadmin/app/model/do/QPServerInfoDB"
	"wnadmin/app/model/dto"
)

func Login(c *gin.Context, params []byte) {
	var loginDto dto.LoginDto
	json.Unmarshal(params, &loginDto)
	var agentInfo QPGameUserDB.AccountsPlatformDo
	var err error
	var checkAgent bo.CheckAgentBo
	json.Unmarshal(params, &checkAgent)
	if agentInfo, err = CheckAgent(checkAgent); err != nil {
		response.FailWithMessage(8004, err.Error(), c)
		return
	}

	langs := garray.NewStrArrayFrom([]string{"cn", "en", "vn", "th", "ko", "tw", "ja", "pt", "Bahasa", "spa"})

	if langs.Contains(loginDto.Lang) == false {
		loginDto.Lang = "en"
	}
	if loginDto.Uid == "" {
		response.FailWithMessage(2001, "uid is empty", c)
		return
	}
	if loginDto.GameCode == "" {
		response.FailWithMessage(2002, "gameCode is empty", c)
	}
	GameKindItemDo := QPServerInfoDB.GameKindItemDo{}
	if err = GameKindItemDo.Table().Where("gameCode", loginDto.GameCode).First(&GameKindItemDo).Error; err != nil {
		response.FailWithMessage(1010, "gameCode is error", c)
		return
	}
	gameListEnable := []string{}
	if err = json.Unmarshal([]byte(*agentInfo.GameListEnable), &gameListEnable); err != nil {
		response.FailWithMessage(1011, "gameListEnable is error", c)
		return
	}
	if !garray.NewStrArrayFrom(gameListEnable).Contains(gconv.String(GameKindItemDo.KindID)) {
		response.FailWithMessage(1011, fmt.Sprintf("Game is closed:%s", GameKindItemDo.KindID), c)
		return
	}
	if *agentInfo.UsedMoneyLimit >= *agentInfo.MoneyLimit && *agentInfo.MoneyLimit != 0 && agentInfo.AutoKick == true {
		response.FailWithMessage(1006, "Platform restrictions,someone reach limit", c)
		return
	}
	if agentInfo.IsLogin == 2 {
		response.FailWithMessage(1006, "Platform restrictions,agent dont support login", c)
		return
	}
	loginFrom := "login"
	token, err := getTokenFromServer(loginFrom, loginDto, agentInfo)
	if err != nil {
		response.FailWithMessage(1006, err.Error(), c)
		return
	}
	url, err := geturl(token, loginDto, agentInfo, GameKindItemDo)
	if err != nil {
		response.FailWithMessage(1006, err.Error(), c)
		return
	}
	result := map[string]string{
		"path": url,
	}
	response.OkWithData(result, c)

}
func getTokenFromServer(loginFrom string, loginDto dto.LoginDto, agentInfo QPGameUserDB.AccountsPlatformDo) (string, error) {
	var data map[string]string
	switch loginFrom {
	case "demo":
		data = map[string]string{
			"action":    "103",
			"MoneyType": gconv.String(agentInfo.MoneyType),
			"RegionID":  gconv.String(*agentInfo.RegionID),
		}
		break
	default:
		AccountsInfoDo := QPGameUserDB.AccountsInfoDo{}
		if err := AccountsInfoDo.Table().
			Where("Accounts", loginDto.Uid).
			Where("Platform_Id", agentInfo.ID).
			First(&AccountsInfoDo).Error; err != nil {
			return "", fmt.Errorf("agent of user not bind:%s:%s", loginDto.Uid, agentInfo.ID)
		}
		phoneType := "0"
		if loginDto.PhoneType != "" {
			phoneType = loginDto.PhoneType
		}
		data = map[string]string{
			"action":    "102",
			"MoneyType": gconv.String(agentInfo.MoneyType),
			"RegionID":  gconv.String(*agentInfo.RegionID),
			"UserID":    gconv.String(*AccountsInfoDo.UserID),
			"PhoneType": phoneType,
		}
		break
	}

	res, err := utils.SendHttpPost(*agentInfo.ServerUrl, data)

	if err != nil {
		return "", err
	}
	switch gconv.Int(res["returnCode"]) {
	case 0:
		token := res["data"].(map[string]interface{})["token"]
		if token != nil {
			//todo 上线通知
			return token.(string), nil
		}
	case 1:
		return "", fmt.Errorf("User does not exist.10021:%s", res["returnMsg"])
		break
	case 86:
		return "", fmt.Errorf("User does not exist.10022:%s", res["returnMsg"])
		break
	case 73:
		return "", fmt.Errorf("Platform restrictions.10061:%s", res["returnMsg"])
		break
	case 660:
		return "", fmt.Errorf("regionID error:%s", res["returnMsg"])

	}
	return "", fmt.Errorf("get token error:%s", res["returnMsg"])
}
func geturl(token string, loginDto dto.LoginDto, agentInfo QPGameUserDB.AccountsPlatformDo, GameKindItemDo QPServerInfoDB.GameKindItemDo) (url string, err error) {
	isClose := "1"
	if agentInfo.IsClose != "" {
		isClose = agentInfo.IsClose
	}
	payload := map[string]interface{}{}
	payload["token"] = token
	payload["lang"] = loginDto.Lang
	payload["userAccount"] = loginDto.Uid
	payload["agentIdPath"] = ""
	payload["ip"] = ""
	payload["rate"] = agentInfo.UserBankRate
	payload["loginUrl"] = global.CONFIG.App.GAME_LOGIN_URL
	payload["apiUrl"] = global.CONFIG.App.GAME_OPEN_URL
	payload["backUrl"] = ""
	payload["optional"] = loginDto.Optional
	payload["closeType"] = loginDto.CloseType
	payload["closeable"] = isClose
	payload["is_https"] = loginDto.IsHttps
	payload["moneytype"] = agentInfo.MoneyType
	payload["MoneyBag"] = agentInfo.MoneyBag
	payload["SeamlessWalletType"] = agentInfo.SeamlessWalletType
	payload["agentId"] = agentInfo.ID
	payload["logoGroupID"] = agentInfo.LogoGroupID
	payload["gamename"] = getGameNameByKindId(&GameKindItemDo, loginDto.Lang)
	payload["gameCode"] = GameKindItemDo.GameCode

	ServerList, _ := getGatewayIpByMoneyType(agentInfo.MoneyType, *agentInfo.RegionID)

	historyUrl := getHistoryUrl(agentInfo, GameKindItemDo, "", loginDto)
	payload["loginConfig"] = map[string]interface{}{
		"ServerList": ServerList,
		"HistoryUrl": historyUrl,
	}
	payload["ws"] = ServerList
	payload["history"] = historyUrl

	payload["liveVideoIsActive"] = 1
	payload["liveVideoUrl"] = fmt.Sprintf("%s/?agent_token=%s&item_id=%d&agent_id=%s", global.CONFIG.App.LIVE_VIDEO_WEB_URL, token, GameKindItemDo.KindID, global.CONFIG.App.LIVE_VIDEO_AGENT_ID)
	if agentInfo.GameRecommenBtn == true {
		payload["hotUrl"] = fmt.Sprintf("%s/hot?token=%s&lang=%s", global.CONFIG.App.RES_URL, token, loginDto.Lang)
	}

	logoGroup := QPGameUserDB.BrandDo{}
	payload["LoadingType"] = *GameKindItemDo.LoadingType
	payload["isStartPage"] = *GameKindItemDo.IsStartPage
	payload["isBuyFree"] = GameKindItemDo.IsBuyFree
	payload["bgmBody"] = fmt.Sprintf("%sres/loadingbg/%s/bg.jpg", global.CONFIG.App.GAME_URL, GameKindItemDo.GameCode)
	payload["bgmGame"] = fmt.Sprintf("%sres/loadingbg/%s/bg_qian.jpg", global.CONFIG.App.GAME_URL, GameKindItemDo.GameCode)
	payload["bgmBottom"] = fmt.Sprintf("%sres/loadingbg/bg_bottom.png", global.CONFIG.App.GAME_URL)
	payload["buyFree"] = fmt.Sprintf("%sres/loadingbg/%s/buyfree.png", global.CONFIG.App.GAME_URL, GameKindItemDo.GameCode)

	if err = global.DB.QPGameUserDB.Table("Brand").Joins("left join BrandTemplate on Brand.groupID = BrandTemplate.groupID and Brand.loading_type=BrandTemplate.loading_type").Where("Brand.groupID", *agentInfo.BrandID).First(&logoGroup).Error; err == nil {
		if *logoGroup.LoadingType == "jy2" {
			payload["bgmBottom"] = fmt.Sprintf("%sres/loadingbg/jy2_bottom.png")
		}
		payload["logoGroup"] = logoGroup
	}

	gameAddress := QPServerInfoDB.GameAddressDo{}
	if err := gameAddress.Table().First(&gameAddress).Error; err != nil {
		return "", fmt.Errorf("gameaddress error")
	}
	host := gameAddress.URL
	if loginDto.IsHttps == "1" {
		host = strings.Replace(host, "http://", "https://", 1)
	}

	token = utils.MyHash(payload)
	if payload["logoGroup"] != nil {
		url := fmt.Sprintf("%slogogroup/indexStyle/%s/index.html?token=%s", host, *logoGroup.LoadingType, token)
		return url, nil
	}
	switch GameKindItemDo.GameCode {
	case "slots-MajiongWays", "slots-MajiongWays2":
		url = host + "logogroup/indexStyle/pg/index.html?token=" + token
		break
	default:
		switch payload["LoadingType"].(string) {
		case "1":
			url = host + "logogroup/indexStyle/pt/index.html?token=" + token
			break
		case "2":
			url = host + "logogroup/indexStyle/ptw/index.html?token=" + token
			break
		default:
			token = utils.DesECBEncrypt(payload)
			url = fmt.Sprintf("%s%s/web-mobile/index.html?closeable=%s&closeType=%s&lang=%s&gamename=%s&token=%s", host, GameKindItemDo.GameCode, isClose, loginDto.CloseType, loginDto.Lang, GameKindItemDo.GameName, token)
			break
		}
	}

	return url, nil
}
func getGameNameByKindId(GameKindItemDo *QPServerInfoDB.GameKindItemDo, lang string) string {
	GameI18nDo := QPServerInfoDB.GameI18nDo{}
	if err := GameI18nDo.Table().Where("KindID", GameKindItemDo.KindID).First(&GameI18nDo).Error; err != nil {
		fmt.Println("getGameNameByKindId error")
	}
	names := map[string]string{}
	if err := json.Unmarshal([]byte(GameI18nDo.Name), &names); err != nil {
		fmt.Println("getGameNameByKindId error")
	}
	var name string
	if names[lang] == "" {
		name = names["en"]
	} else {
		name = names[lang]
	}
	GameKindItemDo.GameName = name
	return name
}

func getGatewayIpByMoneyType(moneyType int, regionID int) ([]string, error) {

	var GatewayIP string
	var Port int
	var IsHttps int
	var ProxyUrl string

	sql := fmt.Sprintf("EXEC HttpCheck2 @MoneyType=%d, @RegionID=%d", moneyType, regionID)

	global.DB.Login.Raw(sql).Row().Scan(&GatewayIP, &Port, &IsHttps, &ProxyUrl)
	var gatewayIPs []string

	if ProxyUrl != "" {
		gatewayIPs = append(gatewayIPs, ProxyUrl)
	} else if IsHttps == 1 {
		gatewayIPs = append(gatewayIPs, fmt.Sprintf("wss://%s:%d", GatewayIP, Port))
	} else {
		gatewayIPs = append(gatewayIPs, fmt.Sprintf("ws://%s:%d", GatewayIP, Port))
	}

	return gatewayIPs, nil
}
func getHistoryUrl(agentInfo QPGameUserDB.AccountsPlatformDo, gameInfo QPServerInfoDB.GameKindItemDo, token string, loginDto dto.LoginDto) string {
	gameCodes := garray.NewIntArrayFrom([]int{500, 501, 502, 503, 506, 509, 510, 511, 512, 513, 520, 517, 519, 520, 523, 524, 505, 521, 526})
	historyUrl := global.CONFIG.App.HISTORY_URL
	if gameCodes.Contains(gameInfo.KindID) {
		historyUrl = fmt.Sprintf("%s%s", historyUrl, gameInfo.GameCode)
	}
	if token != "" {
		historyUrl = fmt.Sprintf("%s?token=%s&kindid=%d&lang=%s&gamename=%s&moneytype=%d", historyUrl, token, gameInfo.KindID, loginDto.Lang, gameInfo.GameName, agentInfo.MoneyType)
	}
	return historyUrl
}
