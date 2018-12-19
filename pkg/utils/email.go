package utils

import (
	"errors"
	"fmt"
	"gopkg.in/gomail.v2"
	"net"
	"strconv"
	"strings"
)

func sendMailUsingGomail(user, password, host, to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", user)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	onlyHost, onlyPort, _ := net.SplitHostPort(host)
	var port int
	var err error
	if port, err = strconv.Atoi(onlyPort); err != nil {
		Log.Errorln("strconv.Atoi error [%s]", err)
		return err
	}

	d := gomail.NewDialer(onlyHost, port, user, password)

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		Log.Errorln("DialAndSend error [%s]", err)
		return err
	}
	return nil
}

func SendRandomCodeToMail(to, randomCode, timeout string) error {
	var user = Config.GetString("email.sender")
	if ! IsValidEmail(user) {
		Log.Errorln("Wrong configuration: email.sender [%v], not a valid email addr.", user)
		return errors.New("wrong configuration: email.sender")
	}
	var password = Config.GetString("email.password")
	if len(password) == 0 {
		Log.Errorln("Wrong configuration: email.password [%v], it empty.", password)
		return errors.New("wrong configuration: email.password")
	}
	var smtpSvr = Config.GetString("email.smtpsvr")
	if ! strings.Contains(smtpSvr, ":") {
		Log.Errorln("Wrong configuration: email.smtpsvr [%v], must contains host and port.", smtpSvr)
		return errors.New("wrong configuration: email.smtpsvr")
	}

	subject := "您的Yuudidi注册验证码"

	body := `
		<html>
		<body>
		<h3>
		您的Yuudidi注册验证码为：  ` + randomCode + ` ，有效期 ` + timeout + ` 分钟，请您尽快验证。
		</h3>
		</body>
		</html>
		`
	fmt.Println("send email")
	if err := sendMailUsingGomail(user, password, smtpSvr, to, subject, body); err != nil {
		Log.Errorln("SendRandomCodeToMail fail [%v]", err)
		return err
	}
	Log.Infoln("Send email to [%v] success", to)
	return nil
}
