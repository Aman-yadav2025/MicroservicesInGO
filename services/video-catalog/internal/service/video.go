package service

import (
	"video-catalog/internal/model"
	"video-catalog/internal/repository"
)

func CreateVideo(title, description, videoUrl, thumbnailUrl string, userId uint) (model.Video,error){
	video := model.Video{
		Title:        title,
		Description:  description,
		VideoUrl:     videoUrl,
		ThrumnailUrl: thumbnailUrl,
		UserId:       userId,
		Views:        0,
	}
	err :=repository.DB.Create(&video).Error
	
	return video,err
}

func GetAllVideos() ([]model.Video,error){
	var videos []model.Video
	err := repository.DB.Order("created_at desc").Find(&videos).Error
	return videos,err
}
