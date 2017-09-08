package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/attwad/cdf-fe/server/hsts"
	"github.com/attwad/cdf-fe/server/search"
	"github.com/attwad/cdf/data"
	"github.com/attwad/cdf/db"
	"github.com/attwad/cdf/health"
	"github.com/attwad/cdf/stats/io"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var (
	hostPort       = flag.String("listen_addr", "127.0.0.1:8080", "Address to listen on.")
	projectID      = flag.String("project_id", "college-de-france", "Google cloud project.")
	elasticAddress = flag.String("elastic_address", "", "HTTP address to elastic instance")
	enableHSTS     = flag.Bool("enable_hsts", true, "Enable HSTS header in HTTP responses")
)

const (
	pageSize = 10
)

type server struct {
	ctx         context.Context
	db          db.Wrapper
	searcher    search.Searcher
	statsReader io.Reader
}

func (s *server) APIServeLessons(w http.ResponseWriter, r *http.Request) {
	filter := db.FilterNone
	if r.URL.Query().Get("filter") == "converted" {
		filter = db.FilterOnlyConverted
	}
	lessons, cursor, err := s.db.GetLessons(s.ctx, r.URL.Query().Get("cursor"), filter, pageSize)
	if err != nil {
		log.Println("Could not read lessons from db:", err)
		http.Error(w, "Could not read lessons from DB", http.StatusInternalServerError)
		return
	}
	type response struct {
		Cursor  string                `json:"cursor"`
		Lessons []data.ExternalCourse `json:"lessons"`
	}
	resp := &response{Cursor: cursor, Lessons: make([]data.ExternalCourse, 0)}
	for _, lesson := range lessons {
		resp.Lessons = append(resp.Lessons, data.ExternalCourse{
			Course:            lesson.Course,
			FormattedDate:     fmt.Sprintf("%d/%d/%d", lesson.Date.Day(), lesson.Date.Month(), lesson.Date.Year()),
			FormattedDuration: fmt.Sprintf("%d min.", lesson.DurationSec/60),
		})
	}
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		log.Println("Could not write json output:", err)
		http.Error(w, "Could not write json", http.StatusInternalServerError)
		return
	}
}

func (s *server) APIServeSearch(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if strings.TrimSpace(q) == "" {
		http.Error(w, "empty query", http.StatusBadRequest)
		return
	}
	fStr := r.URL.Query().Get("from")
	from, err := strconv.Atoi(fStr)
	if err != nil {
		http.Error(w, "from param incorrect, must be a positive number", http.StatusBadRequest)
		return
	}
	sStr := r.URL.Query().Get("size")
	size, err := strconv.Atoi(sStr)
	if err != nil {
		http.Error(w, "from param incorrect, must be a positive number", http.StatusBadRequest)
		return
	}
	jsr, err := s.searcher.Search(q, from, size)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	type searchResponse struct {
		Query    string          `json:"query"`
		TookMs   int             `json:"took_ms"`
		TimedOut bool            `json:"timed_out"`
		Total    int             `json:"total"`
		Sources  []search.Source `json:"sources"`
	}
	sr := searchResponse{
		Query:    q,
		TookMs:   jsr.TookMs,
		TimedOut: jsr.TimedOut,
		Total:    jsr.Hits.Total,
		Sources:  make([]search.Source, 0),
	}
	for _, hit := range jsr.Hits.Hits {
		sr.Sources = append(sr.Sources, hit.Source)
	}
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if err := enc.Encode(sr); err != nil {
		log.Println("Could not write json output:", err)
		http.Error(w, "Could not write json", http.StatusInternalServerError)
		return
	}
}

func (s *server) APIServeStats(w http.ResponseWriter, r *http.Request) {
	stats, err := s.statsReader.Read(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	type uiStats struct {
		Computed             string
		NumTotal             int
		NumConverted         int
		ConvertedDurationSec int
		LeftDurationSec      int
		PercentDone          float32
	}
	uis := &uiStats{
		Computed:             stats.Computed.Format("Jan _2 2006"),
		NumTotal:             stats.NumTotal,
		NumConverted:         stats.NumConverted,
		ConvertedDurationSec: stats.ConvertedDurationSec,
		LeftDurationSec:      stats.LeftDurationSec,
	}
	if stats.LeftDurationSec+stats.ConvertedDurationSec > 0 {
		uis.PercentDone = float32(stats.ConvertedDurationSec) * 100 / float32(stats.LeftDurationSec+stats.ConvertedDurationSec)
	}
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if err := enc.Encode(uis); err != nil {
		log.Println("Could not write json output:", err)
		http.Error(w, "Could not write json", http.StatusInternalServerError)
		return
	}
}

func main() {
	flag.Parse()
	ctx := context.Background()
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	dbWrapper, err := db.NewDatastoreWrapper(ctx, *projectID)
	if err != nil {
		log.Fatalf("creating db wrapper: %v", err)
	}
	log.Println("Will connect to elastic instance @", *elasticAddress)
	sr, err := io.NewDatastoreReader(ctx, *projectID)
	if err != nil {
		log.Fatal(err)
	}
	s := &server{
		ctx:         ctx,
		db:          dbWrapper,
		searcher:    search.NewElasticSearcher(*elasticAddress),
		statsReader: sr,
	}
	r := mux.NewRouter()
	r.HandleFunc("/api/lessons", s.APIServeLessons).Methods("GET")
	r.HandleFunc("/api/search", s.APIServeSearch).Methods("GET")
	r.HandleFunc("/api/stats", s.APIServeStats).Methods("GET")
	r.Handle("/healthz", health.NewElasticHealthChecker(*elasticAddress)).Methods("GET")
	appHandler := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "dist/index.html")
	}
	for _, route := range []string{"/search", "/lesson{*}", "/about", "/"} {
		r.Handle(route, http.HandlerFunc(appHandler)).Methods("GET")
	}
	r.Handle(
		"/{[a-z0-9.]+.(js|html|css)}",
		http.FileServer(http.Dir("dist"))).Methods("GET")

	log.Println("Serving on", *hostPort)
	srv := &http.Server{
		Handler: handlers.CombinedLoggingHandler(os.Stderr, handlers.CompressHandler(hsts.NewHandler(r))),
		Addr:    *hostPort,
		// Good practice: enforce timeouts.
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
