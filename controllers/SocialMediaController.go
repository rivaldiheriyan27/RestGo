package controllers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SocialMediaDB struct {
	DB *gorm.DB
}

func (db *SocialMediaDB) CreateSocialMedia(c *gin.Context) {

}
