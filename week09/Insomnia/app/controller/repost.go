package controller

import (
	. "Insomnia/app/api/request"
	. "Insomnia/app/api/response"
	. "Insomnia/app/infrastructure/helper"
	"Insomnia/app/infrastructure/redis"
	"Insomnia/app/models"
	"Insomnia/app/service"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

type RePost struct{}

var repostService *service.RePostService

// CreateRePost 创建re回复
// @Summary 用户创建re回复
// @Description 用户创建re回复接口
// @Tags RePost
// @Accept json
// @Produce json
// @Param tUuid query string true "帖子id"
// @Param pUuid query string true "回复id"
// @Param body query string true "re回复内容"
// @Param Authorization header string true "jwt验证"
// @Success 200 {object} GetRePostResponse "re回复创建成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 500 {object} ErrorResponse "内部错误"
// @Router /api/v1/rePost/createRePost [post]
func (r *RePost) CreateRePost(c *gin.Context) {
	//定义一个创建帖子请求的结构体
	req := &CreateRePostReq{}

	//使用ShouldBind去解析获得的结构体,蛙趣真清晰啊
	if err := c.ShouldBind(req); err != nil {
		FailMsg(c, fmt.Sprintf("params invalid error: %v", err))
		return
	}

	Uuid, _ := c.Get("Uuid")
	uuid := Uuid.(string)
	repost, err := repostService.CreateRePost(uuid, *req)
	if err != nil {
		FailMsgData(c, fmt.Sprintf("创建re回复失败: %v", err), models.RePost{})
		return
	}

	rsp := GetRePostResponse{
		CreatedAt: repost.CreatedAt.Format("2006-01-02 15:04"),
		TUuid:     repost.TUuid,
		UuId:      repost.Uuid,
		PUuid:     repost.PUuid,
		Likes:     repost.Likes,
		Body:      repost.Body,
		RUuid:     repost.RUuid,
	}
	err = redis.DelResp(repost.TUuid)
	if err != nil {
		Danger(err, "删除对应post缓存数据失败")
	}
	OkMsgData(c, "创建回复成功", rsp)
	return
}

// ReadRePost 创建re回复
// @Summary 用户创建re回复
// @Description 用户创建re回复接口
// @Tags RePost
// @Accept json
// @Produce json
// @Param rUuid query string true "re回复id"
// @Param Authorization header string true "jwt验证"
// @Success 200 {object} GetRePostResponse "re回复获取成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 500 {object} ErrorResponse "内部错误"
// @Router /api/v1/rePost/readRePost [post]
func (r *RePost) ReadRePost(c *gin.Context) {
	//定义一个寻找到帖子请求的结构体
	req := &FindRePostReq{}

	//使用ShouldBind去解析获得的结构体,蛙趣真清晰啊
	if err := c.ShouldBind(req); err != nil {
		FailMsg(c, fmt.Sprintf("params invalid error: %v", err))
		return
	}

	//从缓存中优先获取
	cache, err := redis.GetResp(req.RUuid)
	if err == nil {
		rsp := GetRePostResponse{}
		//解析缓存为对应格式
		if err := json.Unmarshal(cache, &rsp); err == nil {
			//获取点赞状态
			Uuid, _ := c.Get("Uuid")
			uuid := Uuid.(string)
			exist, err := models.CheckLike(rsp.RUuid, uuid)
			if err == nil {
				rsp.Exist = strconv.FormatBool(exist)
				OkMsgData(c, "从缓存中获取回复成功", rsp)
				return
			}
		}
	}

	//未命中正常获取
	repost, err := repostService.ReadRePost(req.RUuid)
	if err != nil {
		FailMsgData(c, fmt.Sprintf("获取re回复失败: %v", err), models.Post{})
		return
	}

	//获取点赞状态
	Uuid, _ := c.Get("Uuid")
	uuid := Uuid.(string)
	exist, err := models.CheckLike(repost.RUuid, uuid)
	if err != nil {
		return
	}

	rsp := GetRePostResponse{
		CreatedAt: repost.CreatedAt.Format("2006-01-02 15:04"),
		TUuid:     repost.TUuid,
		UuId:      repost.Uuid,
		PUuid:     repost.PUuid,
		RUuid:     repost.RUuid,
		Likes:     repost.Likes,
		Body:      repost.Body,
		Exist:     strconv.FormatBool(exist),
	}

	OkMsgData(c, "获取re回复成功", rsp)
	return
}

// DestroyRePost 删除re回复
// @Summary 用户删除re回复
// @Description 用户删除re回复接口
// @Tags RePost
// @Accept json
// @Produce json
// @Param rUuid query string true "re回复id"
// @Param Authorization header string true "jwt验证"
// @Success 200 {object} MessageResponse "re回复删除成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 500 {object} ErrorResponse "内部错误"
// @Router /api/v1/rePost/destroyRePost [post]
func (r *RePost) DestroyRePost(c *gin.Context) {
	//定义一个寻找到帖子请求的结构体
	req := &FindRePostReq{}

	//使用ShouldBind去解析获得的结构体,蛙趣真清晰啊
	if err := c.ShouldBind(req); err != nil {
		FailMsg(c, fmt.Sprintf("params invalid error: %v", err))
		return
	}

	err := repostService.DestroyRePost(req.RUuid)
	if err != nil {
		FailMsg(c, fmt.Sprintf("%v", err))
		return
	}

	OkMsg(c, "删除re回复成功")
	return
}

// GetRePosts 获取指定回复下的re回复
// @Summary 用户获取指定回复下的re回复
// @Description 用户获取指定回复下的re回复接口
// @Tags RePost
// @Accept json
// @Produce json
// @Param pUuid query string true "回复id"
// @Param Authorization header string true "jwt验证"
// @Success 200 {object} GetRePostResponse "re回复获取成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 500 {object} ErrorResponse "内部错误"
// @Router /api/v1/rePost/getRePosts [post]
func (r *RePost) GetRePosts(c *gin.Context) {

	//定义一个获取回复请求的结构体
	req := &GetRePostsReq{}

	//使用ShouldBind去解析获得的结构体,蛙趣真清晰啊
	if err := c.ShouldBind(req); err != nil {
		FailMsg(c, fmt.Sprintf("params invalid error: %v", err))
		return
	}

	//获取回复
	reposts, err := repostService.GetRePosts(req.PUuid)
	if err != nil {
		FailMsgData(c, fmt.Sprintf("%v", err), models.RePost{})
		return
	}

	var rsp []GetRePostResponse
	for _, repost := range reposts {
		rp := GetRePostResponse{
			CreatedAt: repost.CreatedAt.Format("2006-01-02 15:04"),
			TUuid:     repost.TUuid,
			UuId:      repost.Uuid,
			PUuid:     repost.PUuid,
			Likes:     repost.Likes,
			Body:      repost.Body,
			RUuid:     repost.RUuid,
		}
		rsp = append(rsp, rp)
	}

	OkMsgData(c, "获取re回复成功", rsp)
	return
}

// LikeRePost 获取指定回复下的re回复
// @Summary 用户获取指定回复下的re回复
// @Description 用户获取指定回复下的re回复接口
// @Tags RePost
// @Accept json
// @Produce json
// @Param uid query string true "这里对应的就是rUuid,但是方便你复制粘贴,帖子唯一标识"
// @Param Authorization header string true "jwt验证"
// @Success 200 {object} LikesResponse "点赞状态切换成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 500 {object} ErrorResponse "内部错误"
// @Router /api/v1/rePost/likeRePost [post]
func (r *RePost) LikeRePost(c *gin.Context) {
	//定义一个创建帖子请求的结构体
	req := &LikesReq{}
	//使用ShouldBind去解析获得的结构体,蛙趣真清晰啊
	if err := c.ShouldBind(req); err != nil {
		FailMsg(c, fmt.Sprintf("params invalid error: %v", err))
		return
	}

	Uuid, _ := c.Get("Uuid")
	uuid := Uuid.(string)

	//切换点赞状态
	exist, err := repostService.LikeRePosts(req.Uid, uuid)
	if err != nil {
		FailMsgData(c, fmt.Sprintf("点赞切换失败:%v", err), LikesResponse{})
		return
	}

	OkMsgData(c, fmt.Sprintf("点赞状态切换成功!"), LikesResponse{Exist: strconv.FormatBool(exist)})
	return
}
