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

	if len(greatMan.GreatManInfos) == 0 {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"name":       greatMan.GreatManInfos[0].Name,
		"comment":    greatMan.GreatManInfos[0].Comment,
		"headImgUrl": greatMan.HeadImgUrl,
	})
}

func CreatNewMan(c *gin.Context) {
	man := Man{}

	if err := c.ShouldBind(&man); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	var infos []model.GreatManInfo

	for _, i := range man.ManInfos {
		infos = append(infos, model.GreatManInfo{
			Language: i.Language,
			Comment:  i.Comment,
			Name:     i.Name,
		})
	}

	g := model.GreatMan{
		HeadImgUrl:    man.HeadImgUrl,
		GreatManInfos: infos,
	}

	if err := common.GetDB().Create(&g).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusOK, g.ID)
}

func UpdateMan(c *gin.Context) {
	id, err := util.ParseUint(c.Param("id"))
	man := Man{}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	if err := c.ShouldBind(&man); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	var infos []model.GreatManInfo

	for _, i := range man.ManInfos {
		infos = append(infos, model.GreatManInfo{
			Language: i.Language,
			Comment:  i.Comment,
			Name:     i.Name,
		})
	}

	if err := common.
		GetDB().
		Omit("TalkRecords").
		Session(&gorm.Session{FullSaveAssociations: true}).
		Updates(&model.GreatMan{
			Model: gorm.Model{
				ID: id,
			},
			HeadImgUrl:    man.HeadImgUrl,
			GreatManInfos: infos,
		}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func DeleteMan(c *gin.Context) {
	id, err := util.ParseUint(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	if err := common.GetDB().Delete(&model.GreatMan{}, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func GetGreatMen(c *gin.Context) {
	var s SearchParam
	var men []model.GreatMan
	var results []uint

	if err := c.ShouldBind(&s); err != nil || s.Limit > 20 {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if err := common.
		GetDB().
		Order("created_at DESC").
		Limit(int(s.Limit)).
		Offset(int(s.Limit * s.Page)).
		Find(&men).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	for _, m := range men {
		results = append(results, m.ID)
	}

	c.JSON(http.StatusOK, results)
}

type Man struct {
	HeadImgUrl *string   `json:"headImgUrl"`
	ManInfos   []ManInfo `json:"manInfos" binding:"required"`
}

type ManInfo struct {
	Language string  `json:"language" binding:"required"`
	Name     string  `json:"name" binding:"required"`
	Comment  *string `json:"comment"`
}

type SearchParam struct {
	Limit uint `form:"limit"`
	Page  uint `form:"page"`
}
