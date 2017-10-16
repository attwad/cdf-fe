package donation

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	stripe "github.com/stripe/stripe-go"
)

const pubKey = "a publishable key"

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

type fakeBroker struct {
	balance      int
	balanceError error
}

func (b *fakeBroker) GetBalance(ctx context.Context) (int, error) {
	return b.balance, b.balanceError
}

func (b *fakeBroker) ChangeBalance(ctx context.Context, delta int) error {
	if b.balanceError != nil {
		return b.balanceError
	}
	b.balance += delta
	return nil
}

type fakeReporter struct {
	reported error
}

func (f *fakeReporter) Report(err error) {
	f.reported = err
}

func (f *fakeReporter) Close() error {
	return nil
}

func TestHTTPHandlerPOST(t *testing.T) {
	tests := []struct {
		msg         string
		req         postRequest
		wantStatus  int
		api         fakeStripeAPIWrapper
		requestType string
		broker      fakeBroker
		wantBalance int
		wantErr     string
	}{
		{
			msg:         "charge succeeds",
			req:         postRequest{StripeEmail: "foo@bar.com", StripeToken: "tok", AmountUsdCents: minPaymentsUsdCents * 2},
			requestType: "POST",
			wantStatus:  200,
			wantBalance: minPaymentsUsdCents * 2,
		}, {
			msg:         "too small amount",
			req:         postRequest{StripeEmail: "foo@bar.com", StripeToken: "tok", AmountUsdCents: minPaymentsUsdCents - 1},
			requestType: "POST",
			wantStatus:  400,
		}, {
			msg:         "fail on new charge",
			req:         postRequest{StripeEmail: "foo@bar.com", StripeToken: "tok", AmountUsdCents: minPaymentsUsdCents * 2},
			api:         fakeStripeAPIWrapper{errCharge: errors.New("invalid charge")},
			requestType: "POST",
			wantStatus:  500,
			wantErr:     "invalid charge",
		}, {
			msg:         "fail on new customer",
			req:         postRequest{StripeEmail: "foo@bar.com", StripeToken: "tok", AmountUsdCents: minPaymentsUsdCents * 2},
			api:         fakeStripeAPIWrapper{errNewCustomer: errors.New("invalid customer")},
			requestType: "POST",
			wantStatus:  500,
			wantErr:     "invalid customer",
		}, {
			msg:         "wrong request type",
			req:         postRequest{StripeEmail: "foo@bar.com", StripeToken: "tok", AmountUsdCents: minPaymentsUsdCents * 2},
			requestType: "TRACE",
			wantStatus:  405,
		}, {
			msg:         "change balance fails",
			req:         postRequest{StripeEmail: "foo@bar.com", StripeToken: "tok", AmountUsdCents: minPaymentsUsdCents * 2},
			requestType: "POST",
			broker:      fakeBroker{balanceError: errors.New("wrong balance")},
			wantStatus:  500,
			wantErr:     "wrong balance",
		},
	}
	for _, test := range tests {
		fr := &fakeReporter{reported: errors.New("")}
		h := &donateHandler{&test.api, &test.broker, pubKey, fr}
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

		if got, want := test.broker.balance, test.wantBalance; got != want {
			t.Errorf("[%s] balance got=%d, want=%d", test.msg, got, want)
		}
		if got, want := fr.reported.Error(), test.wantErr; !strings.Contains(got, want) {
			t.Errorf("[%s] error not contained, got=%q, want=%q", test.msg, got, want)
		}
	}
}

func TestHTTPHandlerGET(t *testing.T) {
	fr := &fakeReporter{}
	h := &donateHandler{&fakeStripeAPIWrapper{}, nil, pubKey, fr}
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
		t.Errorf("one hour amountUsdCents got=%d, want=%d", got, want)
	}

	if got, want := gr.StripePublishableKey, pubKey; got != want {
		t.Errorf("stripe publishable key preparation, got=%q, want=%q", got, want)
	}
}
