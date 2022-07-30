package main

import (
	"encoding/json"
	"fmt"
	"time"

	"argonaut/client/email"
	"argonaut/client/rabbitmq"
)

func main() {
	emailProducer := rabbitmq.NewProducer("emails")
	defer emailProducer.Close()

	emailProducer.Start()

	for {
		time.Sleep(time.Second * 10)
		email := &email.Email{
			FromName:    "shubham",
			FromEmail:   "code.shubhamv@gmail.com",
			ToName:      "shubham",
			ToEmail:     "code.shubhamv@gmail.com",
			Subject:     "Test Email",
			TextContent: "Email Body",
		}

		body, err := json.Marshal(email)
		if err != nil {
			fmt.Sprintf("Error marshalling email: %s\n", err.Error())
		}
		emailProducer.Send(body)
	}
}
