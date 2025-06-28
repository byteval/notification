package notification

type EmailMessage struct {
	UID         uint32
	Subject     string
	From        string
	To          []string
	Date        string
	TextBody    string
	HTMLBody    string
	Attachments []EmailAttachment
}

type EmailAttachment struct {
	Name        string
	ContentType string
	Content     []byte
}
