package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (app *application) VirtualTerminal(w http.ResponseWriter, r *http.Request) {

	if err := app.renderTemplate(w, r, "terminal", &templateData{}, "stripe-js"); err != nil {
		app.errorLog.Println(err)
	}
}

// PaymentIntentProxy proxies payment intent requests to the API server
func (app *application) PaymentIntentProxy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		app.errorLog.Printf("Error reading request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{
			"ok":      false,
			"message": "Bad Request: Invalid request body",
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	defer r.Body.Close()

	// Create a new request to the API server
	apiURL := fmt.Sprintf("%s/api/payment-intent", app.config.api)
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(body))
	if err != nil {
		app.errorLog.Printf("Error creating request: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]interface{}{
			"ok":      false,
			"message": "Internal Server Error: Failed to create request",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Copy headers
	req.Header.Set("Content-Type", "application/json")

	// Make the request to API server
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		app.errorLog.Printf("Error making request to API: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]interface{}{
			"ok":      false,
			"message": fmt.Sprintf("Internal Server Error: Failed to connect to API server: %v", err),
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	defer resp.Body.Close()

	// Copy response headers (but keep Content-Type as JSON)
	for key, values := range resp.Header {
		if key != "Content-Type" {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
	}

	// Ensure Content-Type is JSON
	w.Header().Set("Content-Type", "application/json")

	// Set status code
	w.WriteHeader(resp.StatusCode)

	// Copy response body
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		app.errorLog.Printf("Error copying response: %v", err)
	}
}

func (app *application) PaymentSucceeded(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	// read posted data
	cardHolder := r.Form.Get("cardholder-name")
	cardHolderEmail := r.Form.Get("cardholder-email")
	paymentIntent := r.Form.Get("payment_intent")
	paymentMethod := r.Form.Get("payment_method")
	paymentAmount := r.Form.Get("payment_amount")
	paymentCurrency := r.Form.Get("payment_currency")

	data := make(map[string]interface{})
	data["cardholder"] = cardHolder
	data["cardholder-email"] = cardHolderEmail
	data["pi"] = paymentIntent
	data["pm"] = paymentMethod
	data["pa"] = paymentAmount
	data["pc"] = paymentCurrency

	if err := app.renderTemplate(w, r, "succeeded", &templateData{
		Data: data,
	}); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) BuyOnce(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	macId, err := strconv.Atoi(id)
	if err != nil {
		app.errorLog.Printf("Invalid ID parameter: %v", err)
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	mac, err := app.DB.GetMac(macId)
	if err != nil {
		app.errorLog.Printf("Error getting Mac with ID %d: %v", macId, err)
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	data := make(map[string]interface{})
	data["mac"] = mac

	if err := app.renderTemplate(w, r, "buy-once", &templateData{
		Data: data,
	}, "stripe-js"); err != nil {
		app.errorLog.Printf("Error rendering template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
