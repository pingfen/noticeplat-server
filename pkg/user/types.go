package user

import "time"

const (
	USER_KEYPREFIX     = "/user/data"
	APPOPENID_KEYPREIX = "/user/.index/app-openid"
	SRVOPENID_KEYPREIX = "/user/.index/srv-openid"
)

type User struct {
	Id        string `json:"id"`
	AppOpenId string `json:"appOpenId"`
	SrvOpenId string `json:"srvOpenId"`

	Name      string `json:"name"`
	Province  string `json:"province,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Email     string `json:"email,omitempty"`
	AvatarUrl string `json:"avatarUrl,omitempty"`

	CreateTime  time.Time `json:"createTime"`
	UpdateTime  time.Time `json:"updateTime"`
	SrvBindTime time.Time `json:"srvBindTime"`

	Groups []string `json:"groups"`
}

func (u *User) BindSrvOpenId(id string) {
	u.SrvOpenId = id
	u.SrvBindTime = time.Now()
}
