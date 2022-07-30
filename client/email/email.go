package email

type Email struct {
	FromName    string
	FromEmail   string
	ToName      string
	ToEmail     string
	Subject     string
	TextContent string
	HTMLContent string
}
