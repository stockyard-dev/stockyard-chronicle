package server
import ("encoding/json";"log";"net/http";"github.com/stockyard-dev/stockyard-chronicle/internal/store")
type Server struct{db *store.DB;mux *http.ServeMux;limits Limits}
func New(db *store.DB,limits Limits)*Server{s:=&Server{db:db,mux:http.NewServeMux(),limits:limits}
s.mux.HandleFunc("POST /api/events",s.emit);s.mux.HandleFunc("GET /api/events",s.query)
s.mux.HandleFunc("GET /api/types",s.topTypes);s.mux.HandleFunc("GET /api/sources",s.sources)
s.mux.HandleFunc("GET /api/stats",s.stats);s.mux.HandleFunc("GET /api/health",s.health)
s.mux.HandleFunc("GET /ui",s.dashboard);s.mux.HandleFunc("GET /ui/",s.dashboard);s.mux.HandleFunc("GET /",s.root);return s}
func(s *Server)ServeHTTP(w http.ResponseWriter,r *http.Request){s.mux.ServeHTTP(w,r)}
func wj(w http.ResponseWriter,c int,v any){w.Header().Set("Content-Type","application/json");w.WriteHeader(c);json.NewEncoder(w).Encode(v)}
func we(w http.ResponseWriter,c int,m string){wj(w,c,map[string]string{"error":m})}
func(s *Server)root(w http.ResponseWriter,r *http.Request){if r.URL.Path!="/"{http.NotFound(w,r);return};http.Redirect(w,r,"/ui",302)}
func(s *Server)emit(w http.ResponseWriter,r *http.Request){var e store.Event;json.NewDecoder(r.Body).Decode(&e);if e.Type==""{we(w,400,"type required");return};s.db.Emit(&e);wj(w,201,e)}
func(s *Server)query(w http.ResponseWriter,r *http.Request){q:=r.URL.Query();wj(w,200,map[string]any{"events":oe(s.db.Query(q.Get("type"),q.Get("source"),q.Get("severity"),q.Get("search"),200))})}
func(s *Server)topTypes(w http.ResponseWriter,r *http.Request){wj(w,200,map[string]any{"types":oe(s.db.TopTypes(20))})}
func(s *Server)sources(w http.ResponseWriter,r *http.Request){wj(w,200,map[string]any{"sources":oe(s.db.Sources())})}
func(s *Server)stats(w http.ResponseWriter,r *http.Request){wj(w,200,s.db.Stats())}
func(s *Server)health(w http.ResponseWriter,r *http.Request){st:=s.db.Stats();wj(w,200,map[string]any{"status":"ok","service":"chronicle","events":st.Total,"today":st.Today})}
func oe[T any](s []T)[]T{if s==nil{return[]T{}};return s}
func init(){log.SetFlags(log.LstdFlags|log.Lshortfile)}
