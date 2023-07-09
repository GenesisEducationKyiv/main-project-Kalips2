package presenter

import (
	"btc-app/pkg/domain/model"
	"fmt"
	"net/http"
)

func PresentRate(w http.ResponseWriter, rate *model.CurrencyRate) {
	fmt.Fprint(w, rate.ToString())
}

func PresentString(w http.ResponseWriter, response string) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, response)
}

func PresentErrorByBadRequest(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusBadRequest)
}

func PresentErrorByConflict(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusConflict)
}

func PresentErrorByInternalServerError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
