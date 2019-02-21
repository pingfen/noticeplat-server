package wechat

import (
	"net/http"
	"os"

	"fmt"
	"github.com/pingfen/noticeplat-server/pkg/msg"
	"gopkg.in/chanxuehong/wechat.v2/mp/core"
	"gopkg.in/chanxuehong/wechat.v2/mp/message/template"
)

var (
	client *core.Client
)

func init() {
	appid := os.Getenv("WX_SRV_APPID")
	if appid == "" {
		panic("WX_SRV_APPID was required")
	}

	app_secret := os.Getenv("WX_SRV_APPSECRET")
	if app_secret == "" {
		panic("WX_SRV_APPSECRET was required")
	}

	http_client := http.DefaultClient
	ats := core.NewDefaultAccessTokenServer(
		appid,
		app_secret,
		http_client,
	)
	client = core.NewClient(ats, http_client)
}

func SendTemplateMsg(srvOpenId string, todo *msg.Todo) error {
	msg := &template.TemplateMessage2{
		ToUser:     srvOpenId, // "o3HH903-YBeNkhLmdomElPrmxLZI"
		TemplateId: "Vkj5nukhBnbzuxmHA4ggf6y_vYhwbKjEyTIz_n1ygAU",
		MiniProgram: &template.MiniProgram{
			//AppId:    "wxb0168e8389c0e56f", //通知台
			AppId:    "wx9d0569208850e892", //实时公交车
			PagePath: "/pages/index/index",
		},
		Data: NewTemplData([]TemplDataItem{
			{Value: fmt.Sprintf("【%s】%s", msg.MsgTypeName(todo.MsgType), todo.Title), Color: "#dc143c"},
			{Value: todo.ID},
			{Value: todo.Project},
			{Value: msg.LevelName(todo.Level)},
			{Value: msg.MsgTypeName(todo.MsgType)},
			{Value: todo.Content},
			{Value: "点击编辑或查看详情!", Color: "#173177"},
		}),
	}
	_, err := template.Send(client, msg)
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
