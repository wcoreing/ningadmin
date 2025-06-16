package mylog

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	// Log 全局日志实例
	Log *logrus.Logger
	// 默认日志级别
	defaultLogLevel = logrus.InfoLevel
)

// 自定义格式化器
type customFormatter struct {
	logrus.TextFormatter
}

// Format 自定义日志格式
func (f *customFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// 获取调用者信息
	var file string
	var line int
	var function string

	if entry.HasCaller() {
		file = filepath.Base(entry.Caller.File)
		line = entry.Caller.Line
		function = filepath.Base(entry.Caller.Function)
	}

	// 日志时间格式
	timestamp := entry.Time.Format("2006-01-02 15:04:05.000")

	// 获取请求ID
	requestID := ""
	if id, exists := entry.Data["request_id"]; exists {
		requestID = fmt.Sprintf("[ReqID:%s]", id)
	}

	// 格式化其他字段
	var fields string
	if len(entry.Data) > 0 {
		items := make([]string, 0, len(entry.Data))
		for k, v := range entry.Data {
			if k != "request_id" {
				items = append(items, fmt.Sprintf("%s=%v", k, v))
			}
		}
		if len(items) > 0 {
			fields = " " + strings.Join(items, " ")
		}
	}

	// 构建日志消息
	var caller string
	if file != "" {
		caller = fmt.Sprintf("[%s:%d %s]", file, line, function)
	}

	logMsg := fmt.Sprintf("time=%s level=%s%s %s%s %s\n",
		timestamp, entry.Level, requestID, caller, fields, entry.Message)

	return []byte(logMsg), nil
}

// LogLevel 解析日志级别
func LogLevel(level string) logrus.Level {
	switch strings.ToLower(level) {
	case "trace":
		return logrus.TraceLevel
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warn":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	case "fatal":
		return logrus.FatalLevel
	case "panic":
		return logrus.PanicLevel
	default:
		return defaultLogLevel
	}
}

// NewLogger 创建新的日志实例
func NewLogger(filename string, level string) *logrus.Logger {
	// 创建日志目录
	logDir := "./logs"
	if err := os.MkdirAll(logDir, 0777); err != nil {
		log.Fatalf("创建日志目录失败: %v", err)
	}

	// 创建日志实例
	logger := logrus.New()

	// 设置格式化器
	logger.SetFormatter(&customFormatter{
		TextFormatter: logrus.TextFormatter{
			DisableColors:   true,
			TimestampFormat: "2006-01-02 15:04:05.000",
			FullTimestamp:   true,
			DisableQuote:    true,
		},
	})

	// 设置调用者信息
	logger.SetReportCaller(true)

	// 设置日志级别
	logger.SetLevel(LogLevel(level))

	// 生成日志文件名
	date := time.Now().Format("20060102")
	logFile := fmt.Sprintf("%s/%s_%s.log", logDir, filename, date)

	// 打开日志文件
	fileWriter, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("创建日志文件失败: %v", err)
	}

	// 同时输出到控制台和文件
	logger.SetOutput(io.MultiWriter(os.Stdout, fileWriter))

	return logger
}

// WithRequestID 创建带有请求ID的日志条目
func WithRequestID(requestID string) *logrus.Entry {
	return Log.WithField("request_id", requestID)
}

// GetCaller 获取调用者信息
func GetCaller(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown"
	}
	return fmt.Sprintf("%s:%d", filepath.Base(file), line)
}

// SetLogLevel 动态设置日志级别
func SetLogLevel(level string) {
	Log.SetLevel(LogLevel(level))
}

func init() {
	// 默认使用Info级别
	Log = NewLogger("wnadmin", "info")
}
