package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"greatmanbackend/common"
	"greatmanbackend/model"
	"greatmanbackend/util"
	"net/http"
)

func GetGreatMan(c *gin.Context) {
	var greatMan model.GreatMan
	greatManId, err := util.ParseUint(c.Param("id"))
	lang := c.Query("lang")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	if err = common.GetDB().
		Preload("GreatManInfos", &model.GreatManInfo{
			Language: lang,
		}).
		Take(&greatMan, greatManId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"name":       greatMan.GreatManInfos[0].Name,
		"comment":    greatMan.GreatManInfos[0].Comment,
		"headImgUrl": greatMan.HeadImgUrl,
	})
}
