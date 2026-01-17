package mysql

import (
	"bluebell/models"
	"database/sql"
	"errors"
	"strings"

	"github.com/jmoiron/sqlx"
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

func GetPostsByIDs(ids []string) (list []*models.Post, err error) {
	sqlStr := `SELECT post_id, title, content, author_id, community_id,status,create_time
		FROM post
		WHERE post_id IN (?)
		ORDER BY FIND_IN_SET(post_id, ?)`
	// https://www.liwenzhou.com/posts/Go/sqlx/
	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	if err != nil {
		return nil, err
	}
	query = db.Rebind(query)
	err = db.Select(&list, query, args...) // !!!
	return
}
