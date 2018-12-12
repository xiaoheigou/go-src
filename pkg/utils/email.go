package utils

import (
	"fmt"
	"net/smtp"
	"strings"
)

func sendToMail(user, password, host, to, subject, body string) error {
	hostPort := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hostPort[0])

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = user
	headers["To"] = to
	headers["Subject"] = subject
	headers["Content-Type"] = "text/html; charset=UTF-8"

	// Setup message
	message := ""
	for k,v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	msg := []byte(message)
	sendTo := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, user, sendTo, msg)
	return err
}

func SendRandomCodeToMail(to, randomCode, timeout string) error {
	var user = Config.GetString("email.sender")
	if ! IsValidEmail(user) {
		Log.Errorln("Wrong configuration: email.sender [%v], not a valid email addr.", user)
		return nil
	}
	var password = Config.GetString("email.password")
	if len(password) == 0 {
		Log.Errorln("Wrong configuration: email.password [%v], it empty.", password)
		return nil
	}
	var smtpSvr = Config.GetString("email.smtpsvr")
	if strings.Contains(smtpSvr, ":") {
		Log.Errorln("Wrong configuration: email.host [%v], must contains host and port.", smtpSvr)
		return nil
	}

	subject := "您的Yuudidi注册随机码"

	body := `
		<html>
		<body>
		<h3>
		您的Yuudidi注册随机码为：  ` + randomCode + ` ，有效期 ` + timeout + ` 分钟，请您尽快验证。
		</h3>
		</body>
		</html>
		`
	fmt.Println("send email")
	if err = sendToMail(user, password, smtpSvr, to, subject, body); err != nil {
		Log.Errorln("SendRandomCodeToMail fail [%v]", err)
	}
	return nil
}

