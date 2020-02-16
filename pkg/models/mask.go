package models

import "encoding/json"

type MaskStore struct {
	Area    string `json:"area"`
	Address string `json:"address"`
	OrgCode string `json:"orgcode"`
	Name    string `json:"name"`
	ID      int    `json:"id"`
	Count   int    `json:"kznum"`
}

type MaskUserInfo struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	CardNo string `json:"cardNo"`
	Mobile string `json:"mobile"`

	ExpectStore *MaskStore `json:"expectStore"`
}

func (mui *MaskUserInfo) Bytes() []byte {
	ret, _ := json.Marshal(mui)
	return ret
}
