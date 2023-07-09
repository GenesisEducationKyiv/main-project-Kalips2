package domain

type EmailMessage struct {
	Header string
	Body   string
}

func NewRateMessage(rate *CurrencyRate, currFrom string, currTo string) *EmailMessage {
	emailSubject := "Поточний курс " + currFrom + " до " + currTo + "."
	return &EmailMessage{Header: emailSubject, Body: rate.ToString()}
}
