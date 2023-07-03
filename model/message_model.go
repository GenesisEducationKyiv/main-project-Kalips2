package model

type Message struct {
	Header string
	Body   string
}

func NewRateMessage(rate *Rate, currFrom string, currTo string) *Message {
	emailSubject := "Поточний курс " + currFrom + " до " + currTo + "."
	return &Message{Header: emailSubject, Body: rate.ToString()}
}
