package utils

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"wnadmin/app/provider/mylog"

	"github.com/forgoer/openssl"
	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/google/uuid"
	"github.com/parnurzeal/gorequest"
)

func GenerateFileName(originalName string) string {
	ext := filepath.Ext(originalName)
	fileName := strings.TrimSuffix(originalName, ext)
	return fileName + "_" + RandomString(8) + ext
}

// RandomString 生成给定长度的随机字符串
func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
func GetContentTypeByExtension(fileName string) (string, error) {
	extension := filepath.Ext(fileName)
	extension = strings.ToLower(extension)
	contentTypes := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".txt":  "text/plain",
		".pdf":  "application/pdf",
		// 添加更多文件类型的映射
	}

	contentType, found := contentTypes[extension]
	if !found {
		return "", fmt.Errorf("unsupported file extension: %s", extension)
	}

	return contentType, nil
}
func HttpBuildQuery(params map[string]interface{}) (paramStr string) {
	// sort keys
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	values := url.Values{}
	for _, key := range keys {
		values.Set(key, gconv.String(params[key]))
	}
	paramStr = values.Encode()

	return paramStr
}
func DesECBEncrypt(params map[string]interface{}) string {
	encryptKey := []byte("ebfc8cc9e3b4111142049be708c3b07c")
	if len(encryptKey) > 8 {
		encryptKey = encryptKey[:8]
	}
	queryString := HttpBuildQuery(params)
	signature := openssl.Md5(queryString)
	fmt.Println(hex.EncodeToString(signature))

	// replace sign key using signature and encode to a json string
	params["sign_key"] = strings.ToUpper(hex.EncodeToString(signature))
	paramsFinalByte, _ := json.Marshal(params) // auto sort by key

	paramsEncrypted, _ := openssl.DesECBEncrypt(paramsFinalByte, encryptKey, openssl.PKCS5_PADDING)
	return base64.StdEncoding.EncodeToString(paramsEncrypted)
}
func DesECBDecrypt(x string) []byte {
	encryptKey := []byte("ebfc8cc9e3b4111142049be708c3b07c")
	if len(encryptKey) > 8 {
		encryptKey = encryptKey[:8]
	}
	data, _ := base64.StdEncoding.DecodeString(x)
	paramsDecrypted, _ := openssl.DesECBDecrypt(data, encryptKey, openssl.PKCS5_PADDING)
	return paramsDecrypted
}
func Ptr[T any](value T) *T {
	return &value
}

// getClientIP returns the real client IP address from the given HTTP request.
func GetClientIP(r *http.Request) string {
	// Check the X-Forwarded-For header first, as it may contain the original client IP.
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		// Split the comma-separated list and take the first (leftmost) IP address.
		ips := strings.Split(xff, ",")
		ip := ips[0]
		return strings.TrimSpace(ip)
	}
	// If X-Forwarded-For is not present or empty, fall back to RemoteAddr.
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return ""
	}
	return ip
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

// SendHttpPost 发送HTTP POST请求并返回响应
func SendHttpPost(reqUrl string, requestBody map[string]string) (map[string]interface{}, error) {
	queryString := concatenateURLParams(requestBody)
	request := gorequest.New()
	resp, body, errs := request.Post(reqUrl).
		Send(queryString).
		End()
	mylog.Log.Info("resp", resp, "body", body, "errs", errs)
	result := map[string]interface{}{}
	json.Unmarshal([]byte(body), &result)
	return result, nil
}
func MyHash(data map[string]interface{}) string {
	jsonData, _ := json.Marshal(data)
	jsonString := string(jsonData)
	mylog.Log.Info("jsonString", jsonString)
	encodedData := encodeChineseCharacters(jsonString)
	reversedData := reverseString(encodedData)
	encryptedData := encryptData(reversedData)
	base64Data := base64Encode(encryptedData)

	return base64Data
}
func encodeChineseCharacters(str string) string {
	str = strings.ReplaceAll(str, "%", "%25")
	encoded := ""
	for _, char := range str {
		if int(char) >= 128 {
			encoded += url.QueryEscape(string(char))
		} else {
			encoded += string(char)
		}
	}
	return encoded
}

func reverseString(str string) string {
	runes := []rune(str)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func encryptData(data string) string {
	encrypted := ""
	for _, char := range data {
		encrypted += string(int(char) ^ 99)
	}
	return encrypted
}

func base64Encode(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}

// GetCurrentUsername 获取当前用户名
func GetCurrentUsername(c *gin.Context) string {
	if username, exists := c.Get("username"); exists {
		return username.(string)
	}
	return ""
}

// GetBaseURL 获取基础URL
func GetBaseURL(c *gin.Context) string {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	return scheme + "://" + c.Request.Host
}

// ParseParamID 解析路径参数ID
func ParseParamID(c *gin.Context) int64 {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	return id
}

// GenerateUUID 生成UUID
func GenerateUUID() string {
	return uuid.New().String()
}

// SaveUploadedFile 保存上传的文件
func SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// 确保目标目录存在
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

// RemoveFile 删除文件
func RemoveFile(filePath string) error {
	return os.Remove(filePath)
}

// StrToInt64 字符串转int64
func StrToInt64(str string) int64 {
	if str == "" {
		return 0
	}
	val, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0
	}
	return val
}
