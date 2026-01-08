package models

import "time"

type Community struct {
	ID   int    `db:"community_id" json:"id,string"`
	Name string `db:"community_name" json:"name"`
}

type CommunityDetail struct {
	ID           int       `db:"community_id" json:"id,string"`
	Name         string    `db:"community_name" json:"name"`
	Introduction string    `db:"introduction" json:"introduction,omitempty"`
	CreateTime   time.Time `db:"create_time" json:"create_time"`
}
