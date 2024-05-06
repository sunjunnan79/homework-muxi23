package controller

import (
	. "Insomnia/app/api/request"
	. "Insomnia/app/api/response"
	"Insomnia/app/infrastructure/kafka"
	"Insomnia/app/infrastructure/redis"
	"Insomnia/app/models"
	"Insomnia/app/service"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"sync"
)

type Thread struct{}

var lock sync.Mutex
var threadService *service.ThreadService

// CreateThread 创建帖子
// @Summary 用户创建帖子接口
// @Description 用户创建帖子接口
// @Tags Thread
// @Accept json
// @Produce json
// @Param title query string true "标题"
// @Param topic query string true "主题"
// @Param body query string true "内容"
// @Param images query []string true "图片链接列表"
// @Param Authorization header string true "jwt验证"
// @Success 200 {object} GetThreadResponse "帖子创建成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 500 {object} ErrorResponse "内部错误"
// @Router /api/v1/thread/createThread [post]
func (t *Thread) CreateThread(c *gin.Context) {
	//定义一个创建帖子请求的结构体
	req := &CreateThreadReq{}

	//使用ShouldBind去解析获得的结构体,蛙趣真清晰啊
	if err := c.ShouldBind(req); err != nil {
		FailMsg(c, fmt.Sprintf("无法解析的表单: %v", err))
		return
	}
	fmt.Println(req.Images)
	Uuid, _ := c.Get("Uuid")
	uuid := Uuid.(string)
	thread, err := threadService.CreateThread(uuid, *req)
	if err != nil {
		FailMsgData(c, fmt.Sprintf("创建帖子失败: %v", err), models.Thread{})
		return
	}
	// 反序列化 JSON 字符串并输出
	var retrievedSlice []string
	if err := json.Unmarshal([]byte(thread.Images), &retrievedSlice); err != nil {
		log.Fatal("反序列化数据失败:", err)
	}

	rsp := GetThreadResponse{
		CreatedAt: thread.CreatedAt.Format("2006-01-02 15:04"),
		TUuid:     thread.TUuid,
		Topic:     thread.Topic,
		Title:     thread.Title,
		UuId:      thread.Uuid,
		Likes:     thread.Likes,
		Body:      thread.Body,
		Number:    thread.Number,
		Images:    retrievedSlice,
	}
	OkMsgData(c, "创建帖子成功", rsp)
	return
}

// ReadThread 获取帖子详情
// @Summary 获取帖子详情接口
// @Description 获取指定帖子的详情信息接口
// @Tags Thread
// @Accept json
// @Produce json
// @Param tUuid query string true "帖子唯一标识"
// @Param Authorization header string true "jwt验证"
// @Success 200 {object} GetThreadResponse "获取帖子成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 500 {object} ErrorResponse "内部错误"
// @Router /api/v1/thread/readThread [post]
func (t *Thread) ReadThread(c *gin.Context) {
	//定义一个寻找到帖子请求的结构体
	req := &FindThreadReq{}

	//使用ShouldBind去解析获得的结构体,蛙趣真清晰啊
	if err := c.ShouldBind(req); err != nil {
		FailMsg(c, fmt.Sprintf("params invalid error: %v", err))
		return
	}

	thread, err := threadService.ReadThread(req.TUuid)
	if err != nil {
		FailMsgData(c, fmt.Sprintf("获取帖子失败: %v", err), models.Thread{})
		return
	}

	// 反序列化 JSON 字符串并输出
	var retrievedSlice []string
	if err := json.Unmarshal([]byte(thread.Images), &retrievedSlice); err != nil {
		log.Fatal("反序列化数据失败:", err)
	}

	//获取点赞状态
	Uuid, _ := c.Get("Uuid")
	uuid := Uuid.(string)
	exist, err := models.CheckLike(thread.TUuid, uuid)
	if err != nil {
		return
	}

	rsp := GetThreadResponse{
		CreatedAt: thread.CreatedAt.Format("2006-01-02 15:04"),
		TUuid:     thread.TUuid,
		Topic:     thread.Topic,
		Title:     thread.Title,
		UuId:      thread.Uuid,
		Likes:     thread.Likes,
		Body:      thread.Body,
		Number:    thread.Number,
		Images:    retrievedSlice,
		Exist:     strconv.FormatBool(exist),
	}

	OkMsgData(c, "获取帖子成功", rsp)
	return
}

