package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"taxApi/internal/handlers"
)

func NewRouter(h *handlers.Handler) *mux.Router {
	router := mux.NewRouter()
	unAuth := router.PathPrefix("/registration").Subrouter()
	unAuth.HandleFunc("", h.Registration).Methods(http.MethodPost)
	unAuth.HandleFunc("/token", h.GetTokenToUser).Methods(http.MethodGet)

	AuthCustomer := router.PathPrefix("/customer").Subrouter()
	AuthCustomer.Use(h.Authentication)
	AuthCustomer.HandleFunc("/travel", h.CreateTravel).Methods(http.MethodPost)
	AuthCustomer.HandleFunc("/order", h.GetOrder).Methods(http.MethodGet)

	AuthDriver := router.PathPrefix("/driver").Subrouter()
	AuthDriver.Use(h.Authentication)
	AuthDriver.HandleFunc("/report", h.ReportDriver).Methods(http.MethodPost)
	AuthDriver.HandleFunc("/travel", h.GetTravelList).Methods(http.MethodGet)
	AuthDriver.HandleFunc("/travel/id", h.GetTravelById).Queries("id", "{id}").Methods(http.MethodGet)
	AuthDriver.HandleFunc("/travel/end/id", h.EndTravelById).Queries("id", "{id}").Methods(http.MethodPost)

	AdminAuth := router.PathPrefix("/reports").Subrouter()
	AdminAuth.Use(h.Authentication)
	AdminAuth.HandleFunc("/id", h.ReadById).Queries("id", "{id}").Methods(http.MethodGet)
	AdminAuth.HandleFunc("", h.Report).Methods(http.MethodPost)

	return router
}
