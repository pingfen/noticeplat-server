package models

type Group struct {
	Id         string   `json:"id"`
	OpenGid    string   `json:"openGid"`
	Name       string   `json:"name"`
	HeadimgUrl string   `json:"headimgurl"`
	Managers   []string `json:"managers"`
	Secret     string   `json:"secret,omitempty"`
}

func (g *Group) IsManager(openid string) bool {
	for _, uid := range g.Managers {
		if uid == openid {
			return true
		}
	}

	return false
}
