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

	r.HandleFunc("/payment", h.savePaymentHandler).
		Methods("POST")

	r.HandleFunc("/payment/search", h.searchForPayments).
		Methods("GET")

	r.HandleFunc("/payment/{id:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}", h.getPaymentHandler).
		Methods("GET")

	r.HandleFunc("/__health", h.healthCheckHandler).
		Methods("GET")

	return r
}

type handlers struct {
	s Service
}

func (h *handlers) getPaymentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	w.Header().Set("Content-Type", "application/json")
	p, err := h.s.Get(r.Context(), id)
	if err != nil {
		switch err {
		case ErrNotFound:
			w.WriteHeader(http.StatusNotFound)
			return
		default:
			log.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	if err := json.NewEncoder(w).Encode(p); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *handlers) searchForPayments(w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query()
	id := vars["organisation_id"]
	w.Header().Set("Content-Type", "application/json")
	ps, err := h.s.SearchByOrganisationId(r.Context(), id[0])
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(Payments{Payments: ps}); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *handlers) savePaymentHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	var p Payment

	if err := decoder.Decode(&p); err != nil {
		log.Warn(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := h.s.Save(r.Context(), p)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("/payment/%s", id))
	w.WriteHeader(http.StatusCreated)
	return
}

func (h *handlers) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	hc := h.s.HealthCheck(r.Context())

	if !hc.Healthy {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	if err := json.NewEncoder(w).Encode(hc); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
