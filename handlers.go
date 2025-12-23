package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

type Handler struct {
	db *Database
}

func NewHandler(db *Database) *Handler {
	return &Handler{db: db}
}

type UploadResponse struct {
	Success bool    `json:"success"`
	Message string  `json:"message"`
	FileIDs []int64 `json:"file_ids,omitempty"`
	FileID  int64   `json:"file_id,omitempty"`
}

// HandleUpload handles both single and multiple file uploads
func (h *Handler) HandleUpload(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request: %s %s", r.Method, r.URL.Path)
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form (32 MB max)
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, UploadResponse{
			Success: false,
			Message: "Failed to parse form",
		})
		return
	}

	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		respondJSON(w, http.StatusBadRequest, UploadResponse{
			Success: false,
			Message: "No files uploaded",
		})
		return
	}

	var fileIDs []int64

	for _, fileHeader := range files {
		// Open uploaded file
		file, err := fileHeader.Open()
		if err != nil {
			log.Printf("Error opening file %s: %v", fileHeader.Filename, err)
			continue
		}

		// Read file data
		data, err := io.ReadAll(file)
		file.Close()
		if err != nil {
			log.Printf("Error reading file %s: %v", fileHeader.Filename, err)
			continue
		}

		// Convert to PDF
		pdfData, err := ConvertToPDF(fileHeader.Filename, data)
		if err != nil {
			log.Printf("Error converting file %s: %v", fileHeader.Filename, err)
			respondJSON(w, http.StatusBadRequest, UploadResponse{
				Success: false,
				Message: fmt.Sprintf("Failed to convert %s: %v", fileHeader.Filename, err),
			})
			return
		}

		// Save to database
		fileID, err := h.db.SaveFile(fileHeader.Filename, pdfData)
		if err != nil {
			log.Printf("Error saving file %s: %v", fileHeader.Filename, err)
			continue
		}

		fileIDs = append(fileIDs, fileID)
	}

	if len(fileIDs) == 0 {
		respondJSON(w, http.StatusInternalServerError, UploadResponse{
			Success: false,
			Message: "Failed to process any files",
		})
		return
	}

	// If single file, return single file ID
	if len(fileIDs) == 1 {
		respondJSON(w, http.StatusOK, UploadResponse{
			Success: true,
			Message: "File uploaded and converted successfully",
			FileID:  fileIDs[0],
		})
		return
	}

	// Multiple files
	respondJSON(w, http.StatusOK, UploadResponse{
		Success: true,
		Message: fmt.Sprintf("%d files uploaded and converted successfully", len(fileIDs)),
		FileIDs: fileIDs,
	})
}

// HandleDownload handles single PDF download
func (h *Handler) HandleDownload(w http.ResponseWriter, r *http.Request) {
	// Get file ID from URL
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "File ID required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid file ID", http.StatusBadRequest)
		return
	}

	// Get file from database
	record, err := h.db.GetFile(id)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Generate PDF filename
	pdfFilename := strings.TrimSuffix(record.OriginalName, filepath.Ext(record.OriginalName)) + ".pdf"

	// Set headers for PDF download
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", pdfFilename))
	w.Header().Set("Content-Length", strconv.Itoa(len(record.PDFData)))

	// Write PDF data
	_, err = w.Write(record.PDFData)
	if err != nil {
		log.Printf("Error writing PDF data: %v", err)
	}
}

// HandleDownloadZip handles multiple PDF downloads as ZIP
func (h *Handler) HandleDownloadZip(w http.ResponseWriter, r *http.Request) {
	// Get file IDs from URL (comma-separated)
	idsStr := r.URL.Query().Get("ids")
	if idsStr == "" {
		http.Error(w, "File IDs required", http.StatusBadRequest)
		return
	}

	idStrs := strings.Split(idsStr, ",")
	var ids []int64

	for _, idStr := range idStrs {
		id, err := strconv.ParseInt(strings.TrimSpace(idStr), 10, 64)
		if err != nil {
			continue
		}
		ids = append(ids, id)
	}

	if len(ids) == 0 {
		http.Error(w, "No valid file IDs", http.StatusBadRequest)
		return
	}

	// Get files from database
	records, err := h.db.GetFiles(ids)
	if err != nil {
		http.Error(w, "Error retrieving files", http.StatusInternalServerError)
		return
	}

	if len(records) == 0 {
		http.Error(w, "No files found", http.StatusNotFound)
		return
	}

	// Prepare files map for ZIP
	files := make(map[string][]byte)
	for _, record := range records {
		pdfFilename := strings.TrimSuffix(record.OriginalName, filepath.Ext(record.OriginalName)) + ".pdf"
		files[pdfFilename] = record.PDFData
	}

	// Create ZIP
	zipData, err := CreateZip(files)
	if err != nil {
		http.Error(w, "Error creating ZIP", http.StatusInternalServerError)
		log.Printf("Error creating ZIP: %v", err)
		return
	}

	// Set headers for ZIP download
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=\"converted_pdfs.zip\"")
	w.Header().Set("Content-Length", strconv.Itoa(len(zipData)))

	// Write ZIP data
	_, err = w.Write(zipData)
	if err != nil {
		log.Printf("Error writing ZIP data: %v", err)
	}
}

// respondJSON sends a JSON response
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
