package mail

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

const (
	smptAuthAddress   = "smtp.gmail.com" // host
	smtpServerAddress = "smtp.gmail.com:587"
)

type EmailSender interface {
	SendEmail(
		subject string,
		content string,
		to []string, // list email address to send email to
		cc []string,
		bcc []string,
		attachFiles []string, // attach some files to the email. list of attach files name
	) error // return error if it fails to send an email
}

type GmailSender struct {
	name              string // recipient will see this as the sender of the email
	fromEmailAddress  string
	fromEmailPassword string // password to access that sender account // kali ini kita akan menggunakan dummy password // pada kasus sesungguhnya kita akan menggunakan app password yang digenerate
}

func NewGmailSender(name string, fromEmailAddress string, fromEmailPassword string) EmailSender {
	return &GmailSender{
		name:              name,
		fromEmailAddress:  fromEmailAddress,
		fromEmailPassword: fromEmailPassword,
	}
}

func (sender *GmailSender) SendEmail(
	subject string,
	content string,
	to []string,
	cc []string,
	bcc []string,
	attachFiles []string,
) error {
	e := email.NewEmail() // create new email object

	// e.from kombinasi dari nama email pengirim dan alamatnya
	e.From = fmt.Sprintf("%s <%s>", sender.name, sender.fromEmailAddress)

	e.Subject = subject

	// konten dari email
	// kita harus melakukan conversi terlebih dulu, karna content merupakan string sementara HTML field adalah []byte slice
	e.HTML = []byte(content)

	e.To = to
	e.Cc = cc
	e.Bcc = bcc

	for _, file := range attachFiles {
		// for each file we call e.attachfile
		_, err := e.AttachFile(file)
		if err != nil {
			return fmt.Errorf("gagal to attach file %s: %w ", file, err)
		}
	}

	// authenticating with SMTP server
	// ussualy first arg can be left empty,
	smtpAuth := smtp.PlainAuth("", sender.fromEmailAddress, sender.fromEmailPassword, smptAuthAddress)

	// send the email
	return e.Send(smtpServerAddress, smtpAuth)
}
