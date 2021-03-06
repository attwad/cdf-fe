package donation

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/attwad/cdf/errorreport"
	"github.com/attwad/cdf/money"
	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/customer"
)

type donateHandler struct {
	wrapper              stripeWrapper
	broker               money.Broker
	stripePublishableKey string
	reporter             errorreport.Reporter
}

const minPaymentsUsdCents = 100

// NewStripeHandler handles HTTP requests for donations via Stripe.
func NewStripeHandler(privateKey, publishableKey string, broker money.Broker, reporter errorreport.Reporter) http.Handler {
	stripe.Key = privateKey
	return &donateHandler{
		wrapper:              &stripeAPIWrapper{},
		broker:               broker,
		stripePublishableKey: publishableKey,
		reporter:             reporter,
	}
}

type stripeWrapper interface {
	NewCharge(*stripe.ChargeParams) (*stripe.Charge, error)
	NewCustomer(*stripe.CustomerParams) (*stripe.Customer, error)
}

type stripeAPIWrapper struct {
}

func (s *stripeAPIWrapper) NewCharge(p *stripe.ChargeParams) (*stripe.Charge, error) {
	return charge.New(p)
}

func (s *stripeAPIWrapper) NewCustomer(p *stripe.CustomerParams) (*stripe.Customer, error) {
	return customer.New(p)
}

type postRequest struct {
	StripeToken    string `json:"stripeToken"`
	StripeEmail    string `json:"stripeEmail"`
	AmountUsdCents uint64 `json:"amountUsdCents"`
}

type postResponse struct {
}

type getResponse struct {
	OneHourAmountUsdCents int    `json:"one_hour_amount_usd_cents"`
	StripePublishableKey  string `json:"stripe_publishable_key"`
}

// logAndReport is used when an email needs to be sent, the error logs in
// stackdriver and a dev having a look REAL QUICK (aka. something in the
// donation flow didn't work)
func (h *donateHandler) logAndReport(w http.ResponseWriter, err error) {
	log.Println(err)
	http.Error(w, err.Error(), http.StatusInternalServerError)
	h.reporter.Report(err)
}

func (h *donateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		gr := getResponse{
			OneHourAmountUsdCents: money.DurationToUsdCents(1 * time.Hour),
			StripePublishableKey:  h.stripePublishableKey,
		}
		enc := json.NewEncoder(w)
		if err := enc.Encode(&gr); err != nil {
			log.Println("Could not write json output:", err)
			http.Error(w, "Could not write json", http.StatusInternalServerError)
			return
		}
		return
	}
	// TODO: Stripe does its own kind of xsrf protection by not allowing reuse of
	// a token, we should still do it ourselves to avoid a first token being used.
	if r.Method != "POST" {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}
	var req postRequest
	d := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := d.Decode(&req); err != nil {
		log.Println("Could not decode post request:", err)
		http.Error(w, "Could not decode post request", http.StatusInternalServerError)
		return
	}
	log.Println("Request parsed:", req)

	if req.AmountUsdCents < minPaymentsUsdCents {
		log.Println("amount must be >=", minPaymentsUsdCents)
		http.Error(w, fmt.Sprintf("amount must be >= %d", minPaymentsUsdCents), http.StatusBadRequest)
		return
	}

	customerParams := &stripe.CustomerParams{Email: req.StripeEmail}
	customerParams.SetSource(req.StripeToken)

	newCustomer, err := h.wrapper.NewCustomer(customerParams)
	if err != nil {
		h.logAndReport(w, fmt.Errorf("Error in new customer: %v", err))
		return
	}
	log.Println("Customer created")

	chargeParams := &stripe.ChargeParams{
		Amount:   req.AmountUsdCents,
		Currency: "usd",
		Desc:     "college-audio.science audio transcriptions",
		Customer: newCustomer.ID,
	}

	log.Println("Creating charge")
	if _, err := h.wrapper.NewCharge(chargeParams); err != nil {
		h.logAndReport(w, fmt.Errorf("Error in new charge: %v", err))
		return
	}

	log.Println("Charge successful")

	// Hopefully nobody is crazy enough to give me more than max_int usd cents...
	if err := h.broker.ChangeBalance(r.Context(), int(req.AmountUsdCents)); err != nil {
		h.logAndReport(w, fmt.Errorf("Error increasing balance of account: %v", err))
		return
	}

	log.Println("Balance increased by", req.AmountUsdCents)

	enc := json.NewEncoder(w)
	if err := enc.Encode(&postResponse{}); err != nil {
		h.logAndReport(w, fmt.Errorf("Error increasing balance of account: %v", err))
		return
	}
}
