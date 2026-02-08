package handler

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/zeromicro/go-zero/rest/httpx"
	"life-system-backend/internal/svc"
	"life-system-backend/internal/types"
)

const maxUploadSize = 5 << 20 // 5MB

var allowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

func UploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

		if err := r.ParseMultipartForm(maxUploadSize); err != nil {
			httpx.OkJson(w, types.CommonResp{Code: 400, Message: "文件过大，最大5MB"})
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			httpx.OkJson(w, types.CommonResp{Code: 400, Message: "请选择文件"})
			return
		}
		defer file.Close()

		// Validate file type
		contentType := header.Header.Get("Content-Type")
		if !allowedImageTypes[contentType] {
			httpx.OkJson(w, types.CommonResp{Code: 400, Message: "仅支持 JPEG、PNG、GIF、WebP 格式"})
			return
		}

		// Generate unique filename
		randBytes := make([]byte, 16)
		if _, err := rand.Read(randBytes); err != nil {
			httpx.OkJson(w, types.CommonResp{Code: 500, Message: "上传失败"})
			return
		}

		ext := filepath.Ext(header.Filename)
		if ext == "" {
			switch contentType {
			case "image/jpeg":
				ext = ".jpg"
			case "image/png":
				ext = ".png"
			case "image/gif":
				ext = ".gif"
			case "image/webp":
				ext = ".webp"
			}
		}
		ext = strings.ToLower(ext)
		filename := hex.EncodeToString(randBytes) + ext

		// Ensure upload directory exists
		uploadDir := "data/uploads"
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			httpx.OkJson(w, types.CommonResp{Code: 500, Message: "上传失败"})
			return
		}

		// Save file
		dstPath := filepath.Join(uploadDir, filename)
		dst, err := os.Create(dstPath)
		if err != nil {
			httpx.OkJson(w, types.CommonResp{Code: 500, Message: "上传失败"})
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			httpx.OkJson(w, types.CommonResp{Code: 500, Message: "上传失败"})
			return
		}

		// Return the URL path
		urlPath := fmt.Sprintf("/uploads/%s", filename)
		httpx.OkJson(w, types.CommonResp{
			Code:    0,
			Message: "success",
			Data: map[string]string{
				"url": urlPath,
			},
		})
	}
}
