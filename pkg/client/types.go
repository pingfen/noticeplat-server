package client

import "time"

const (
	CLIENT_KEYPREFIX = "/client/data"
	BYUSER_KEYPREIX  = "/client/.index/byuser"
)

type Client struct {
	Id   string `json:"id"`
	User string `json:"user"`

	CreateTime time.Time `json:"createTime"`
	UpdateTime time.Time `json:"updateTime"`
}
