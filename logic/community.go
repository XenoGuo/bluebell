package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"errors"
)

// GetCommunityList 查所有的社区返回
func GetCommunityList() ([]models.Community, error) {
	data, err := mysql.GetCommunityList()
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetCommunityDetail 通过id获取社区详情
func GetCommunityDetail(id int64) (*models.CommunityDetail, error) {
	data, err := mysql.GetCommunityDetailById(id)
	if err != nil {
		if errors.Is(err, mysql.ErrNotFound) {
			return nil, ErrCommunityNotFound
		}
		return nil, err
	}
	return data, nil
}