// DestroyThread 删除帖子
// @Summary 删除帖子接口
// @Description 用户删除帖子接口
// @Tags Thread
// @Accept json
// @Produce json
// @Param tUuid query string true "帖子唯一标识"
// @Param Authorization header string true "jwt验证"
// @Success 200 {string} string "删除帖子成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 500 {object} ErrorResponse "内部错误"
// @Router /api/v1/thread/destroyThread [post]
func (t *Thread) DestroyThread(c *gin.Context) {
	//定义一个寻找到帖子请求的结构体
	req := &FindThreadReq{}

	//使用ShouldBind去解析获得的结构体,蛙趣真清晰啊
	if err := c.ShouldBind(req); err != nil {
		FailMsg(c, fmt.Sprintf("params invalid error: %v", err))
		return
	}

	err := threadService.DestroyThread(req.TUuid)
	if err != nil {
		FailMsg(c, fmt.Sprintf("%v", err))
		return
	}

	FailMsg(c, "删除帖子成功")
	return
}

// GetThreads 获取帖子列表
// @Summary 获取帖子列表接口
// @Description 获取指定主题的帖子列表接口
// @Tags Thread
// @Accept json
// @Produce json
// @Param topic query string true "主题"
// @Param Authorization header string true "jwt验证"
// @Success 200 {object} []GetThreadResponse "获取帖子成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 500 {object} ErrorResponse "内部错误"
// @Router /api/v1/thread/getThreads [post]
func (t *Thread) GetThreads(c *gin.Context) {
	//定义一个获取帖子请求的结构体
	req := &GetThreadsReq{}

	//使用ShouldBind去解析获得的结构体,蛙趣真清晰啊
	if err := c.ShouldBind(req); err != nil {
		FailMsg(c, fmt.Sprintf("params invalid error: %v", err))
		return
	}

	//获取帖子
	threads, err := threadService.GetThreads(req.Topic)
	if err != nil {
		FailMsgData(c, fmt.Sprintf("%v", err), models.Thread{})
		return
	}

	var rsp []GetThreadResponse
	for i, thread := range threads {
		threads := GetThreadResponse{
			CreatedAt: thread.CreatedAt.Format("2006-01-02 15:04"),
			TUuid:     thread.TUuid,
			Topic:     thread.Topic,
			Title:     thread.Title,
			UuId:      thread.Uuid,
			Likes:     thread.Likes,
			Body:      thread.Body,
			Number:    thread.Number,
		}

		//存储前10个帖子的回复到缓存
		if i < 10 {
			go func() {
				owner := thread.TUuid
				exist, err := redis.ExistResp("posts" + owner)
				if exist != 1 && err == nil {
					// 创建 Kafka 实例
					topics := []string{"cache"}
					group := "cache-group"
					key := "posts"
					k := kafka.NewKafka(topics, group, key)
					k.CreateCacheProducer(owner)
				}
			}()
		}
		rsp = append(rsp, threads)
	}

	OkMsgData(c, "获取帖子成功", rsp)
	return
}

// GetMyThreads 获取用户的帖子列表
// @Summary 获取用户的帖子列表接口
// @Description 获取指定用户的帖子列表接口
// @Tags Thread
// @Accept json
// @Produce json
// @Param Authorization header string true "jwt验证"
// @Success 200 {object} []GetThreadResponse "获取帖子成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 500 {object} ErrorResponse "内部错误"
// @Router /api/v1/thread/getMyThreads [post]
func (t *Thread) GetMyThreads(c *gin.Context) {

	Uuid, _ := c.Get("Uuid")
	uuid := Uuid.(string)
	//获取帖子
	threads, err := threadService.GetMyThreads(uuid)
	if err != nil {
		FailMsgData(c, fmt.Sprintf("%v", err), models.Thread{})
		return
	}

	var rsp []GetThreadResponse
	for _, thread := range threads {
		threads := GetThreadResponse{
			CreatedAt: thread.CreatedAt.Format("2006-01-02 15:04"),
			TUuid:     thread.TUuid,
			Topic:     thread.Topic,
			Title:     thread.Title,
			UuId:      thread.Uuid,
			Likes:     thread.Likes,
			Body:      thread.Body,
			Number:    thread.Number,
		}
		rsp = append(rsp, threads)
	}

	OkMsgData(c, "获取帖子成功", rsp)
	return
}

// LikeThread 点赞/取消点赞帖子
// @Summary 点赞/取消点赞帖子接口
// @Description 用户点赞或取消点赞帖子接口
// @Tags Thread
// @Accept json
// @Produce json
// @Param uid query string true "这里对应的就是tUuid,但是方便你复制粘贴,帖子唯一标识"
// @Param Authorization header string true "jwt验证"
// @Success 200 {object} LikesResponse "点赞状态切换成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 500 {object} ErrorResponse "内部错误"
// @Router /api/v1/thread/likeThread [post]
func (t *Thread) LikeThread(c *gin.Context) {
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
	exist, err := threadService.LikeThreads(req.Uid, uuid)
	if err != nil {
		FailMsgData(c, fmt.Sprintf("点赞切换失败:%v", err), LikesResponse{})
		return
	}

	OkMsgData(c, fmt.Sprintf("点赞状态切换成功!"), LikesResponse{Exist: strconv.FormatBool(exist)})
	return
}
