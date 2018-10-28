package payment

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func GetHandlers(s Service) *mux.Router {
	h := handlers{s: s}
	r := mux.NewRouter()
	r.HandleFunc("/payment", h.SavePaymentHandler).Methods("POST")
	return r
}

type handlers struct {
	s Service
}

func (h *handlers) SavePaymentHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	var p Payment

	if err := decoder.Decode(&p); err != nil {
		log.Warn(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := h.s.Save(p)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("/%s", id))
	w.WriteHeader(http.StatusCreated)
	return
}
