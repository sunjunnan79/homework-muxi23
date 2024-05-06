package service

import (
	. "Insomnia/app/api/request"
	"Insomnia/app/models"
	"fmt"
)

type PostService struct{}

func (p *PostService) CreatePost(uuid string, cp CreatePostReq) (models.Post, error) {
	// 开始事务
	tx := models.Db.Begin()
	defer tx.Commit()
	Post, err := models.CreatePost(uuid, cp)
	if err != nil {
		// 如果出错，回滚事务
		tx.Rollback()
		return models.Post{}, err
	}

	err = models.UpThreadNumbers(Post.TUuid)
	if err != nil {
		// 如果出错，回滚事务
		tx.Rollback()
		return models.Post{}, err
	}

	author, err := models.AuthorByTUUID(cp.TUuid)
	if err != nil {
		// 如果出错，回滚事务
		tx.Rollback()
		return models.Post{}, err
	}

	if author != uuid {
		pm := models.PostMessage{
			TUuid: Post.TUuid,
			Uuid:  author,
			PUuid: Post.PUuid,
			Body:  Post.Body,
			Check: false,
		}
		err = models.CreatePostMessage(pm)
		if err != nil {
			// 如果出错，回滚事务
			tx.Rollback()
			return models.Post{}, err
		}
	}
	return Post, nil
}

func (p *PostService) ReadPost(pUuid string) (models.Post, error) {
	post, err := models.PostByPUUID(pUuid)
	if err != nil {
		return models.Post{}, err
	}
	return post, nil
}

func (p *PostService) DestroyPost(pUuid string) error {
	// 开始事务
	tx := models.Db.Begin()
	defer tx.Commit()

	post, err := models.PostByPUUID(pUuid)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = models.DestroyPost(pUuid)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = models.DownThreadNumbers(post.TUuid)
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (p *PostService) GetPosts(tUuid string) ([]models.Post, error) {
	posts, err := models.Posts(tUuid)
	if err != nil {
		return []models.Post{}, fmt.Errorf("无法获取回复:%v", err)
	}
	return posts, nil
}

func (p *PostService) LikePosts(pUuid string, uuid string) (exist bool, err error) {
	// 开始事务
	tx := models.Db.Begin()
	defer tx.Commit()
	exist, err = models.ChangeLike(pUuid, uuid)
	//检查是否出错
	if err != nil {
		// 如果出错，回滚事务
		tx.Rollback()
		return
	}
	//根据改变后的点赞类型自动增减
	err = models.UpPostLikesData(pUuid, exist)
	if err != nil {
		// 如果出错，回滚事务
		tx.Rollback()
		return
	}
	Post, err := models.PostByPUUID(pUuid)
	if err != nil {
		// 如果出错，回滚事务
		tx.Rollback()
		return
	}
	lm := models.LikeMessage{
		TUuid: Post.TUuid,
		Uuid:  Post.Uuid,
		PUuid: Post.PUuid,
		Body:  Post.Body,
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
	if Post.Uuid != uuid {
		err = models.CreateLikeMessage(lm)
		if err != nil {
			// 如果出错，回滚事务
			tx.Rollback()
			return
		}
	}
	return
}
