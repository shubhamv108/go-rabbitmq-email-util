package email

import (
	"fmt"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func NewSendGridClient() SendGridClient {
	return &SendGridClientImpl{
		apiKey: os.Getenv("SENDGRID_API_KEY"),
	}
}

type SendGridClient interface {
	Send(email *Email) error
}

type SendGridClientImpl struct {
	apiKey string
}

func (sendGridClient *SendGridClientImpl) Send(email *Email) error {
	from := mail.NewEmail(email.FromName, email.FromEmail)
	to := mail.NewEmail(email.ToEmail, email.ToEmail)
	message := mail.NewSingleEmail(from, email.Subject, to, email.TextContent, email.HTMLContent)
	client := sendgrid.NewSendClient(sendGridClient.apiKey)
	response, err := client.Send(message)
	if err != nil {
		return err
	} else {
		fmt.Println(response)
	}
	return nil
}
