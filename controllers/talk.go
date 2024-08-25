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
	greatManId, err1 := util.ParseUint(c.Param("id"))
	lang := c.Query("lang")
	limit, err2 := util.ParseUint(c.Query("limit"))
	page, err3 := util.ParseUint(c.Query("page"))
	talkType := c.Query("type")

	if err1 != nil || err2 != nil || err3 != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	if limit >= 20 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "limit too large",
		})
		return
	}

	skip := limit * page

	var talkRecords []model.TalkRecord

	if err := common.GetDB().
		Order("created_at DESC").
		Limit(int(limit)).
		Offset(int(skip)).
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
