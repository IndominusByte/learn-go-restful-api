package mail

import (
	"bytes"
	"errors"
	"html/template"
	"log"
	"os"
	"strings"

	"gopkg.in/gomail.v2"
)

type Mail struct {
	Server   string
	Port     int
	Username string
	Password string
}

func (m *Mail) SendEmail(files, emails []string, subject, sender, html string, templateVar interface{}) {
	templ := strings.Split(html, "/")
	t := template.New(templ[len(templ)-1])

	t, err := t.ParseFiles(html)
	if err != nil {
		panic(err)
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, templateVar); err != nil {
		panic(err)
	}

	for _, v := range files {
		// check directory or file exists
		if _, err = os.Stat(v); errors.Is(err, os.ErrNotExist) {
			panic(err)
		}
	}

	d := gomail.NewDialer(m.Server, m.Port, m.Username, m.Password)

	s, err := d.Dial()
	if err != nil {
		panic(err)
	}

	mail, result := gomail.NewMessage(), tpl.String()

	for _, r := range emails {
		mail.SetHeader("From", sender)
		mail.SetHeader("To", r)
		mail.SetHeader("Subject", subject)
		mail.SetBody("text/html", result)
		for _, v := range files {
			mail.Attach(v)
		}

		if err := gomail.Send(s, mail); err != nil {
			log.Printf("Could not send email to %q: %v", r, err)
		}
		mail.Reset()
	}
}
