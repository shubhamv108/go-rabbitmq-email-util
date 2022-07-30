package main

import (
	"encoding/json"

	email2 "argonaut/client/email"
	"argonaut/client/rabbitmq"
)

func main() {
	emailMessageProcessor := NewEmailMessageProcessor()
	emailConsumer := rabbitmq.NewConsumer("emails", emailMessageProcessor.messageProcessor)
	defer emailConsumer.Close()
	emailConsumer.Start()
}

func NewEmailMessageProcessor() EmailMessageProcessor {
	return &EmailMessageProcessorImpl{
		emailSender: email2.NewEmailSender(),
	}
}

type EmailMessageProcessor interface {
	messageProcessor(message []byte) error
}

type EmailMessageProcessorImpl struct {
	emailSender email2.EmailSender
}

func (emailConsumer *EmailMessageProcessorImpl) messageProcessor(message []byte) error {
	var email *email2.Email
	if err := json.Unmarshal(message, &email); err != nil {
		return err
	}
	return emailConsumer.emailSender.Send(email)
}
