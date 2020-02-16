package wechat

import (
	"fmt"
	"github.com/pingfen/notify/pkg/models"
	"gopkg.in/chanxuehong/wechat.v2/mp/message/template"
	"time"
)

func SendTemplateMsg(srvOpenId string, t *models.Todo) error {
	msg := &template.TemplateMessage2{
		ToUser:     srvOpenId, // "o3HH903-YBeNkhLmdomElPrmxLZI"
		TemplateId: "1C8OQl70lFxEbcrGlyrHETli4nMrTLJWXSPKAJKNWgw",
		//MiniProgram: &template.MiniProgram{
		//	//AppId:    "wxb0168e8389c0e56f", //通知台
		//	//AppId:    "wx9d0569208850e892", //实时公交车
		//	AppId:    "wx4f1d0a8b20c8ab32", // 零待办
		//	PagePath: "/pages/index/index",
		//},
		Data: NewTemplData([]TemplDataItem{
			{Value: fmt.Sprintf("【%s】%s", t.LevelName(), t.Subject), Color: "#dc143c"},
			{Value: t.Labels["project"]},
			{Value: t.TypeName()},
			{Value: t.Content},
			{Value: time.Unix(t.StartTime, 0).Format("2006-01-02 15:04:06")},
			{Value: "点击编辑或查看详情!", Color: "#173177"},
		}),
	}
	_, err := template.Send(srvClient, msg)
	return err
}

type TemplDataItem struct {
	Value string `json:"value"`
	Color string `json:"color,omitempty"`
}

func NewTemplData(items []TemplDataItem) map[string]TemplDataItem {
	l := len(items)
	ret := make(map[string]TemplDataItem, len(items))
	for i, item := range items {
		switch i {
		case 0:
			ret["first"] = item
		case l - 1:
			ret["remark"] = item
		default:
			ret[fmt.Sprintf("keyword%d", i)] = item
		}
	}

	//fmt.Printf("%+v", ret)
	return ret
}
