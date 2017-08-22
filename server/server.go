package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/attwad/cdf-fe/server/db"
	"github.com/attwad/cdf-fe/server/health"
	"github.com/attwad/cdf-fe/server/search"
	"github.com/attwad/cdf/data"
)

var (
	hostPort       = flag.String("listen_addr", "127.0.0.1:8080", "Address to listen on.")
	projectID      = flag.String("project_id", "college-de-france", "Google cloud project.")
	elasticAddress = flag.String("elastic_address", "", "HTTP address to elastic instance")
)

type server struct {
	ctx      context.Context
	db       db.Wrapper
	searcher search.Searcher
}

func (s *server) APIServeLessons(w http.ResponseWriter, r *http.Request) {
	filter := db.FilterNone
	if r.URL.Query().Get("filter") == "converted" {
		filter = db.FilterOnlyConverted
	}
	lessons, cursor, err := s.db.GetLessons(s.ctx, r.URL.Query().Get("cursor"), filter)
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
	jsr, err := s.searcher.Search(q)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	type searchResponse struct {
		Query    string          `json:"query"`
		TookMs   int             `json:"took_ms"`
		TimedOut bool            `json:"timed_out"`
		Sources  []search.Source `json:"sources"`
	}
	sr := searchResponse{
		Query:    q,
		TookMs:   jsr.TookMs,
		TimedOut: jsr.TimedOut,
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

func main() {
	flag.Parse()
	ctx := context.Background()
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	dbWrapper, err := db.NewDatastoreWrapper(ctx, *projectID)
	if err != nil {
		log.Fatalf("creating db wrapper: %v", err)
	}
	log.Println("Will connect to elastic instance @", *elasticAddress)
	s := &server{
		ctx:      ctx,
		db:       dbWrapper,
		searcher: search.NewElasticSearcher(*elasticAddress),
	}
	http.HandleFunc("/api/lessons", s.APIServeLessons)
	http.HandleFunc("/api/search", s.APIServeSearch)
	http.Handle("/healthz", health.NewElasticHealthChecker(*elasticAddress))
	http.Handle("/", http.FileServer(http.Dir("dist")))
	log.Println("Serving on", *hostPort)

	log.Fatal(http.ListenAndServe(*hostPort, nil))
}
