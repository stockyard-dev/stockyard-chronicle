package store

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

type DB struct{ db *sql.DB }

// Event is an append-only log entry. Once emitted, an event cannot be
// modified or deleted through the API. This is intentional — chronicle is
// a journal, not a CRUD app.
type Event struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Source    string `json:"source,omitempty"`
	Subject   string `json:"subject,omitempty"`
	Data      string `json:"data,omitempty"`
	Severity  string `json:"severity,omitempty"`
	Tags      string `json:"tags,omitempty"`
	CreatedAt string `json:"created_at"`
}

type TypeCount struct {
	Type  string `json:"type"`
	Count int    `json:"count"`
}

func Open(d string) (*DB, error) {
	if err := os.MkdirAll(d, 0755); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", filepath.Join(d, "chronicle.db")+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, err
	}
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS events(
			id TEXT PRIMARY KEY,
			type TEXT NOT NULL,
			source TEXT DEFAULT '',
			subject TEXT DEFAULT '',
			data TEXT DEFAULT '',
			severity TEXT DEFAULT 'info',
			tags TEXT DEFAULT '',
			created_at TEXT DEFAULT(datetime('now'))
		)`,
		`CREATE INDEX IF NOT EXISTS idx_events_type ON events(type)`,
		`CREATE INDEX IF NOT EXISTS idx_events_source ON events(source)`,
		`CREATE INDEX IF NOT EXISTS idx_events_severity ON events(severity)`,
		`CREATE INDEX IF NOT EXISTS idx_events_date ON events(created_at)`,
	}
	for _, q := range migrations {
		if _, err := db.Exec(q); err != nil {
			return nil, fmt.Errorf("migrate: %w", err)
		}
	}
	return &DB{db: db}, nil
}

func (d *DB) Close() error { return d.db.Close() }

func genID() string { return fmt.Sprintf("%d", time.Now().UnixNano()) }
func now() string   { return time.Now().UTC().Format(time.RFC3339) }

// Emit appends a new event to the log. ID and CreatedAt are set automatically.
// Severity defaults to "info" if empty.
func (d *DB) Emit(e *Event) error {
	e.ID = genID()
	e.CreatedAt = now()
	if e.Severity == "" {
		e.Severity = "info"
	}
	_, err := d.db.Exec(
		`INSERT INTO events(id, type, source, subject, data, severity, tags, created_at)
		 VALUES(?, ?, ?, ?, ?, ?, ?, ?)`,
		e.ID, e.Type, e.Source, e.Subject, e.Data, e.Severity, e.Tags, e.CreatedAt,
	)
	return err
}

// Query returns events matching the given filters. Limit defaults to 100.
// Empty filter strings match everything for that field.
func (d *DB) Query(typ, source, severity, search string, limit int) []Event {
	if limit <= 0 {
		limit = 100
	}
	where := []string{"1=1"}
	args := []any{}
	if typ != "" {
		where = append(where, "type=?")
		args = append(args, typ)
	}
	if source != "" {
		where = append(where, "source=?")
		args = append(args, source)
	}
	if severity != "" {
		where = append(where, "severity=?")
		args = append(args, severity)
	}
	if search != "" {
		where = append(where, "(subject LIKE ? OR data LIKE ? OR tags LIKE ?)")
		s := "%" + search + "%"
		args = append(args, s, s, s)
	}
	args = append(args, limit)
	rows, _ := d.db.Query(
		`SELECT id, type, source, subject, data, severity, tags, created_at
		 FROM events WHERE `+strings.Join(where, " AND ")+`
		 ORDER BY created_at DESC LIMIT ?`,
		args...,
	)
	if rows == nil {
		return nil
	}
	defer rows.Close()
	var o []Event
	for rows.Next() {
		var e Event
		rows.Scan(&e.ID, &e.Type, &e.Source, &e.Subject, &e.Data, &e.Severity, &e.Tags, &e.CreatedAt)
		o = append(o, e)
	}
	return o
}

// TopTypes returns the top N event types by count, ordered descending.
func (d *DB) TopTypes(limit int) []TypeCount {
	if limit <= 0 {
		limit = 20
	}
	rows, _ := d.db.Query(
		`SELECT type, COUNT(*) c FROM events GROUP BY type ORDER BY c DESC LIMIT ?`,
		limit,
	)
	if rows == nil {
		return nil
	}
	defer rows.Close()
	var o []TypeCount
	for rows.Next() {
		var t TypeCount
		rows.Scan(&t.Type, &t.Count)
		o = append(o, t)
	}
	return o
}

// Sources returns the distinct source values seen so far, alphabetized.
func (d *DB) Sources() []string {
	rows, _ := d.db.Query(`SELECT DISTINCT source FROM events WHERE source != '' ORDER BY source`)
	if rows == nil {
		return nil
	}
	defer rows.Close()
	var o []string
	for rows.Next() {
		var s string
		rows.Scan(&s)
		o = append(o, s)
	}
	return o
}

// Stats returns aggregate counts for the dashboard cards.
// Today is computed in UTC.
func (d *DB) Stats() map[string]any {
	m := map[string]any{
		"total":       0,
		"types":       0,
		"sources":     0,
		"today":       0,
		"by_severity": map[string]int{},
		"by_type":     map[string]int{},
	}

	var total int
	d.db.QueryRow(`SELECT COUNT(*) FROM events`).Scan(&total)
	m["total"] = total

	var typeCount int
	d.db.QueryRow(`SELECT COUNT(DISTINCT type) FROM events`).Scan(&typeCount)
	m["types"] = typeCount

	m["sources"] = len(d.Sources())

	today := time.Now().UTC().Format("2006-01-02") + "T00:00:00Z"
	var todayCount int
	d.db.QueryRow(`SELECT COUNT(*) FROM events WHERE created_at >= ?`, today).Scan(&todayCount)
	m["today"] = todayCount

	if rows, _ := d.db.Query(`SELECT severity, COUNT(*) FROM events GROUP BY severity`); rows != nil {
		defer rows.Close()
		by := map[string]int{}
		for rows.Next() {
			var k string
			var c int
			rows.Scan(&k, &c)
			by[k] = c
		}
		m["by_severity"] = by
	}

	// by_type returns the same data as TopTypes but as a map for convenience
	// in the dashboard rendering.
	tt := d.TopTypes(50)
	by := map[string]int{}
	for _, t := range tt {
		by[t.Type] = t.Count
	}
	m["by_type"] = by

	return m
}
