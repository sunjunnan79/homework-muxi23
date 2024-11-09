package service

import (
	. "Insomnia/app/api/request"
	"Insomnia/app/models"
	"fmt"
)

type ThreadService struct{}

// CreateThread POST /thread/create
// 执行群组创建逻辑
func (t *ThreadService) CreateThread(uuid string, ct CreateThreadReq) (models.Thread, error) {
	thread, err := models.CreateThread(uuid, ct)
	if err != nil {
		return models.Thread{}, err
	}
	return thread, nil
}

// ReadThread GET /thread/read
// 通过uuid来获取帖子
func (t *ThreadService) ReadThread(tUuid string) (models.Thread, error) {
	thread, err := models.ThreadByTUUID(tUuid)
	if err != nil {
		return models.Thread{}, err
	}
	return thread, nil
}

// DestroyThread 通过tUuid来删除帖子(没写完,要把回复和回复的回复都删了)
func (t *ThreadService) DestroyThread(tUuid string) error {
	err := models.DestroyThread(tUuid)
	if err != nil {
		return err
	}
	return nil
}

// GetThreads 获取所有的帖子,以降序的方式
func (t *ThreadService) GetThreads(topic string) ([]models.Thread, error) {
	thread, err := models.Threads(topic)
	if err != nil {
		return []models.Thread{}, fmt.Errorf("无法获取帖子:%v", err)
	}
	return thread, nil
}

// GetMyThreads 获取我的所有的帖子,以降序的方式
func (t *ThreadService) GetMyThreads(UuId string) ([]models.Thread, error) {
	thread, err := models.ThreadByUuId(UuId)
	if err != nil {
		return []models.Thread{}, fmt.Errorf("无法获取帖子:%v", err)
	}
	return thread, nil
}

// LikeThreads 点赞或者取消点赞
func (t *ThreadService) LikeThreads(tUuid string, uuid string) (exist bool, err error) {
	// 开始事务
	tx := models.Db.Begin()
	defer tx.Commit()

	exist, err = models.ChangeLike(tUuid, uuid)
	//检查是否出错
	if err != nil {
		// 如果出错，回滚事务
		tx.Rollback()
		return
	}

	//根据改变后的点赞类型自动增减
	err = models.UpThreadLikesData(tUuid, exist)
	if err != nil {
		// 如果出错，回滚事务
		tx.Rollback()
		return
	}

	thread, err := models.ThreadByTUUID(tUuid)
	if err != nil {
		// 如果出错，回滚事务
		tx.Rollback()
		return
	}

	lm := models.LikeMessage{
		TUuid: thread.TUuid,
		Uuid:  thread.Uuid,
		Title: thread.Title,
		Body:  thread.Body,
		Check: false,
		Likes: 0,
	}
	//如果不存在则删除该LikeMessage
	if !exist {
		err = models.DeleteLikeMessage(lm)
		if err != nil {
			// 如果出错，回滚事务
			tx.Rollback()
			return false, err
		}
		return
	}
	//如果存在而且不为评论作者,则创建或者恢复该点赞消息
	if thread.Uuid != uuid {
		err = models.CreateLikeMessage(lm)
		if err != nil {
			// 如果出错，回滚事务
			tx.Rollback()
			return
		}
	}
	return
}
