package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Get("/virtual-terminal", app.VirtualTerminal)
	mux.Post("/api/payment-intent", app.PaymentIntentProxy)
	mux.Post("/payment-succeeded", app.PaymentSucceeded)
	mux.Get("/mac/{id}", app.BuyOnce)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
