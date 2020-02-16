package e

type Err struct {
	Code int
	Msg  string
}

func (e *Err) Error() string {
	return e.Msg
}

var (
	COMMON_SUC        = &Err{10000, "success"}
	COMMON_BADREQUEST = &Err{10001, "请求错误，请检查请求是否正确"}
	COMMON_PARAM_ERR  = &Err{10002, "请求参数错误，请更正"}
	COMMON_PARAM_MISS = &Err{10003, "请求参数缺失，请检查"}
	COMMON_NOT_FOUND  = &Err{10004, "请求资源未找到"}
	COMMON_CONFILCT   = &Err{10005, "数据已被别人更改，请重试"}

	COMMON_INTERNAL_ERR             = &Err{15000, "内部错误，请稍后重试"}
	COMMON_INTERNAL_CALLING_ERR     = &Err{15001, "内部处理失败，请稍后重试"}
	COMMON_INTERNAL_CALLING_TIMEOUT = &Err{15002, "内部调用无响应，请稍后重试"}

	AUTH_NOTLOGIN     = &Err{40001, "用户未登录，请登录"}
	AUTH_NOPERMISSION = &Err{40003, "您无此权限，请联系管理员添加角色"}
)
