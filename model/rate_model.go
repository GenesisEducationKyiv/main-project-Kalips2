package model

import "strconv"

type Rate struct {
	Value float64
}

func (rate *Rate) String() string {
	return strconv.FormatFloat(rate.Value, 'f', 5, 64)
}
