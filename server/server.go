package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/attwad/cdf-fe/server/handlers"
	"github.com/attwad/cdf-fe/server/handlers/donation"
	"github.com/attwad/cdf-fe/server/hsts"
	"github.com/attwad/cdf-fe/server/search"
	"github.com/attwad/cdf/db"
	"github.com/attwad/cdf/health"
	"github.com/attwad/cdf/money"
	"github.com/attwad/cdf/stats/io"
	gh "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var (
	hostPort             = flag.String("listen_addr", "127.0.0.1:8080", "Address to listen on.")
	projectID            = flag.String("project_id", "college-de-france", "Google cloud project.")
	elasticAddress       = flag.String("elastic_address", "http://localhost:9002", "HTTP address to elastic instance")
	stripePublishableKey = flag.String("stripe_publishable_key", "pk_test_VpXluQcCGGQVNgg0j0abLR5m", "Publishable key to use with stripe API")
)

func main() {
	flag.Parse()
	ctx := context.Background()
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	dbWrapper, err := db.NewDatastoreWrapper(ctx, *projectID)
	if err != nil {
		log.Fatalf("creating db wrapper: %v", err)
	}
	sr, err := io.NewDatastoreReader(ctx, *projectID)
	if err != nil {
		log.Fatal(err)
	}
	broker, err := money.NewDatastoreBroker(ctx, *projectID)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Will connect to elastic instance @", *elasticAddress)
	r := mux.NewRouter()
	apiRouter := r.PathPrefix("/api").Subrouter()
	apiRouter.Handle("/lessons", handlers.NewLessonsHandler(dbWrapper)).Methods("GET")
	apiRouter.Handle("/search", handlers.NewSearchHandler(search.NewElasticSearcher(*elasticAddress))).Methods("GET")
	apiRouter.Handle("/stats", handlers.NewStatsHandler(sr)).Methods("GET")
	apiRouter.Handle("/donate", donation.NewStripeHandler(os.Getenv("STRIPE_SECRET_KEY"), *stripePublishableKey, broker))
	r.Handle("/healthz", health.NewElasticHealthChecker(*elasticAddress)).Methods("GET")
	appHandler := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "dist/index.html")
	}
	for _, route := range []string{"/searc{*}", "/lesson{*}", "/about", "/donate", "/"} {
		r.Handle(route, http.HandlerFunc(appHandler)).Methods("GET")
	}
	r.Handle(
		"/{[a-z0-9.]+.(js|html|css)}",
		http.FileServer(http.Dir("dist"))).Methods("GET")

	log.Println("Serving on", *hostPort)
	srv := &http.Server{
		Handler: gh.RecoveryHandler(gh.PrintRecoveryStack(true))(gh.CombinedLoggingHandler(os.Stderr, gh.CompressHandler(hsts.NewHandler(r)))),
		Addr:    *hostPort,
		// Good practice: enforce timeouts.
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
