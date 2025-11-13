package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/omerkoc09/go-stripe/internal/cards"
)

type stripePayload struct {
	Currency string `json:"currency"`
	Amount   string `json:"amount"`
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message,omitempty"`
	Content string `json:"content,omitempty"`
	ID      int    `json:"id,omitempty"`
}

func (app *application) GetPaymentIntent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Read body for debugging
	bodyBytes, _ := io.ReadAll(r.Body)
	app.infoLog.Printf("Received request body: %s", string(bodyBytes))

	// Reset body for decoder
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var payload stripePayload

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		app.errorLog.Printf("JSON decode error: %v, body: %s", err, string(bodyBytes))
		j := jsonResponse{
			OK:      false,
			Message: fmt.Sprintf("Invalid request payload: %v", err),
		}
		out, _ := json.Marshal(j)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(out)
		return
	}

	app.infoLog.Printf("Parsed payload: Amount=%s, Currency=%s", payload.Amount, payload.Currency)

	amount, err := strconv.Atoi(payload.Amount)
	if err != nil {
		app.errorLog.Println(err)
		j := jsonResponse{
			OK:      false,
			Message: "Invalid amount",
		}
		out, _ := json.Marshal(j)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(out)
		return
	}

	// Check if Stripe secret key is configured
	if app.config.stripe.secret == "" {
		app.errorLog.Println("Stripe secret key is not configured")
		j := jsonResponse{
			OK:      false,
			Message: "Stripe secret key is not configured. Please set STRIPE_SECRET environment variable.",
		}
		out, _ := json.Marshal(j)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(out)
		return
	}

	card := cards.Card{
		Secret:   app.config.stripe.secret,
		Key:      app.config.stripe.key,
		Currency: payload.Currency,
	}

	pi, msg, err := card.CreatePaymentIntent(payload.Currency, amount)
	if err != nil {
		app.errorLog.Println(err)
		j := jsonResponse{
			OK:      false,
			Message: msg,
			Content: "",
		}
		out, err := json.Marshal(j)
		if err != nil {
			app.errorLog.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(out)
		return
	}

	// Success - return payment intent
	out, err := json.MarshalIndent(pi, "", "   ")
	if err != nil {
		app.errorLog.Println(err)
		j := jsonResponse{
			OK:      false,
			Message: "Error processing payment intent",
		}
		out, _ := json.Marshal(j)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(out)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

func (app *application) GetMacWithId(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	macId, _ := strconv.Atoi(id)

	mac, err := app.DB.GetMac(macId)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	out, err := json.MarshalIndent(mac, "", "   ")
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)

}
