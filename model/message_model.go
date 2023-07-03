package model

import "strconv"

type Message struct {
	EmailFrom string
	Header    string
	Body      string
}

func NewRateMessage(rate float64, emailFrom string, currFrom string, currTo string) *Message {
	rateFormatted := strconv.FormatFloat(rate, 'f', 5, 64)
	emailSubject := "Поточний курс " + currFrom + " до " + currTo + "."
	return &Message{EmailFrom: emailFrom, Header: emailSubject, Body: rateFormatted}
}
