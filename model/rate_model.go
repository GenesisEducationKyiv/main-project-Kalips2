package model

import "strconv"

type Rate struct {
	Value float64
}

func (rate *Rate) ToString() string {
	return strconv.FormatFloat(rate.Value, 'f', 5, 64)
}

func (rate *Rate) SetValue(value string) {
	rate.Value, _ = strconv.ParseFloat(value, 64)
}
