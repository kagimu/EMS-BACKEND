package emailSender

import (
	"bytes"
	"crypto/tls"
	"html/template"
	"log"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

type Info struct {
	Email       string
	Token       string
	MailType    string
	UserName    string
	CurrentYear int
	Type        string
	Data        interface{}
	BaseUrl     string
}

func (i Info) SendEmail() error {
	templateUrl := ""
	i.BaseUrl = "https://test.wholesome.noregrets.ug"
	if i.Type == "registration" {
		templateUrl = "/fileRepo/templates/registrationTemplate.html"
	} else if i.Type == "scheduled-interview" {
		templateUrl = "/fileRepo/templates/notifyInterview.html"
	} else if i.Type == "application-submitted" {
		templateUrl = "/fileRepo/templates/applicationSubmittedTemplate.html"
	} else {
		templateUrl = "/fileRepo/templates/resetPasswordTemplate.html"
	}

	t, err := template.ParseFiles(templateUrl)

	if err != nil {
		log.Println(err)
		//log.Println("Error one")
	}
	buff := new(bytes.Buffer)
	//var tpl bytes.Buffer
	if err := t.Execute(buff, i); err != nil {
		log.Println(err)
		//log.Println("Error two")
	}

	result := buff.String()

	// EMS Release 1 (security fix B2): SMTP credentials are read from the
	// environment (MAIL_FROM / MAIL_HOST / MAIL_PORT / MAIL_USERNAME /
	// MAIL_PASSWORD) instead of being hardcoded in source. No other email
	// behaviour is changed.
	from := os.Getenv("MAIL_FROM")
	host := os.Getenv("MAIL_HOST")
	username := os.Getenv("MAIL_USERNAME")
	password := os.Getenv("MAIL_PASSWORD")

	port, err := strconv.Atoi(os.Getenv("MAIL_PORT"))
	if err != nil {
		port = 587
	}

	m := gomail.NewMessage()
	if from != "" {
		m.SetHeader("From", from)
	} else {
		m.SetHeader("From", username)
	}
	m.SetHeader("To", i.Email)
	m.SetHeader("Subject", i.MailType)
	m.SetBody("text/html", result)

	d := gomail.NewDialer(host, port, username, password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		log.Println("Error sending email:", err)
	}

	return nil
}
