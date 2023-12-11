package emailSender

import (
	"bytes"
	"html/template"
	"log"

	"gopkg.in/gomail.v2"
)

type Info struct {
	Email       string
	Token       string
	MailType    string
	UserName    string
	CurrentYear int
	Type        string
}

func (i Info) SendEmail() {
	templateUrl := ""
	if i.Type == "registration" {
		templateUrl = "./fileRepo/templates/registrationTemplate.html"
	} else {
		templateUrl = "./fileRepo/templates/resetPasswordTemplate.html"
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

	m := gomail.NewMessage()
	m.SetHeader("From", "")
	m.SetHeader("To", i.Email)
	m.SetHeader("Subject", i.MailType)
	//m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	//m.SetBody("text/html", "Hello Andy steve please find this email! "+i.Token)
	m.SetBody("text/html", result)
	//m.Attach("/home/Alex/lolcat.jpg")

	d := gomail.NewDialer("smtp.gmail.com", 587, "", "")

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		//panic(err)
	}
}
