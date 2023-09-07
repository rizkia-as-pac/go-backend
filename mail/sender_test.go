package mail

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tech_school/simple_bank/utils/conf"
)

func TestSendEmailWithGmail(t *testing.T) {
	if testing.Short() {
		t.Skip() // send email test wont executed if short flag is set 
	}

	config, err := conf.LoadConfig("..")
	require.NoError(t, err)

	GmailSender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword) // create new gmail sender

	subject := "A test email"
	// in go `` backtick can be use to write multiline string
	content := `
	<h1>Hello world</h1>
	<p>This is a test message from <a href="http://techschool.guru">Tech School</a></p>
	`
	to := []string{"rizkia.as.pac@gmail.com"}
	attachFiles := []string{"../app.env"}

	err = GmailSender.SendEmail(subject, content, to, nil, nil, attachFiles)
	require.NoError(t, err)
}
