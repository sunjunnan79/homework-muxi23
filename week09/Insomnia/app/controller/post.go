package controller

import (
	request2 "Insomnia/app/api/request"
	. "Insomnia/app/api/response"
	. "Insomnia/app/infrastructure/helper"
	"Insomnia/app/infrastructure/kafka"
	"Insomnia/app/infrastructure/redis"
	"Insomnia/app/models"
	"Insomnia/app/service"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"sync"
)

var lock1 sync.Mutex

type Post struct{}

var postService *service.PostService

// CreatePost 创建回复
// @Summary 用户创建回复接口
// @Description 用户创建回复接口
// @Tags Post
// @Accept json
// @Produce json
// @Param tUuid query string true "帖子id"
// @Param body query string true "回复内容"
// @Param Authorization header string true "jwt验证"
// @Success 200 {object} GetPostResponse "回复创建成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 500 {object} ErrorResponse "内部错误"
// @Router /api/v1/post/createPost [post]
func (p *Post) CreatePost(c *gin.Context) {
	//定义一个创建帖子请求的结构体
	req := &request2.CreatePostReq{}

	//使用ShouldBind去解析获得的结构体,蛙趣真清晰啊
	if err := c.ShouldBind(req); err != nil {
		FailMsg(c, fmt.Sprintf("params invalid error: %v", err))
		return
	}

	Uuid, _ := c.Get("Uuid")
	uuid := Uuid.(string)
	post, err := postService.CreatePost(uuid, *req)
	if err != nil {
		FailMsgData(c, fmt.Sprintf("创建回复失败: %v", err), models.Post{})
		return
	}

	rsp := GetPostResponse{
		CreatedAt: post.CreatedAt.Format("2006-01-02 15:04"),
		TUuid:     post.TUuid,
		UuId:      post.Uuid,
		PUuid:     post.PUuid,
		Likes:     post.Likes,
		Body:      post.Body,
	}

	err = redis.DelResp(post.TUuid)
	if err != nil {
		Danger(err, "删除对应post缓存数据失败")
	}

	OkMsgData(c, "创建回复成功", rsp)
	return
}

// ReadPost 读取回复
// @Summary 用户读取回复接口
// @Description 用户读取回复接口
// @Tags Post
// @Accept json
// @Produce json
// @Param pUuid query string true "回复id"
// @Param Authorization header string true "jwt验证"
// @Success 200 {object} GetPostResponse "回复读取成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 500 {object} ErrorResponse "内部错误"
// @Router /api/v1/post/readPost [post]
func (p *Post) ReadPost(c *gin.Context) {
	//定义一个寻找到帖子请求的结构体
	req := &request2.FindPostReq{}

	//使用ShouldBind去解析获得的结构体,蛙趣真清晰啊
	if err := c.ShouldBind(req); err != nil {
		FailMsg(c, fmt.Sprintf("params invalid error: %v", err))
		return
	}

	//缓存未命中的话就从数据库中拿取
	post, err := postService.ReadPost(req.PUuid)
	if err != nil {
		FailMsgData(c, fmt.Sprintf("获取回复失败: %v", err), models.Post{})
		return
	}

	//获取点赞状态
	Uuid, _ := c.Get("Uuid")
	uuid := Uuid.(string)
	exist, err := models.CheckLike(post.PUuid, uuid)
	if err != nil {
		return
	}

	rsp := GetPostResponse{
		CreatedAt: post.CreatedAt.Format("2006-01-02 15:04"),
		TUuid:     post.TUuid,
		UuId:      post.Uuid,
		PUuid:     post.PUuid,
		Likes:     post.Likes,
		Body:      post.Body,
		Exist:     strconv.FormatBool(exist),
	}
	OkMsgData(c, "获取回复成功", rsp)
	return
}

// DestroyPost 删除回复
// @Summary 用户删除回复接口
// @Description 用户删除回复接口
// @Tags Post
// @Accept json
// @Produce json
// @Param pUuid query string true "回复id"
// @Param Authorization header string true "jwt验证"
// @Success 200 {object} MessageResponse "删除成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 500 {object} ErrorResponse "内部错误"
// @Router /api/v1/post/destroyPost [post]
func (p *Post) DestroyPost(c *gin.Context) {
	//定义一个寻找到帖子请求的结构体
	req := &request2.FindPostReq{}

	//使用ShouldBind去解析获得的结构体,蛙趣真清晰啊
	if err := c.ShouldBind(req); err != nil {
		FailMsg(c, fmt.Sprintf("params invalid error: %v", err))
		return
	}

	err := postService.DestroyPost(req.PUuid)
	if err != nil {
		FailMsg(c, fmt.Sprintf("%v", err))
		return
	}

	OkMsg(c, "删除回复成功")
	return
}

