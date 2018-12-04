// Copyright 2013 wetalk authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package mailer

import (
	"encoding/base64"
	"fmt"
	"net/smtp"
	"strings"
	"usercenter/setting"

	"crypto/tls"
	"net"

	"github.com/astaxie/beego"
)

type Message struct {
	To      []string
	From    string
	Subject string
	Body    string
	User    string
	Type    string
	Info    string
}

// create mail content
func (m Message) Content() string {
	// set mail type
	contentType := "text/plain; charset=UTF-8"
	if m.Type == "html" {
		contentType = "text/html; charset=UTF-8"
	}
	bs64 := base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")
	subject := fmt.Sprintf("=?UTF-8?B?%s?=", bs64.EncodeToString([]byte(m.Subject)))
	// create mail content
	content := "From: " + m.User + "<" + m.From +
		">\r\nSubject: " + subject + "\r\nContent-Type: " + contentType + "\r\n\r\n" + m.Body
	return content
}

func Send(msg Message) (int, error) {
	host, _, _ := net.SplitHostPort(setting.MailHost)

	// get message body
	content := msg.Content()

	auth := smtp.PlainAuth("", setting.MailAuthUser, setting.MailAuthPass, host)

	if len(msg.To) == 0 {
		return 0, fmt.Errorf("empty receive emails")
	}

	if len(msg.Body) == 0 {
		return 0, fmt.Errorf("empty email body")
	}

	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	conn, err := tls.Dial("tcp", setting.MailHost, tlsconfig)
	if err != nil {
		return 0, err
	}

	defer conn.Close()

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return 0, err
	}

	defer c.Quit()
	if err = c.Auth(auth); err != nil {
		return 0, err
	}

	if err = c.Mail(msg.From); err != nil {
		return 0, err
	}

	for _, to := range msg.To {
		if err = c.Rcpt(to); err != nil {
			return 0, err
		}
	}

	w, err := c.Data()
	if err != nil {
		return 0, err
	}

	body := []byte("To: " + strings.Join(msg.To, ";") + "\r\n" + content)

	if _, err = w.Write(body); err != nil {
		return 0, err
	}

	err = w.Close()
	if err != nil {
		return 0, err
	}
	return 1, c.Quit()
}

// Direct Send mail message
/*
func Send(msg Message) (int, error) {
	host := strings.Split(setting.MailHost, ":")

	// get message body
	content := msg.Content()

	auth := smtp.PlainAuth("", setting.MailAuthUser, setting.MailAuthPass, host[0])

	if len(msg.To) == 0 {
		return 0, fmt.Errorf("empty receive emails")
	}

	if len(msg.Body) == 0 {
		return 0, fmt.Errorf("empty email body")
	}

	if msg.Massive {
		// send mail to multiple emails one by one
		num := 0
		for _, to := range msg.To {
			body := []byte("To: " + to + "\r\n" + content)
			err := smtp.SendMail(setting.MailHost, auth, msg.From, []string{to}, body)
			if err != nil {
				return num, err
			}
			num++
		}
		return num, nil
	} else {
		body := []byte("To: " + strings.Join(msg.To, ";") + "\r\n" + content)

		// send to multiple emails in one message
		err := smtp.SendMail(setting.MailHost, auth, msg.From, msg.To, body)
		if err != nil {
			return 0, err
		} else {
			return 1, nil
		}
	}
}
*/

// Async Send mail message
func SendAsync(msg Message) {
	// TODO may be need pools limit concurrent nums
	go func() {
		if num, err := Send(msg); err != nil {
			tos := strings.Join(msg.To, "; ")
			info := ""
			if len(msg.Info) > 0 {
				info = ", info: " + msg.Info
			}
			// log failed
			beego.Error(fmt.Sprintf("Async send email %d succeed, not send emails: %s%s err: %s", num, tos, info, err))
		}
	}()
}

// Create html mail message
func NewHtmlMessage(To []string, From, Subject, Body string) Message {
	return Message{
		To:      To,
		From:    From,
		Subject: Subject,
		Body:    Body,
		Type:    "html",
	}
}
