package email

type EmailSender interface {
	Send(email *Email) error
}

type EmailSenderImpl struct {
	sendGridClient SendGridClient
}

func (sender EmailSenderImpl) Send(email *Email) error {
	return sender.sendGridClient.Send(email)
}

func NewEmailSender() EmailSender {
	return &EmailSenderImpl{
		sendGridClient: NewSendGridClient(),
	}
}
