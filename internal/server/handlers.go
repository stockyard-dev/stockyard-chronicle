package server
import("encoding/json";"net/http";"strconv";"github.com/stockyard-dev/stockyard-chronicle/internal/store")
func(s *Server)handleListProjects(w http.ResponseWriter,r *http.Request){list,_:=s.db.ListProjects();if list==nil{list=[]store.Project{}};writeJSON(w,200,list)}
func(s *Server)handleCreateProject(w http.ResponseWriter,r *http.Request){var p store.Project;json.NewDecoder(r.Body).Decode(&p);if p.Name==""{writeError(w,400,"name required");return};if err:=s.db.CreateProject(&p);err!=nil{writeError(w,500,err.Error());return};writeJSON(w,201,p)}
func(s *Server)handleDeleteProject(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);s.db.DeleteProject(id);writeJSON(w,200,map[string]string{"status":"deleted"})}
func(s *Server)handleListEntries(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);list,_:=s.db.ListEntries(id);if list==nil{list=[]store.Entry{}};writeJSON(w,200,list)}
func(s *Server)handleCreateEntry(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);var e store.Entry;json.NewDecoder(r.Body).Decode(&e);e.ProjectID=id;if e.Title==""{writeError(w,400,"title required");return};if err:=s.db.CreateEntry(&e);err!=nil{writeError(w,500,err.Error());return};writeJSON(w,201,e)}
func(s *Server)handleDeleteEntry(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);s.db.DeleteEntry(id);writeJSON(w,200,map[string]string{"status":"deleted"})}
func(s *Server)handleStats(w http.ResponseWriter,r *http.Request){n,_:=s.db.CountEntries();writeJSON(w,200,map[string]interface{}{"entries":n})}
