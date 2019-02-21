package group

const (
	GROUP_KEYPREFIX        = "/group/data"
	GROUPMEMBERS_KEYPREFIX = "/group/members"
)

type GroupBase struct {
	OpenGId string `json:"openGId"`
	Name    string `json:"name"`
	Avatar  string `json:"avatar"`
}

type Group struct {
	*GroupBase
	Members []Member `json:"members,omitempty"`
}

type Member struct {
	OpenId    string `json:"openId"`
	TimeStamp int64  `json:"timestamp"`
}
