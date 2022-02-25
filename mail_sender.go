package main

import (
	"log"
	"net/smtp"
	"strconv"
)

var auth smtp.Auth

func init() {
	// authenticate with host
	auth = smtp.PlainAuth("", config.Sender, config.SenderPass, config.SmptHost)
}

func SendMail(msg string) {
	body := []byte(msg)
	// send email to list of recipients
	err := smtp.SendMail(config.SmptHost+":"+strconv.Itoa(config.SmtpPort), auth, config.Sender, config.Recipients, body)
	if err != nil {
		log.Println(err)
	}
	log.Printf("Mail sent to %v\n", config.Recipients)
}
