package health

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"
)

type elasticHealthCheck struct {
	client         *http.Client
	elasticAddress string
}

// NewElasticHealthChecker returns a new HTTP handler that returns a status 200
// if elastic search is healthy.
func NewElasticHealthChecker(elasticAddress string) http.Handler {
	return &elasticHealthCheck{
		client: &http.Client{
			Timeout: time.Second * 2,
		},
		elasticAddress: elasticAddress,
	}
}

func (h *elasticHealthCheck) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(h.elasticAddress)
	if err != nil {
		http.Error(w, "Invalid elasticsearch address", http.StatusPreconditionFailed)
		return
	}
	u.Path = "_cluster/health"
	resp, err := h.client.Get(u.String())
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	type healthResponse struct {
		Status string `json:"status"`
	}
	var hr healthResponse
	if err := json.NewDecoder(resp.Body).Decode(&hr); err != nil {
		log.Println("Elasticsearch is unhealthy:", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	if hr.Status != "green" && hr.Status != "yellow" {
		log.Printf("Elasticsearch is unhealthy:\n%+v\n", hr)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}
