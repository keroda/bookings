package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/keroda/bookings/internal/models"
	mail "github.com/xhit/go-simple-mail/v2"
)

func listenForMail() {
	go func() {
		//infinite loop
		for {
			msg := <-app.MailChan
			sendMsg(msg)
		}
	}()
}

func sendMsg(m models.MailData) {
	mailServer := mail.NewSMTPClient()
	mailServer.Host = "localhost"
	mailServer.Port = 1025
	mailServer.KeepAlive = false
	mailServer.ConnectTimeout = 10 * time.Second
	mailServer.SendTimeout = 10 * time.Second
	//usernamne
	//password
	//encryption

	client, err := mailServer.Connect()
	if err != nil {
		errorLog.Println(err)
	}

	email := mail.NewMSG()
	email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)

	var msgToSend string
	if m.Template == "" {
		msgToSend = m.Content
	} else {
		//load template from file
		data, err := ioutil.ReadFile(fmt.Sprintf("/email-templates/%s", m.Template))
		if err != nil {
			app.ErrorLog.Println(err)
		}
		mailTemplate := string(data)
		msgToSend = strings.Replace(mailTemplate, "[%body%]", m.Content, 1)
	}
	email.SetBody(mail.TextHTML, msgToSend)

	err = email.Send(client)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Email sent")
	}
}
