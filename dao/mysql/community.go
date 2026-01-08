package mysql

import (
	"bluebell/models"
	"database/sql"
	"errors"
	"fmt"
)

func GetCommunityList() ([]models.Community, error) {
	const sqlStr = `select community_id,community_name from community`
	list := make([]models.Community, 0)
	if err := db.Select(&list, sqlStr); err != nil {
		return nil, fmt.Errorf("dao: GetCommunityList: %w", err)
	}
	return list, nil
}

func GetCommunityDetailById(id int64) (*models.CommunityDetail, error) {
	sqlStr := `select community_id,community_name,introduction,create_time from community where community_id = ?`
	detail := new(models.CommunityDetail)
	if err := db.Get(detail, sqlStr, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("dao: GetCommunityDetailById: %w", ErrNotFound)
		}
		return nil, fmt.Errorf("dao: GetCommunityDetailById: %w", err)
	}
	return detail, nil
}
