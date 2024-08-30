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

	if len(talk.TalkContents) == 0 {
		c.JSON(http.StatusNotFound, gin.H{})

		return
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

func NewTalk(c *gin.Context) {
	var talk Talk
	id, err := util.ParseUint(c.Param("id"))

	if err := c.ShouldBind(&talk); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	r := convertToModel(&talk)

	if err := common.
		GetDB().
		Model(&model.GreatMan{
			Model: gorm.Model{
				ID: id,
			},
		}).
		Session(&gorm.Session{FullSaveAssociations: true}).
		Association("TalkRecords").
		Append(&r); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{})
			return
		}

		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, r.ID)
}

func UpdateTalk(c *gin.Context) {
	var talk Talk
	id, err := util.ParseUint(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	if err := c.ShouldBind(&talk); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	t := convertToModel(&talk)
	t.ID = id
	if err := common.
		GetDB().
		Session(&gorm.Session{FullSaveAssociations: true}).
		Updates(&t).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func DeleteTalk(c *gin.Context) {
	id, err := util.ParseUint(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	if err := common.GetDB().Delete(&model.TalkRecord{}, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func convertToModel(talk *Talk) model.TalkRecord {
	var tcs []model.TalkContent
	var as []model.Attachment

	for _, t := range talk.TalkContents {
		tcs = append(tcs, model.TalkContent{
			Language:    t.Language,
			MainBody:    t.MainBody,
			Title:       t.Title,
			Interviewer: t.Interviewer,
			Source:      t.Source,
		})
	}

	for _, a := range talk.Attachments {
		as = append(as, model.Attachment{
			Type:  model.AttachmentType(a.Type),
			Value: a.Value,
		})
	}

	return model.TalkRecord{
		Type:         talk.Type,
		TalkContents: tcs,
		Attachments:  as,
	}
}

type Talk struct {
	Type         string `json:"type" binding:"required"`
	TalkContents []struct {
		Language    string  `json:"language" binding:"required"`
		Title       string  `json:"title" binding:"required"`
		MainBody    string  `json:"mainBody" binding:"required"`
		Interviewer *string `json:"interviewer"`
		Source      *string `json:"source"`
	} `json:"talkContents"`
	Attachments []AttachmentData `json:"attachments"`
}

type AttachmentData struct {
	Type  string `json:"type" binding:"required"`
	Value string `json:"value" binding:"required"`
}

type SimpleTalkInfo struct {
	TalkRecordId uint      `json:"id" binding:"required"`
	Title        string    `json:"title" binding:"required"`
	CreatedAt    time.Time `json:"date" binding:"required"`
}
