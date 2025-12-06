package main

import (
	"database/sql"
	"log"
	"time"

	"fmt"
	"os"

	// _ "github.com/mattn/go-sqlite3"
	_ "github.com/lib/pq"
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
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Create table if not exists
	// Note: Postgres uses SERIAL for auto-incrementing integer
	// and BYTEA for binary data
	createTableSQL := `CREATE TABLE IF NOT EXISTS files (
		id SERIAL PRIMARY KEY,
		original_name TEXT NOT NULL,
		pdf_data BYTEA NOT NULL,
		uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, err
	}

	return &Database{db: db}, nil
}

func (d *Database) SaveFile(originalName string, pdfData []byte) (int64, error) {
	var id int64
	// Postgres uses $1, $2 placeholders and RETURNING to get the id
	err := d.db.QueryRow(
		"INSERT INTO files (original_name, pdf_data) VALUES ($1, $2) RETURNING id",
		originalName, pdfData,
	).Scan(&id)

	if err != nil {
		return 0, err
	}
	return id, nil
}

func (d *Database) GetFile(id int64) (*FileRecord, error) {
	var record FileRecord
	// Postgres ises $1 placeholder
	err := d.db.QueryRow(
		"SELECT id, original_name, pdf_data, uploaded_at FROM files WHERE id = $1",
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
