package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"video-catalog/internal/service"
)

type UploadVideoRequest struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	VideoUrl     string `json:"video_url"`
	ThumbnailUrl string `json:"thumbnail_url"`
}

var MAX_UPLOAD_SIZE = 500 << 20 // max 500MB

func UploadVideo(w http.ResponseWriter, r *http.Request) {
	var req UploadVideoRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}
	userId := r.Context().Value("user_id").(uint)

	video, err := service.CreateVideo(req.Title, req.Description, req.VideoUrl, req.ThumbnailUrl, userId)
	if err != nil {
		http.Error(w, "Failed to Create video entry", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(video)
}

func ListVideo(w http.ResponseWriter, r *http.Request) {
	videos, err := service.GetAllVideos()

	if err != nil {
		http.Error(w, "failed to fetch videos", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(videos)
}

func UploadVideoFile(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, int64(MAX_UPLOAD_SIZE))

	if err := r.ParseMultipartForm(int64(MAX_UPLOAD_SIZE)); err != nil {
		http.Error(w, "File must be under 500 MB", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("video-file")

	if err != nil {
		http.Error(w, "Failed to retrive Video file", http.StatusBadRequest)
		return
	}
	defer file.Close()
	//give unique name
	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(handler.Filename))
	//unique save path
	savePath := filepath.Join("uploads", "videos", fileName)

	dst, err := os.Create(savePath)

	if err != nil {
		http.Error(w, "failed to save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "failed to write file", http.StatusInternalServerError)
		return
	}
	title := r.FormValue("title")
	description := r.FormValue("description")

	userID := r.Context().Value("user_id").(uint)

	video, err := service.CreateVideo(title, description, savePath, "", userID)

	if err != nil {
		http.Error(w, "failed to save to database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(fmt.Appendf(nil, "Video %s is uploaded successfully", video.Title))
}
