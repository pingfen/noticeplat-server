package mask

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pingfen/notify/pkg/client"
	"github.com/pingfen/notify/pkg/models"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const (
	URL_RESERVE = "http://kzyynew.qingdao.gov.cn:81/kz/addYdorder"
)

func DealReserveRequest(ctx context.Context, uid string, rr *models.ReserveRequest) error {
	u, err := GetContact(ctx, uid, rr.Contact)
	if err != nil {
		return err
	}

	if rr.MaskStore == nil {
		rr.MaskStore = u.ExpectStore
	}

	return Reserve(ctx, u, rr.MaskStore, rr.Captcha)
}

func Reserve(ctx context.Context, u *models.MaskUserInfo, ms *models.MaskStore, captcha *models.Captcha) error {
	uv := url.Values{}
	uv.Set("code", u.CardNo)
	uv.Set("mobile", u.Mobile)
	uv.Set("name", u.Name)

	uv.Set("ydid", fmt.Sprintf("%d", ms.ID))
	uv.Set("ydname", ms.Name)
	uv.Set("ydaddress", ms.Address)

	uv.Set("capval", captcha.Data)

	req_url := URL_RESERVE + "?" + uv.Encode()
	req, err := http.NewRequest(http.MethodGet, req_url, nil)
	if err != nil {
		return err
	}

	req = req.WithContext(ctx)
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 11_0 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Mobile/15A372 Safari/604.1")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3")
	for _, c := range captcha.Cookies {
		req.AddCookie(c)
	}

	resp, err := client.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取下单操作结果失败%w", err)
	}

	log.Printf("%s", bytes)

	if resp.StatusCode != 200 {
		return fmt.Errorf("下单失败(%d)", resp.StatusCode)
	}

	msr := new(BaseResp)
	err = json.Unmarshal(bytes, msr)
	if err != nil {
		return err
	}

	if msr.Code != "0" {
		return fmt.Errorf("%s", msr.Msg)
	}

	return nil
}
