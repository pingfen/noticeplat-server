package models

import (
	"encoding/base64"
	"net/http"
)

type Captcha struct {
	IsBase64    bool           `json:"isBase64"`
	Data        string         `json:"data"`
	Cookies     []*http.Cookie `json:"cookies"`
	ContentType string         `json:"contentType"`
}

func (c *Captcha) Base64() *Captcha {
	c_new := *c
	if c.IsBase64 {
		return &c_new
	}

	c_new.IsBase64 = true
	c_new.Data = base64.StdEncoding.EncodeToString([]byte(c.Data))

	return &c_new
}

type ReserveRequest struct {
	Contact   string     `json:"contact"`
	MaskStore *MaskStore `json:"maskStore"`
	Captcha   *Captcha   `json:"captcha"`
}
