package httpapi

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

const maxPhotoUploadSize = 5 << 20

func (r *Router) handleUploadFragrancePhoto(w http.ResponseWriter, req *http.Request) {
	if _, ok := r.requireAdmin(w, req); !ok {
		return
	}

	req.Body = http.MaxBytesReader(w, req.Body, maxPhotoUploadSize)
	if err := req.ParseMultipartForm(maxPhotoUploadSize); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_upload", "Photo must be a multipart file up to 5 MB.")
		return
	}

	file, _, err := req.FormFile("photo")
	if err != nil {
		writeError(w, http.StatusBadRequest, "photo_required", "Multipart field photo is required.")
		return
	}
	defer file.Close()

	header := make([]byte, 512)
	bytesRead, err := file.Read(header)
	if err != nil && err != io.EOF {
		writeError(w, http.StatusBadRequest, "invalid_photo", "Failed to read uploaded photo.")
		return
	}
	header = header[:bytesRead]

	contentType := http.DetectContentType(header)
	extension, ok := allowedPhotoExtension(contentType)
	if !ok {
		writeError(w, http.StatusBadRequest, "unsupported_photo_type", "Photo must be JPEG, PNG, or WEBP.")
		return
	}

	if err := os.MkdirAll(r.cfg.UploadDir, 0o755); err != nil {
		r.logger.Error("failed to create upload directory", "error", err)
		writeError(w, http.StatusInternalServerError, "internal_error", "Failed to prepare upload directory.")
		return
	}

	fileName := uuid.NewString() + extension
	destinationPath := filepath.Join(r.cfg.UploadDir, fileName)
	destination, err := os.OpenFile(destinationPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		r.logger.Error("failed to create photo file", "error", err)
		writeError(w, http.StatusInternalServerError, "internal_error", "Failed to save uploaded photo.")
		return
	}
	defer destination.Close()

	if _, err := io.Copy(destination, io.MultiReader(bytes.NewReader(header), file)); err != nil {
		r.logger.Error("failed to write photo file", "error", err)
		writeError(w, http.StatusInternalServerError, "internal_error", "Failed to save uploaded photo.")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{
		"imageUrl": "/uploads/" + fileName,
	})
}

func allowedPhotoExtension(contentType string) (string, bool) {
	switch contentType {
	case "image/jpeg":
		return ".jpg", true
	case "image/png":
		return ".png", true
	case "image/webp":
		return ".webp", true
	default:
		return "", false
	}
}
