package mysql

import (
	"bluebell/models"
	"database/sql"
	"errors"
)

func CreatePost(p *models.Post) error {
	sqlStr := `INSERT INTO post(post_id,title,content,author_id,community_id) VALUES (?,?,?,?,?)`
	_, err := db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID)
	return err
}

func GetPostById(id int64) (*models.Post, error) {
	sqlStr := `SELECT post_id,title,content,author_id,community_id,create_time FROM post WHERE post_id = ?`
	post := new(models.Post)
	err := db.Get(post, sqlStr, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return post, nil
}

func GetPostTotalCount() (count int64, err error) {
	sqlStr := `select count(*) from post`
	err = db.Get(&count, sqlStr)
	return
}

func GetPostList(limit int64, offset int64) (list []*models.Post, err error) {
	sqlStr := `SELECT post_id, title, content, author_id, community_id,status,create_time
		FROM post
		ORDER BY create_time DESC
		LIMIT ? OFFSET ?`
	err = db.Select(&list, sqlStr, limit, offset)
	return
}
