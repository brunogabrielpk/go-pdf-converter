package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"path/filepath"
	"strings"

	"github.com/jung-kurt/gofpdf"
)

// ConvertToPDF converts various file types to PDF
func ConvertToPDF(filename string, data []byte) ([]byte, error) {
	ext := strings.ToLower(filepath.Ext(filename))

	switch ext {
	case ".jpg", ".jpeg", ".png":
		return convertImageToPDF(filename, data)
	case ".txt":
		return convertTextToPDF(filename, data)
	default:
		return nil, fmt.Errorf("unsupported file type: %s", ext)
	}
}

// convertImageToPDF converts image files to PDF
func convertImageToPDF(filename string, data []byte) ([]byte, error) {
	// Decode the image to get dimensions
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	bounds := img.Bounds()
	width := float64(bounds.Dx())
	height := float64(bounds.Dy())

	// Create PDF with custom page size matching image aspect ratio
	// Convert pixels to mm (assuming 72 DPI)
	const dpi = 72.0
	const mmPerInch = 25.4
	widthMM := (width / dpi) * mmPerInch
	heightMM := (height / dpi) * mmPerInch

	// Limit page size to A4 if too large
	maxWidth := 210.0  // A4 width in mm
	maxHeight := 297.0 // A4 height in mm

	if widthMM > maxWidth || heightMM > maxHeight {
		ratio := width / height
		if widthMM > heightMM {
			widthMM = maxWidth
			heightMM = widthMM / ratio
		} else {
			heightMM = maxHeight
			widthMM = heightMM * ratio
		}
	}

	pdf := gofpdf.NewCustom(&gofpdf.InitType{
		OrientationStr: "P",
		UnitStr:        "mm",
		SizeStr:        "",
		Size:           gofpdf.SizeType{Wd: widthMM, Ht: heightMM},
	})

	pdf.AddPage()

	// Register the image from memory
	ext := strings.ToLower(filepath.Ext(filename))
	imageType := "jpg"
	if ext == ".png" {
		imageType = "png"
	}

	imageReader := bytes.NewReader(data)
	imageInfo := pdf.RegisterImageOptionsReader(filename, gofpdf.ImageOptions{
		ImageType: imageType,
	}, imageReader)

	if pdf.Error() != nil {
		return nil, fmt.Errorf("failed to register image: %w", pdf.Error())
	}

	// Place image to fill the page
	pdf.ImageOptions(filename, 0, 0, widthMM, heightMM, false, gofpdf.ImageOptions{
		ImageType: imageType,
	}, 0, "")

	if pdf.Error() != nil {
		return nil, fmt.Errorf("failed to add image to PDF: %w", pdf.Error())
	}

	// Output PDF to buffer
	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to output PDF: %w", err)
	}

	_ = imageInfo // Use the variable to avoid unused error

	return buf.Bytes(), nil
}

// convertTextToPDF converts text files to PDF
func convertTextToPDF(filename string, data []byte) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Courier", "", 10)

	// Read text content
	text := string(data)

	// Split text into lines and add to PDF
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		// Handle long lines by wrapping
		pdf.MultiCell(0, 5, line, "", "", false)
	}

	if pdf.Error() != nil {
		return nil, fmt.Errorf("failed to create PDF from text: %w", pdf.Error())
	}

	// Output PDF to buffer
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to output PDF: %w", err)
	}

	return buf.Bytes(), nil
}

// CreateZip creates a ZIP archive containing multiple PDF files
func CreateZip(files map[string][]byte) ([]byte, error) {
	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)

	for filename, data := range files {
		fileWriter, err := zipWriter.Create(filename)
		if err != nil {
			return nil, fmt.Errorf("failed to create file in zip: %w", err)
		}

		_, err = io.Copy(fileWriter, bytes.NewReader(data))
		if err != nil {
			return nil, fmt.Errorf("failed to write file to zip: %w", err)
		}
	}

	if err := zipWriter.Close(); err != nil {
		return nil, fmt.Errorf("failed to close zip: %w", err)
	}

	return buf.Bytes(), nil
}
