package service

import (
	. "Insomnia/app/api/request"
	"Insomnia/app/models"
)

type MyMessageService struct{}

// GetPostMessage 评论提醒
func (m *MyMessageService) GetPostMessage(uuid string) (err error, pm []models.PostMessage) {
	pm, err = models.PostMessageByUuId(uuid)
	if err != nil {
		return
	}
	return
}

// CheckPostMessage 检查评论提醒
func (m *MyMessageService) CheckPostMessage(dpmq CheckPostMessageReq) (err error) {
	err = models.CheckPostMessage(dpmq.Id)
	if err != nil {
		return
	}
	return
}

// GetLikeMessage 点赞提醒
func (m *MyMessageService) GetLikeMessage(uuid string) (err error, lm []models.LikeMessage) {
	lm, err = models.LikeMessageByUuId(uuid)
	if err != nil {
		return
	}
	return
}

// CheckLikeMessage 删除消息提醒
func (m *MyMessageService) CheckLikeMessage(dlmq CheckLikeMessageReq) (err error) {
	err = models.CheckLikeMessage(dlmq.Id)
	if err != nil {
		return
	}
	return
}
