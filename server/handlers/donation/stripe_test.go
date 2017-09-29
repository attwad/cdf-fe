package donation

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	stripe "github.com/stripe/stripe-go"
)

type fakeStripeAPIWrapper struct {
	errCharge      error
	errNewCustomer error
}

func (s *fakeStripeAPIWrapper) NewCharge(p *stripe.ChargeParams) (*stripe.Charge, error) {
	return &stripe.Charge{}, s.errCharge
}

func (s *fakeStripeAPIWrapper) NewCustomer(p *stripe.CustomerParams) (*stripe.Customer, error) {
	return &stripe.Customer{}, s.errNewCustomer
}

func TestHTTPHandlerPOST(t *testing.T) {
	tests := []struct {
		msg         string
		req         postRequest
		wantStatus  int
		api         fakeStripeAPIWrapper
		requestType string
	}{
		{
			msg:         "charge succeeds",
			req:         postRequest{StripeEmail: "foo@bar.com", StripeToken: "tok", Amount: minPaymentsUsdCents * 2},
			requestType: "POST",
			wantStatus:  200,
		}, {
			msg:         "too small amount",
			req:         postRequest{StripeEmail: "foo@bar.com", StripeToken: "tok", Amount: minPaymentsUsdCents - 1},
			requestType: "POST",
			wantStatus:  400,
		}, {
			msg:         "fail on new charge",
			req:         postRequest{StripeEmail: "foo@bar.com", StripeToken: "tok", Amount: minPaymentsUsdCents * 2},
			api:         fakeStripeAPIWrapper{errCharge: errors.New("invalid charge")},
			requestType: "POST",
			wantStatus:  500,
		}, {
			msg:         "fail on new customer",
			req:         postRequest{StripeEmail: "foo@bar.com", StripeToken: "tok", Amount: minPaymentsUsdCents * 2},
			api:         fakeStripeAPIWrapper{errNewCustomer: errors.New("invalid customer")},
			requestType: "POST",
			wantStatus:  500,
		}, {
			msg:         "wrong request type",
			req:         postRequest{StripeEmail: "foo@bar.com", StripeToken: "tok", Amount: minPaymentsUsdCents * 2},
			requestType: "TRACE",
			wantStatus:  405,
		},
	}
	for _, test := range tests {
		h := &donateHandler{&test.api}
		b, err := json.Marshal(&test.req)
		if err != nil {
			t.Errorf("[%s]: marshalling request:%s", test.msg, err)
			continue
		}
		req := httptest.NewRequest(test.requestType, "/", bytes.NewBuffer(b))
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		resp := w.Result()
		if got, want := resp.StatusCode, test.wantStatus; got != want {
			t.Errorf("[%s] status got=%d, want=%d", test.msg, got, want)
		}
	}
}

func TestHTTPHandlerGET(t *testing.T) {
	h := &donateHandler{&fakeStripeAPIWrapper{}}
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	resp := w.Result()
	if got, want := resp.StatusCode, 200; got != want {
		t.Errorf("status got=%d, want=%d", got, want)
	}
	gr := &getResponse{}
	d := json.NewDecoder(resp.Body)
	defer resp.Body.Close()
	if err := d.Decode(gr); err != nil {
		t.Fatalf("unmarshalling response:%s", err)
	}
	if got, want := gr.OneHourAmountUsdCents, 144; got != want {
		t.Errorf("one hour amount got=%d, want=%d", got, want)
	}
}
