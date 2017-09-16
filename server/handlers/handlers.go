package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/attwad/cdf-fe/server/search"
	"github.com/attwad/cdf/data"
	"github.com/attwad/cdf/db"
	"github.com/attwad/cdf/stats/io"
)

const pageSize = 10

// ExternalCourse is what gets sent to clients, it contains formatted durations, dates etc.
type ExternalCourse struct {
	data.Course
	Converted         bool   `json:"converted"`
	FormattedDate     string `json:"date"`
	FormattedDuration string `json:"duration"`
}

type statsHandler struct {
	statsReader io.Reader
}

// NewStatsHandler handles HTTP requests for stats.
func NewStatsHandler(statsReader io.Reader) http.Handler {
	return &statsHandler{statsReader}
}

type lessonsHandler struct {
	db db.Wrapper
}

// NewLessonsHandler handles HTTP requests for lessons browsing.
func NewLessonsHandler(db db.Wrapper) http.Handler {
	return &lessonsHandler{db}
}

type searchHandler struct {
	searcher search.Searcher
}

// NewSearchHandler handles HTTP requests for search queries.
func NewSearchHandler(searcher search.Searcher) http.Handler {
	return &searchHandler{searcher}
}

func (s *lessonsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	filter := db.FilterNone
	if r.URL.Query().Get("filter") == "converted" {
		filter = db.FilterOnlyConverted
	}
	lessons, cursor, err := s.db.GetLessons(r.Context(), r.URL.Query().Get("cursor"), filter, pageSize)
	if err != nil {
		log.Println("Could not read lessons from db:", err)
		http.Error(w, "Could not read lessons from DB", http.StatusInternalServerError)
		return
	}
	type response struct {
		Cursor  string           `json:"cursor"`
		Lessons []ExternalCourse `json:"lessons"`
	}
	resp := &response{Cursor: cursor, Lessons: make([]ExternalCourse, 0)}
	for _, lesson := range lessons {
		resp.Lessons = append(resp.Lessons, ExternalCourse{
			Course:            lesson.Course,
			Converted:         lesson.Converted,
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

func (s *searchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

func (s *statsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
