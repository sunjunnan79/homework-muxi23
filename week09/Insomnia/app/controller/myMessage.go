package controller

import (
	. "Insomnia/app/api/request"
	. "Insomnia/app/api/response"
	. "Insomnia/app/infrastructure/helper"
	"Insomnia/app/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

type MyMessage struct{}

var myMessageService *service.MyMessageService

// GetPostMessage 获取用户的评论消息
// @Summary 获取用户的评论消息接口
// @Description 获取用户的评论消息接口
// @Tags MyMessage
// @Accept json
// @Produce json
// @Param Authorization header string true "jwt验证"
// @Success 200 {object} MessageResponse "获取消息成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 500 {object} ErrorResponse "内部错误"
// @Router /api/v1/myMessage/getPostMessage [post]
func (m *MyMessage) GetPostMessage(c *gin.Context) {
	Uuid, _ := c.Get("Uuid")
	uuid := Uuid.(string)
	err, pm := myMessageService.GetPostMessage(uuid)
	if err != nil {
		FailMsgData(c, fmt.Sprintf("获取评论信息失败: %v", err), nil)
		return
	}
	var rsp []MessageResponse
	for _, message := range pm {
		messages := MessageResponse{
			CreatedAt: message.CreatedAt.Format("2006-01-02 15:04"),
			TUuid:     message.TUuid,
			Title:     message.Title,
			UuId:      message.Uuid,
			Body:      message.Body,
		}
		rsp = append(rsp, messages)
	}
	OkMsgData(c, "获取评论信息成功", rsp)
}

// CheckPostMessage 查看用户的评论消息
// @Summary 查看用户的评论消息接口
// @Description 查看用户的评论消息接口
// @Tags MyMessage
// @Accept json
// @Produce json
// @Param id query string true "返回给你的post消息的id"
// @Param Authorization header string true "jwt验证"
// @Success 200 {object} MessageResponse "查看消息成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 500 {object} ErrorResponse "内部错误"
// @Router /api/v1/myMessage/checkPostMessage [post]
func (m *MyMessage) CheckPostMessage(c *gin.Context) {
	//定义一个Login请求类型的结构体
	req := &CheckPostMessageReq{}

	//使用ShouldBind去解析获得的结构体,蛙趣真清晰啊
	if err := c.ShouldBind(req); err != nil {
		Danger(err, "解析请求结构体失败")
		FailMsgData(c, fmt.Sprintf("params invalid error: %v", err), nil)
		return
	}
	err := myMessageService.CheckPostMessage(*req)
	if err != nil {
		return
	}
	OkMsg(c, "检查消息成功")
}

// GetLikeMessage 获取用户的点赞消息
// @Summary 获取用户的点赞消息接口
// @Description 获取用户的点赞消息接口
// @Tags MyMessage
// @Accept json
// @Produce json
// @Param Authorization header string true "jwt验证"
// @Success 200 {object} MessageResponse "获取帖子成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 500 {object} ErrorResponse "内部错误"
// @Router /api/v1/myMessage/getLikeMessage [post]
func (m *MyMessage) GetLikeMessage(c *gin.Context) {
	Uuid, _ := c.Get("Uuid")
	uuid := Uuid.(string)
	err, lm := myMessageService.GetLikeMessage(uuid)
	if err != nil {
		Danger(err, "获取点赞消息失败")
		FailMsgData(c, fmt.Sprintf("获取点赞消息失败: %v", err), nil)
		return
	}
	var rsp []MessageResponse
	for _, message := range lm {
		messages := MessageResponse{
			Id:        strconv.Itoa(int(message.ID)),
			CreatedAt: message.CreatedAt.Format("2006-01-02 15:04"),
			TUuid:     message.TUuid,
			Title:     message.Title,
			UuId:      message.Uuid,
			Body:      message.Body,
			Check:     message.Check,
		}
		rsp = append(rsp, messages)
	}
	OkMsgData(c, "获取点赞消息成功", rsp)
}

// CheckLikeMessage 查看用户的点赞消息
// @Summary 查看用户的点赞消息接口
// @Description 查看用户的点赞消息接口
// @Tags MyMessage
// @Accept json
// @Produce json
// @Param id query string true "返回给你的like消息的id"
// @Param Authorization header string true "jwt验证"
// @Success 200 {object} MessageResponse "查看消息成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 500 {object} ErrorResponse "内部错误"
// @Router /api/v1/myMessage/checkLikeMessage [post]
func (m *MyMessage) CheckLikeMessage(c *gin.Context) {
	//定义一个Login请求类型的结构体
	req := &CheckLikeMessageReq{}

	//使用ShouldBind去解析获得的结构体,蛙趣真清晰啊
	if err := c.ShouldBind(req); err != nil {
		Danger(err, "解析请求结构体失败")
		FailMsgData(c, fmt.Sprintf("删除对应的消息失败: %v", err), nil)
		return
	}
	err := myMessageService.CheckLikeMessage(*req)
	if err != nil {
		Danger(err, "删除对应的消息失败")
		FailMsgData(c, fmt.Sprintf("删除对应的消息失败: %v", err), nil)
		return
	}
	OkMsg(c, "检查消息成功")
}
