package mail

import (
	"fmt"
	"log"
	"net/smtp"
	"strings"
	"wnadmin/app/config"
)

// Send 把询盘发到指定邮箱。未配置发信服务器或收件邮箱时直接返回。
func Send(to, subject, body string) error {
	cfg := config.App.Mail
	to = strings.TrimSpace(to)
	if strings.TrimSpace(cfg.Host) == "" || to == "" {
		log.Printf("询盘已保存，邮件未发送：发信服务器或收件邮箱还没填")
		return nil
	}
	port := cfg.Port
	if port == 0 {
		port = 587
	}
	from := strings.TrimSpace(cfg.From)
	if from == "" {
		from = cfg.User
	}
	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n" +
		"\r\n" + body)
	auth := smtp.PlainAuth("", cfg.User, cfg.Password, cfg.Host)
	err := smtp.SendMail(fmt.Sprintf("%s:%d", cfg.Host, port), auth, from, []string{to}, msg)
	if err != nil {
		log.Printf("询盘邮件发送失败: %v", err)
	}
	return err
}
