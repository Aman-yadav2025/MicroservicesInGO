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

func GetVideos(page int , limit int) ([]model.Video,int64,error){
	var videos []model.Video
	var totalCount int64

	repository.DB.Model(model.Video{}).Count(&totalCount)

	offset := (page -1 )*limit

	err := repository.DB.Order("created_at desc").Limit(limit).Offset(offset).Find(&videos).Error

	return videos,totalCount,err
}
