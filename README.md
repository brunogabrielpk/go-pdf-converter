# PDF Converter Web Application

A web application built in Go for converting files (images and text documents) to PDF format.

## Features

- **Single File Upload**: Upload one file at a time and download the converted PDF
- **Multiple File Upload**: Upload multiple files and download them as a ZIP archive
- **Supported Formats**:
  - Images: JPG, JPEG, PNG
  - Documents: TXT
- **SQLite Storage**: All converted PDFs are stored in an SQLite database
- **Modern Web Interface**: Drag-and-drop file upload with real-time feedback

## Project Structure

```
pdf-converter/
├── main.go           # HTTP server and routing
├── database.go       # SQLite database operations
├── converter.go      # PDF conversion logic
├── handlers.go       # HTTP request handlers
├── static/
│   └── index.html    # Frontend interface
├── go.mod            # Go module file
└── files.db          # SQLite database (created at runtime)
```

## Installation

1. Make sure you have Go installed (version 1.16 or later)

2. Navigate to the project directory:
```bash
cd pdf-converter
```

3. Dependencies are already installed. If needed, run:
```bash
go mod download
```

## Running the Application

1. Start the server:
```bash
./pdf-converter
```

Or build and run:
```bash
go run .
```

2. Open your browser and navigate to:
```
http://localhost:8080
```

## Usage

### Single File Conversion
1. Click the upload area or drag and drop a file
2. Select a supported file (JPG, JPEG, PNG, or TXT)
3. Click "Convert to PDF"
4. Download the converted PDF file

### Multiple File Conversion
1. Click the upload area or drag and drop multiple files
2. Select multiple supported files
3. Click "Convert to PDF"
4. Download all converted PDFs as a ZIP archive

## API Endpoints

### POST /upload
Upload one or more files for conversion to PDF.

**Request**: multipart/form-data with "files" field

**Response**:
```json
{
  "success": true,
  "message": "File uploaded and converted successfully",
  "file_id": 1
}
```

For multiple files:
```json
{
  "success": true,
  "message": "3 files uploaded and converted successfully",
  "file_ids": [1, 2, 3]
}
```

### GET /download?id={file_id}
Download a single converted PDF file.

**Parameters**:
- `id`: File ID returned from upload

### GET /download-zip?ids={file_ids}
Download multiple PDFs as a ZIP archive.

**Parameters**:
- `ids`: Comma-separated list of file IDs (e.g., "1,2,3")

## Technical Details

### Dependencies
- `github.com/mattn/go-sqlite3`: SQLite database driver
- `github.com/jung-kurt/gofpdf`: PDF generation library

### Image to PDF Conversion
- Images are decoded and embedded in PDF pages
- Page size automatically adjusts to image dimensions
- Maximum size limited to A4 to prevent oversized PDFs

### Text to PDF Conversion
- Text files are rendered with Courier font
- Automatic line wrapping for long lines
- Standard A4 page format

### Database Schema
```sql
CREATE TABLE files (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    original_name TEXT NOT NULL,
    pdf_data BLOB NOT NULL,
    uploaded_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

## Configuration

- **Port**: Default is 8080 (can be changed in main.go:28)
- **Database**: SQLite file at `./files.db`
- **Max Upload Size**: 32 MB (can be changed in handlers.go:33)

## Limitations

- Maximum upload size: 32 MB per request
- Supported file formats: JPG, JPEG, PNG, TXT
- Text files use fixed font (Courier)
- Images are resized to fit A4 if larger

## License

This project is open source and available for educational purposes.
