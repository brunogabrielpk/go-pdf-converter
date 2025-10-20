package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type FileRecord struct {
	ID           int
	OriginalName string
	PDFData      []byte
	UploadedAt   time.Time
}

type Database struct {
	db *sql.DB
}

func NewDatabase(filepath string) (*Database, error) {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, err
	}

	// Create table if not exists
	createTableSQL := `CREATE TABLE IF NOT EXISTS files (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		original_name TEXT NOT NULL,
		pdf_data BLOB NOT NULL,
		uploaded_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, err
	}

	return &Database{db: db}, nil
}

func (d *Database) SaveFile(originalName string, pdfData []byte) (int64, error) {
	result, err := d.db.Exec(
		"INSERT INTO files (original_name, pdf_data) VALUES (?, ?)",
		originalName, pdfData,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (d *Database) GetFile(id int64) (*FileRecord, error) {
	var record FileRecord
	err := d.db.QueryRow(
		"SELECT id, original_name, pdf_data, uploaded_at FROM files WHERE id = ?",
		id,
	).Scan(&record.ID, &record.OriginalName, &record.PDFData, &record.UploadedAt)

	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (d *Database) GetFiles(ids []int64) ([]*FileRecord, error) {
	records := make([]*FileRecord, 0, len(ids))

	for _, id := range ids {
		record, err := d.GetFile(id)
		if err != nil {
			log.Printf("Error fetching file ID %d: %v", id, err)
			continue
		}
		records = append(records, record)
	}

	return records, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}
