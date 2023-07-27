package constant

import (
	"gopkg.in/gomail.v2"
	"strconv"
)

func SendMailSelf(body string) {
	userName := "XXXX@163.com"
	authCode := "XXXXX"
	host := "smtp.163.com"
	portStr := "465"
	mailTo := "XXXX@qq.com"
	sendName := "XXXX@163.com"
	subject := "服务器重启"
	SendMail(userName, authCode, host, portStr, mailTo, sendName, subject, body)
}

func SendMail(userName, authCode, host, portStr, mailTo, sendName string, subject, body string) error {
	port, _ := strconv.Atoi(portStr)
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(userName, sendName))
	m.SetHeader("To", mailTo)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	d := gomail.NewDialer(host, port, userName, authCode)
	err := d.DialAndSend(m)
	return err
}
