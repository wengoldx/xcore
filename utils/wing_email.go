// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2019/05/22   yangping       New version
// -------------------------------------------------------------------

package utils

import (
	"strings"

	"gopkg.in/gomail.v2"
)

// MailAgent mail agent informations
//
// `Useage`
//
// ---
//
//	mailTempleteSubject = "You hava a mail！"
//	mailTemplateFormat  = `
//		<p>
//			<span font-weight:bold; style="font-size:16px; color:#363636">Dear</span><br><br>
//			<span style="font-size:14px; color:#484848">Your account %s have not activate, please click the follow link to activate it.</span><br>
//			<span style="font-size:14px; color:#484848">%s</span><br>
//		</p>
//		<p align="right">
//			<span style="font-size:12px; color:#484848">From %s</span><br>
//			<span style="font-size:10px; color:#636363">%s</span>
//		</p>
//		`
//
//	 mailagent = &utils.NewMailAgent(account, password, smtphost, smtpport)
//	 subject := mailTempleteSubject
//	 message := fmt.Sprintf(mailTemplateFormat, account, link, who,
//	     time.Now().Format(templateTimeFormat))
//	 // send mail with attachment
//	 // return ma.SendMail(to, subject, message, fileName)
//	 return ma.SendMail(to, subject, message)
type MailAgent struct {
	acc  string // username - mail address
	pwd  string // account password
	host string // stmp/pop3 host
	port int    // stmp/pop3 port
}

// EmailContent email template
type EmailContent struct {
	Subject string // email title or subject
	Body    string // email body content
}

// Create a MailAgent object for send email.
func NewMailAgent(acc, pwd, host string, port int) *MailAgent {
	return &MailAgent{acc: acc, pwd: pwd, host: host, port: port}
}

// SendMail send email by mail account, it may set attachment from local file
func (a *MailAgent) SendMail(to []string, subject, body string, attach ...string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", a.acc)
	m.SetHeader("To", to[0])
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	if att := Variable(attach, ""); att != "" {
		m.Attach(att)
	}

	d := gomail.NewDialer(a.host, a.port, a.acc, a.pwd)
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

// SendCode send verify email with code
//
// The SMS templetes same as:
//
// ---
//
//	TplEmailRegister = EmailContent{"Account Verify Of XXX", `
//	<html>
//	  <body>
//	    <h3> Dear NAME </h3>
//	    <p> Thank you for register XXX, the registration verification code is : <h3> TOKEN </h3>, please activate your account in time.</br>
//	        Please DO NOT forward this code to others. If not myself, please delete this email.</p>
//	    </br>
//	    <h5>XXX Technology Co., Ltd</h5>
//	 </body>
//	</html>`}
func (a *MailAgent) SendCode(email EmailContent, mailto string, code string) error {
	to := []string{mailto}
	body := strings.Replace(email.Body, "NAME", mailto, 1)
	body = strings.Replace(body, "TOKEN", code, 1)
	return a.SendMail(to, email.Subject, body)
}

// SendFormat send mail with formated map
func (a *MailAgent) SendFormat(email EmailContent, mailto string, format ...map[string]string) error {
	return a.SendAttach(email, mailto, "", format...)
}

// SendAttach send mail with formated map and attach file
func (a *MailAgent) SendAttach(email EmailContent, mailto, attach string, format ...map[string]string) error {
	to := []string{mailto}
	body := email.Body
	if len(format) > 0 {
		for key, content := range format[0] {
			if key == "" || content == "" {
				continue
			}
			body = strings.Replace(body, key, content, 1)
		}
	}
	return a.SendMail(to, email.Subject, body, attach)
}
