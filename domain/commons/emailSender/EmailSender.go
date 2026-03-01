package emailSender

import (
	"bytes"
	"crypto/tls"
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

	m := gomail.NewMessage()
	m.SetHeader("From", "andysteve2222@gmail.com")
	m.SetHeader("To", i.Email)
	m.SetHeader("Subject", i.MailType)
	//m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	//m.SetBody("text/html", "Hello Andy steve please find this email! "+i.Token)
	m.SetBody("text/html", result)
	//m.Attach("/home/Alex/lolcat.jpg")

	d := gomail.NewDialer("smtp.gmail.com", 587, "andysteve2222@gmail.com", "rhknebyyuqbqfity")
	//d := gomail.NewDialer(os.Getenv("MAIL_HOST"), os.Getenv("MAIL_PORT"), os.Getenv("MAIL_USERNAME"), os.Getenv("MAIL_PASSWORD"))
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		log.Println("Error sending email:", err)
	}

	return nil
}
