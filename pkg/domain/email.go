package domain

type Email struct {
	Address string
}

func NewEmail(emailAddress string) Email {
	return Email{emailAddress}
}

func (e *Email) GetAddress() string {
	return e.Address
}
