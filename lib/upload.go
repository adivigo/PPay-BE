package lib

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func UploadImage(c *gin.Context, file *multipart.FileHeader, allowedExts []string, maxSize int64, uploadDir string) (string, error) {
	// 1. Validasi ukuran file
	if file.Size > maxSize {
		return "", fmt.Errorf("file size exceeds the maximum limit of %d MB", maxSize/(1024*1024))
	}

	// 2. Validasi ekstensi file
	ext := strings.ToLower(filepath.Ext(file.Filename))
	valid := false
	for _, allowedExt := range allowedExts {
		if ext == allowedExt {
			valid = true
			break
		}
	}
	if !valid {
		return "", fmt.Errorf("invalid file type: only %v are allowed", allowedExts)
	}

	// 3. Proses nama file untuk menghapus spasi
	cleanFileName := strings.ReplaceAll(file.Filename, " ", "_") // Ganti spasi dengan underscore
	newFileName := fmt.Sprintf("%d-%s", time.Now().Unix(), cleanFileName)

	// 4. Buat direktori penyimpanan jika belum ada
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %v", err)
	}

	// 5. Simpan file menggunakan c.SaveUploadedFile
	savePath := filepath.Join(uploadDir, newFileName)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		return "", fmt.Errorf("failed to save file: %v", err)
	}

	return savePath, nil
}
