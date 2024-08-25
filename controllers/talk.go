package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"greatmanbackend/common"
	"greatmanbackend/model"
	"greatmanbackend/util"
	"net/http"
	"time"
)

func GetTalks(c *gin.Context) {
	var s SearchParam
	greatManId, err := util.ParseUint(c.Param("id"))
	lang := c.Query("lang")
	talkType := c.Query("type")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	if err := c.ShouldBind(&s); err != nil || s.Limit > 20 {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	var talkRecords []model.TalkRecord

	if err := common.GetDB().
		Order("created_at DESC").
		Limit(int(s.Limit)).
		Offset(int(s.Limit * s.Page)).
		Select("ID").
		Where(&model.TalkRecord{
			GreatManId: greatManId,
			Type:       talkType,
		}).
		Find(&talkRecords).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	var ids []uint
	for _, r := range talkRecords {
		ids = append(ids, r.ID)
	}

	var talkInfos []SimpleTalkInfo

	if err := common.GetDB().
		Model(&model.TalkContent{}).
		Where(&model.TalkContent{
			Language: lang,
		}).
		Find(&talkInfos, ids).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusOK, talkInfos)
}

func GetTalkDetail(c *gin.Context) {
	talkId, err := util.ParseUint(c.Param("id"))
	lang := c.Query("lang")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	var talk model.TalkRecord

	if err := common.GetDB().
		Preload("TalkContents", &model.TalkContent{Language: lang}).
		Preload("Attachments").
		Take(&talk, talkId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	var attachments []AttachmentData

	for _, a := range talk.Attachments {
		attachments = append(attachments, AttachmentData{
			Type:  string(a.Type),
			Value: a.Value,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"title":       talk.TalkContents[0].Title,
		"content":     talk.TalkContents[0].MainBody,
		"interviewer": talk.TalkContents[0].Interviewer,
		"source":      talk.TalkContents[0].Source,
		"attachments": attachments,
		"date":        talk.CreatedAt,
	})
}

type AttachmentData struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type SimpleTalkInfo struct {
	TalkRecordId uint      `json:"id"`
	Title        string    `json:"title"`
	CreatedAt    time.Time `json:"date"`
}
