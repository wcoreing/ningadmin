package mylog

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"os"
	"time"
)

var Log *logrus.Logger

func Newlog(filename string) *logrus.Logger {
	_, err := os.Stat("./logs")
	if err != nil {
		if !os.IsExist(err) {
			err := os.MkdirAll("./logs", 0777)
			if err != nil {
				log.Fatalf("haverst.logger-> 创建日志目录 %s", "./logs")
			}
		}
	}

	lg := logrus.New()
	//logrus.SetFormatter(&logrus.JSONFormatter{})
	lg.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		DisableColors:   true,
		FullTimestamp:   true,
		DisableQuote:    true,
	},
	)
	year, month, _ := time.Now().Date()
	//writer1 := &bytes.Buffer{}
	writer2 := os.Stdout
	logfile := fmt.Sprintf("./logs/"+filename+"%d%d.log", year, month)
	writer3, _ := os.OpenFile(logfile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0755)
	lg.SetOutput(io.MultiWriter(writer2))
	lg.SetOutput(io.MultiWriter(writer3))
	//设置最低loglevel
	lg.SetLevel(logrus.InfoLevel)
	return lg
}
func init() {
	Log = Newlog("wnadmin")
}
