package store
import ("database/sql";"fmt";"os";"path/filepath";"strings";"time";_ "modernc.org/sqlite")
type DB struct{db *sql.DB}
type Event struct{ID string `json:"id"`;Type string `json:"type"`;Source string `json:"source,omitempty"`;Subject string `json:"subject,omitempty"`;Data string `json:"data,omitempty"`;Severity string `json:"severity,omitempty"`;Tags string `json:"tags,omitempty"`;CreatedAt string `json:"created_at"`}
type TypeCount struct{Type string `json:"type"`;Count int `json:"count"`}
func Open(d string)(*DB,error){if err:=os.MkdirAll(d,0755);err!=nil{return nil,err};db,err:=sql.Open("sqlite",filepath.Join(d,"chronicle.db")+"?_journal_mode=WAL&_busy_timeout=5000");if err!=nil{return nil,err}
db.Exec(`CREATE TABLE IF NOT EXISTS events(id TEXT PRIMARY KEY,type TEXT NOT NULL,source TEXT DEFAULT '',subject TEXT DEFAULT '',data TEXT DEFAULT '',severity TEXT DEFAULT 'info',tags TEXT DEFAULT '',created_at TEXT DEFAULT(datetime('now')))`)
db.Exec(`CREATE INDEX IF NOT EXISTS idx_events_type ON events(type)`)
db.Exec(`CREATE INDEX IF NOT EXISTS idx_events_source ON events(source)`)
db.Exec(`CREATE INDEX IF NOT EXISTS idx_events_date ON events(created_at)`)
return &DB{db:db},nil}
func(d *DB)Close()error{return d.db.Close()}
func genID()string{return fmt.Sprintf("%d",time.Now().UnixNano())}
func now()string{return time.Now().UTC().Format(time.RFC3339)}
func(d *DB)Emit(e *Event)error{e.ID=genID();e.CreatedAt=now();if e.Severity==""{e.Severity="info"}
_,err:=d.db.Exec(`INSERT INTO events VALUES(?,?,?,?,?,?,?,?)`,e.ID,e.Type,e.Source,e.Subject,e.Data,e.Severity,e.Tags,e.CreatedAt);return err}
func(d *DB)Query(typ,source,severity,search string,limit int)[]Event{if limit<=0{limit=100};where:=[]string{"1=1"};args:=[]any{}
if typ!=""{where=append(where,"type=?");args=append(args,typ)}
if source!=""{where=append(where,"source=?");args=append(args,source)}
if severity!=""{where=append(where,"severity=?");args=append(args,severity)}
if search!=""{where=append(where,"(subject LIKE ? OR data LIKE ?)");s:="%"+search+"%";args=append(args,s,s)}
rows,_:=d.db.Query(`SELECT * FROM events WHERE `+strings.Join(where," AND ")+` ORDER BY created_at DESC LIMIT ?`,append(args,limit)...)
if rows==nil{return nil};defer rows.Close()
var o []Event;for rows.Next(){var e Event;rows.Scan(&e.ID,&e.Type,&e.Source,&e.Subject,&e.Data,&e.Severity,&e.Tags,&e.CreatedAt);o=append(o,e)};return o}
func(d *DB)TopTypes(limit int)[]TypeCount{if limit<=0{limit=20};rows,_:=d.db.Query(`SELECT type,COUNT(*) c FROM events GROUP BY type ORDER BY c DESC LIMIT ?`,limit);if rows==nil{return nil};defer rows.Close()
var o []TypeCount;for rows.Next(){var t TypeCount;rows.Scan(&t.Type,&t.Count);o=append(o,t)};return o}
func(d *DB)Sources()[]string{rows,_:=d.db.Query(`SELECT DISTINCT source FROM events WHERE source!='' ORDER BY source`);if rows==nil{return nil};defer rows.Close();var o []string;for rows.Next(){var s string;rows.Scan(&s);o=append(o,s)};return o}
type Stats struct{Total int `json:"total"`;Types int `json:"types"`;Sources int `json:"sources"`;Today int `json:"today"`}
func(d *DB)Stats()Stats{var s Stats;d.db.QueryRow(`SELECT COUNT(*) FROM events`).Scan(&s.Total);d.db.QueryRow(`SELECT COUNT(DISTINCT type) FROM events`).Scan(&s.Types);s.Sources=len(d.Sources())
today:=time.Now().Format("2006-01-02");d.db.QueryRow(`SELECT COUNT(*) FROM events WHERE created_at>=?`,today+"T00:00:00Z").Scan(&s.Today);return s}
