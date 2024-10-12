package kafka

type EmailSendindRequestedPayload struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Content string `json:"content"`
}
