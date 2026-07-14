package emailSender

import (
	"bytes"
	"crypto/tls"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"

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

func resolveTemplatePath(path string) string {
	cleanPath := strings.TrimPrefix(path, "/")
	candidates := []string{
		cleanPath,
		path,
		"./" + cleanPath,
	}

	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)
		candidates = append(candidates, filepath.Join(exeDir, cleanPath))
	}

	for i := 1; i <= 4; i++ {
		candidates = append(candidates, strings.Repeat("../", i)+cleanPath)
	}

	candidates = append(candidates, filepath.Join("c:/workspace/go/wholesomeEdBackEnd", cleanPath))

	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c
		}
	}
	return cleanPath
}

func (i Info) SendEmail() error {
	templateUrl := ""
	i.BaseUrl = "https://test.wholesome.noregrets.ug"
	if i.Type == "registration" {
		templateUrl = "fileRepo/templates/registrationTemplate.html"
	} else if i.Type == "scheduled-interview" {
		templateUrl = "fileRepo/templates/notifyInterview.html"
	} else if i.Type == "application-submitted" {
		templateUrl = "fileRepo/templates/applicationSubmittedTemplate.html"
	} else if i.Type == "staff-deployed" {
		templateUrl = "fileRepo/templates/staffDeployedTemplate.html"
	} else {
		templateUrl = "fileRepo/templates/resetPasswordTemplate.html"
	}

	resolvedPath := resolveTemplatePath(templateUrl)
	t, err := template.ParseFiles(resolvedPath)
	if err != nil {
		log.Println("Error parsing email template:", err)
		return err
	}
	buff := new(bytes.Buffer)
	if err := t.Execute(buff, i); err != nil {
		log.Println("Error executing email template:", err)
		return err
	}

	result := buff.String()

	m := gomail.NewMessage()
	m.SetHeader("From", "andysteve2222@gmail.com")
	m.SetHeader("To", i.Email)
	m.SetHeader("Subject", i.MailType)
	m.SetBody("text/html", result)

	d := gomail.NewDialer("smtp.gmail.com", 587, "andysteve2222@gmail.com", "rhknebyyuqbqfity")
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		log.Println("Error sending email:", err)
		return err
	}

	return nil
}

