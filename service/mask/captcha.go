package mask

import (
	"context"
	"fmt"
	"github.com/pingfen/notify/pkg/client"
	"github.com/pingfen/notify/pkg/models"
	"io/ioutil"
	"net/http"
)

const (
	URL_CAPTCHA = "http://kzyynew.qingdao.gov.cn:81/kz/captcha"
)

func GetCaptcha(ctx context.Context) (captcha *models.Captcha, err error) {
	var req *http.Request
	req, err = http.NewRequest(http.MethodGet, URL_CAPTCHA, nil)
	if err != nil {
		return
	}

	req = req.WithContext(ctx)
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 11_0 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Mobile/15A372 Safari/604.1")

	var resp *http.Response
	resp, err = client.HttpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = fmt.Errorf("加载验证码失败(%d)", resp.StatusCode)
		return
	}

	captcha = new(models.Captcha)

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("读取验证码图片失败%w", err)
		return
	}
	captcha.Data = string(bytes)

	captcha.Cookies = resp.Cookies()
	captcha.ContentType = resp.Header.Get("Content-Type")
	captcha = captcha.Base64()

	return
}
