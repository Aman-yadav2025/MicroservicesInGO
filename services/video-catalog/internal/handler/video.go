package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
	"video-catalog/internal/service"

	"video-catalog/internal/storage"
)

type UploadVideoRequest struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	VideoUrl     string `json:"video_url"`
	ThumbnailUrl string `json:"thumbnail_url"`
}

var MAX_UPLOAD_SIZE = 500 << 20 // max 500MB

var Storage storage.Provider

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

	file, header, err := r.FormFile("video_file")

	if err != nil {
		http.Error(w, "Failed to retrive Video file", http.StatusBadRequest)
		return
	}
	defer file.Close()
	// //give unique name
	// fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(header.Filename))
	// //unique save path for local storage
	// savePath := filepath.Join("uploads", "videos", fileName)

	// dst, err := os.Create(savePath)

	// if err != nil {
	// 	http.Error(w, "failed to save file", http.StatusInternalServerError)
	// 	return
	// }
	// defer dst.Close()

	// if _, err := io.Copy(dst, file); err != nil {
	// 	http.Error(w, "failed to write file", http.StatusInternalServerError)
	// 	return
	// }
	//common for both minio and local storage
	title := r.FormValue("title")
	description := r.FormValue("description")
	if title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}
	uniqueFileName := fmt.Sprintf("%d%s",time.Now().UnixNano(),filepath.Ext(header.Filename))
	size := header.Size

	contentType := header.Header.Get("Content-Type")

	if contentType == ""{
		contentType = "application/octet-stream"
	}

	videoUrl,err := Storage.UploadFile(r.Context(),file,uniqueFileName,size,contentType)
	
	if err != nil {
        http.Error(w, "Failed to upload file to storage", http.StatusInternalServerError)
        return
    }


	userID := r.Context().Value("user_id").(uint)

	video, err := service.CreateVideo(title, description, videoUrl, "", userID)

	if err != nil {
		http.Error(w, "failed to save to database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(fmt.Appendf(nil, "Video %s is uploaded successfully", video.Title))
}

func GetPage(w http.ResponseWriter, r *http.Request) {
	page := 1
	limit := 10
	// GET /videos?page=2&limit=5\
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if ParsePage, err := strconv.Atoi(pageStr); err == nil && ParsePage > 0 {
			page = ParsePage
		}
	}
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limitParse, err := strconv.Atoi(limitStr); err == nil && limitParse > 0 && limitParse < 50 {
			limit = limitParse
		}
	}

	videos, total, err := service.GetVideos(page, limit)

	if err != nil {
		http.Error(w, "Failed to fetch Videos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":        videos,
		"page":        page,
		"limit":       limit,
		"total_count": total,
	})
}
