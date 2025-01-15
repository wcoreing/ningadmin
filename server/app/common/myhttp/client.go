package myhttp

import (
	"fmt"
	"github.com/parnurzeal/gorequest"
	"net/url"
	"strings"
	"wnadmin/app/global"
)

func Post(reqUrl string, requestBody map[string]string) ([]byte, error) {
	global.Log.Info("reqUrl", reqUrl, "requestBody", requestBody)
	queryString := concatenateURLParams(requestBody)
	request := gorequest.New()
	resp, body, errs := request.Post(reqUrl).
		Send(queryString).
		End()
	if errs != nil {
		global.Log.Error("errs", errs)
		return nil, errs[0]
	}
	global.Log.Info("resp", resp, "body", body, "errs", errs)
	return []byte(body), nil
}
func concatenateURLParams(params map[string]string) string {
	var queryString []string

	for k, v := range params {
		// URL编码键和值，确保特殊字符被正确处理
		encodedKey := url.QueryEscape(k)
		encodedValue := url.QueryEscape(v)
		// 拼接成 "key=value" 形式的字符串并添加到结果列表中
		queryString = append(queryString, fmt.Sprintf("%s=%s", encodedKey, encodedValue))
	}

	// 使用 "&" 连接所有参数对
	return strings.Join(queryString, "&")
}
