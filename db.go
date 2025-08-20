package main

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

type AppDB struct {
	db *sql.DB
}


func NewAppDB(path string) *AppDB {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}

	if err := ensureSchema(db); err != nil {
		log.Fatalf("ensure schema: %v", err)
	}

	return &AppDB{db: db}
}

func ensureSchema(db *sql.DB) error {
	schema := `
CREATE TABLE IF NOT EXISTS messages (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  remote_addr TEXT NOT NULL,
  content TEXT NOT NULL,
  created_at DATETIME NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f','now'))
);
PRAGMA journal_mode=WAL;
PRAGMA synchronous=NORMAL;
`
	_, err := db.Exec(schema)
	return err
}

func (a *AppDB) InsertMessage(remote, content string) error {
	_, err := a.db.Exec(`INSERT INTO messages (remote_addr, content) VALUES (?, ?)`, remote, content)
	return err
}

func (a *AppDB) Close() {
	a.db.Close()
}

func (a *AppDB) GetAllMessages() ([]Message, error) {
	rows, err := a.db.Query(`
		SELECT id, remote_addr, content, created_at 
		FROM messages 
		ORDER BY id DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.ID, &msg.RemoteAddr, &msg.Content, &msg.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	// Cek error setelah iterasi
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

