package service

import (
	. "Insomnia/app/api/request"
	"Insomnia/app/models"
	"fmt"
)

type RePostService struct{}

func (r *RePostService) CreateRePost(uuid string, crp CreateRePostReq) (models.RePost, error) {
	// 开始事务
	tx := models.Db.Begin()
	defer tx.Commit()

	repost, err := models.CreateRePost(uuid, crp)
	if err != nil {
		// 如果出错，回滚事务
		tx.Rollback()
		return models.RePost{}, err
	}

	err = models.UpThreadNumbers(crp.TUuid)
	if err != nil {
		// 如果出错，回滚事务
		tx.Rollback()
		return models.RePost{}, err
	}

	author, err := models.AuthorByPUUID(crp.PUuid)
	if err != nil {
		// 如果出错，回滚事务
		tx.Rollback()
		return models.RePost{}, err
	}

	if uuid != author {
		pm := models.PostMessage{
			TUuid: repost.TUuid,
			Uuid:  author,
			PUuid: repost.PUuid,
			Body:  repost.Body,
			Check: false,
		}
		err = models.CreatePostMessage(pm)
		if err != nil {
			// 如果出错，回滚事务
			tx.Rollback()
			return models.RePost{}, err
		}
	}
	return repost, nil
}

func (r *RePostService) ReadRePost(rUuid string) (models.RePost, error) {
	repost, err := models.RePostByRUUID(rUuid)
	if err != nil {
		return models.RePost{}, err
	}
	return repost, nil
}

func (r *RePostService) DestroyRePost(rUuid string) error {
	// 开始事务
	tx := models.Db.Begin()
	defer tx.Commit()

	repost, err := models.RePostByRUUID(rUuid)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = models.DestroyRePost(rUuid)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = models.DownThreadNumbers(repost.TUuid)
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (r *RePostService) GetRePosts(pUuid string) ([]models.RePost, error) {
	reposts, err := models.RePosts(pUuid)
	if err != nil {
		return []models.RePost{}, fmt.Errorf("无法获取re回复:%v", err)
	}
	return reposts, nil
}

func (r *RePostService) LikeRePosts(rUuid string, uuid string) (exist bool, err error) {
	// 开始事务
	tx := models.Db.Begin()
	defer tx.Commit()
	exist, err = models.ChangeLike(rUuid, uuid)
	//检查是否出错
	if err != nil {
		// 如果出错，回滚事务
		tx.Rollback()
		return
	}
	//根据改变后的点赞类型自动增减
	err = models.UpRePostLikesData(rUuid, exist)
	if err != nil {
		// 如果出错，回滚事务
		tx.Rollback()
		return
	}
	RePost, err := models.RePostByRUUID(rUuid)
	if err != nil {
		// 如果出错，回滚事务
		tx.Rollback()
		return
	}
	lm := models.LikeMessage{
		TUuid: RePost.TUuid,
		Uuid:  RePost.Uuid,
		PUuid: RePost.PUuid,
		RUuid: RePost.RUuid,
		Body:  RePost.Body,
		Check: false,
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
	if RePost.Uuid != uuid {
		err = models.CreateLikeMessage(lm)
		if err != nil {
			// 如果出错，回滚事务
			tx.Rollback()
			return
		}
	}
	return
}