// GetPosts 获取该帖子的回复
// @Summary 获取该帖子的回复接口
// @Description 获取该帖子的回复接口
// @Tags Post
// @Accept json
// @Produce json
// @Param tUuid query string true "帖子id"
// @Param Authorization header string true "jwt验证"
// @Success 200 {object} []GetPostResponse "回复数据"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 500 {object} ErrorResponse "内部错误"
// @Router /api/v1/post/getPosts [post]
func (p *Post) GetPosts(c *gin.Context) {
	//定义一个获取回复请求的结构体
	req := &request2.GetPostsReq{}

	//使用ShouldBind去解析获得的结构体,蛙趣真清晰啊
	if err := c.ShouldBind(req); err != nil {
		FailMsg(c, fmt.Sprintf("params invalid error: %v", err))
		return
	}
	//从缓存中优先获取
	cache, err := redis.GetResp("posts" + req.TUuid)
	if err == nil {
		var rsp []GetPostResponse
		//解析缓存为对应格式
		if err := json.Unmarshal(cache, &rsp); err == nil {
			go func() {
				for i, post := range rsp {
					//存储前10个帖子的回复到缓存
					if i < 10 {
						owner := "reposts" + post.PUuid
						exist, err := redis.ExistResp(owner)
						if exist != 1 && err == nil {
							lock1.Lock()
							// 创建 Kafka 实例
							topics := []string{"cache"}
							group := "cache-group"
							key := "reposts"
							k := kafka.NewKafka(topics, group, key)
							k.CreateCacheProducer(owner)
							lock1.Unlock()
						}
					}
				}
			}()
			OkMsgData(c, "从缓存中获取回复成功", rsp)
			return
		}
	}
	var rsp []GetPostResponse
	//获取回复
	posts, err := postService.GetPosts(req.TUuid)
	if err != nil {
		FailMsgData(c, fmt.Sprintf("%v", err), models.Post{})
		return
	}

	for i, post := range posts {
		p := GetPostResponse{
			CreatedAt: post.CreatedAt.Format("2006-01-02 15:04"),
			TUuid:     post.TUuid,
			UuId:      post.Uuid,
			PUuid:     post.PUuid,
			Likes:     post.Likes,
			Body:      post.Body,
		}
		//存储前10个帖子的回复到缓存
		if i < 10 {
			go func() {
				owner := "reposts" + post.PUuid
				exist, err := redis.ExistResp(owner)
				if exist != 1 && err == nil {
					lock1.Lock()
					// 创建 Kafka 实例
					topics := []string{"cache"}
					group := "cache-group"
					key := "reposts"
					k := kafka.NewKafka(topics, group, key)
					k.CreateCacheProducer(owner)
					lock1.Unlock()
				}
			}()
		}
		rsp = append(rsp, p)
	}

	OkMsgData(c, "获取回复成功", rsp)
	return
}

// LikePost 点赞/取消点赞回复
// @Summary 点赞/取消点赞回复接口
// @Description 用户点赞/取消点赞回复接口
// @Tags Post
// @Accept json
// @Produce json
// @Param uid query string true "这里对应的就是pUuid,但是方便你复制粘贴,帖子唯一标识"
// @Param Authorization header string true "jwt验证"
// @Success 200 {object} LikesResponse "点赞状态切换成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 500 {object} ErrorResponse "内部错误"
// @Router /api/v1/post/likePost [post]
func (p *Post) LikePost(c *gin.Context) {
	//定义一个创建帖子请求的结构体
	req := &request2.LikesReq{}
	//使用ShouldBind去解析获得的结构体,蛙趣真清晰啊
	if err := c.ShouldBind(req); err != nil {
		FailMsg(c, fmt.Sprintf("params invalid error: %v", err))
		return
	}

	Uuid, _ := c.Get("Uuid")
	uuid := Uuid.(string)

	//切换点赞状态
	exist, err := postService.LikePosts(req.Uid, uuid)
	if err != nil {
		FailMsgData(c, fmt.Sprintf("点赞切换失败:%v", err), LikesResponse{})
		return
	}

	OkMsgData(c, fmt.Sprintf("点赞状态切换成功!"), LikesResponse{Exist: strconv.FormatBool(exist)})
	return
}
