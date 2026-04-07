package server

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/stockyard-dev/stockyard-chronicle/internal/store"
)

type Server struct {
	db      *store.DB
	mux     *http.ServeMux
	limits  Limits
	dataDir string
	pCfg    map[string]json.RawMessage
}

func New(db *store.DB, limits Limits, dataDir string) *Server {
	s := &Server{
		db:      db,
		mux:     http.NewServeMux(),
		limits:  limits,
		dataDir: dataDir,
	}
	s.loadPersonalConfig()

	// Event log (append-only — no PUT, no DELETE)
	s.mux.HandleFunc("POST /api/events", s.emit)
	s.mux.HandleFunc("GET /api/events", s.query)

	// Aggregates
	s.mux.HandleFunc("GET /api/types", s.topTypes)
	s.mux.HandleFunc("GET /api/sources", s.sources)
	s.mux.HandleFunc("GET /api/stats", s.stats)
	s.mux.HandleFunc("GET /api/health", s.health)

	// Personalization
	s.mux.HandleFunc("GET /api/config", s.configHandler)

	// Tier
	s.mux.HandleFunc("GET /api/tier", func(w http.ResponseWriter, r *http.Request) {
		wj(w, 200, map[string]any{
			"tier":        s.limits.Tier,
			"upgrade_url": "https://stockyard.dev/chronicle/",
		})
	})

	// Dashboard
	s.mux.HandleFunc("GET /ui", s.dashboard)
	s.mux.HandleFunc("GET /ui/", s.dashboard)
	s.mux.HandleFunc("GET /", s.root)

	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

// ─── helpers ──────────────────────────────────────────────────────

func wj(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}

func we(w http.ResponseWriter, code int, msg string) {
	wj(w, code, map[string]string{"error": msg})
}

func oe[T any](s []T) []T {
	if s == nil {
		return []T{}
	}
	return s
}

func (s *Server) root(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, "/ui", 302)
}

// ─── personalization ──────────────────────────────────────────────

func (s *Server) loadPersonalConfig() {
	path := filepath.Join(s.dataDir, "config.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	var cfg map[string]json.RawMessage
	if err := json.Unmarshal(data, &cfg); err != nil {
		log.Printf("chronicle: warning: could not parse config.json: %v", err)
		return
	}
	s.pCfg = cfg
	log.Printf("chronicle: loaded personalization from %s", path)
}

func (s *Server) configHandler(w http.ResponseWriter, r *http.Request) {
	if s.pCfg == nil {
		wj(w, 200, map[string]any{})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.pCfg)
}

// ─── event log ────────────────────────────────────────────────────

func (s *Server) emit(w http.ResponseWriter, r *http.Request) {
	if s.limits.MaxItems > 0 {
		st := s.db.Stats()
		if t, ok := st["total"].(int); ok && t >= s.limits.MaxItems {
			we(w, 402, "Free tier limit reached. Upgrade at https://stockyard.dev/chronicle/")
			return
		}
	}
	var e store.Event
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		we(w, 400, "invalid json")
		return
	}
	if e.Type == "" {
		we(w, 400, "type required")
		return
	}
	if err := s.db.Emit(&e); err != nil {
		we(w, 500, "emit failed")
		return
	}
	wj(w, 201, e)
}

func (s *Server) query(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	limit := 200
	if l := q.Get("limit"); l != "" {
		if n, err := strconv.Atoi(l); err == nil && n > 0 && n <= 1000 {
			limit = n
		}
	}
	wj(w, 200, map[string]any{
		"events": oe(s.db.Query(
			q.Get("type"),
			q.Get("source"),
			q.Get("severity"),
			q.Get("search"),
			limit,
		)),
	})
}

func (s *Server) topTypes(w http.ResponseWriter, r *http.Request) {
	wj(w, 200, map[string]any{"types": s.db.TopTypes(20)})
}

func (s *Server) sources(w http.ResponseWriter, r *http.Request) {
	wj(w, 200, map[string]any{"sources": oe(s.db.Sources())})
}

func (s *Server) stats(w http.ResponseWriter, r *http.Request) {
	wj(w, 200, s.db.Stats())
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	st := s.db.Stats()
	wj(w, 200, map[string]any{
		"status":  "ok",
		"service": "chronicle",
		"events":  st["total"],
		"today":   st["today"],
	})
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
