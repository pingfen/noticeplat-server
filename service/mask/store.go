package mask

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pingfen/notify/pkg/client"
	"github.com/pingfen/notify/pkg/models"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

var (
	all_mask_store []*models.MaskStore
)

func init() {
	all_mask_store = make([]*models.MaskStore, 0)
	err := json.Unmarshal([]byte(MaskStoreJson), &all_mask_store)
	if err != nil {
		panic(err)
	}
}

func ListMaskStores() []*models.MaskStore {
	return all_mask_store
}

func ListMaskByArea(area string) []*models.MaskStore {
	ret := make([]*models.MaskStore, 0)

	for _, s := range all_mask_store {
		if s.Area == s.Area {
			ret = append(ret, s)
		}
	}

	return ret
}

const URL_MASK_COUNT = "http://kzyynew.qingdao.gov.cn:81/kz/listYd"

type BaseResp struct {
	Msg  string `json:"msg"`
	Code string `json:"code"`
}
type MaskStoreResp struct {
	*BaseResp
	Data *struct {
		List []*models.MaskStore `json:"list"`
	} `json:"data"`
}

var (
	mutex            sync.RWMutex
	lastQuerySucTime time.Time
	cacheMasStore    []*models.MaskStore
)

func GetMaskStoreInventory(ctx context.Context) ([]*models.MaskStore, error) {
	ms := readMaskStoreInventoryFromCache()
	if ms != nil {
		return ms, nil
	}

	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	ms, err := getMaskStoreInventory(ctx)
	if err != nil {
		return nil, err
	}

	resetCacheMasStore(ms)

	return ms, nil
}

func readMaskStoreInventoryFromCache() []*models.MaskStore {
	mutex.RLock()
	defer mutex.RUnlock()

	if time.Now().Sub(lastQuerySucTime) < 3*time.Second {
		return cacheMasStore
	}

	return nil
}

func resetCacheMasStore(ms []*models.MaskStore) {
	mutex.Lock()
	defer mutex.Unlock()

	cacheMasStore = ms
	lastQuerySucTime = time.Now()
}

func getMaskStoreInventory(ctx context.Context) ([]*models.MaskStore, error) {
	req, err := http.NewRequest(http.MethodGet, URL_MASK_COUNT, nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 11_0 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Mobile/15A372 Safari/604.1")

	resp, err := client.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("获取口罩库存失败(%d)", resp.StatusCode)
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取口罩库存数据失败%w", err)
	}

	msr := new(MaskStoreResp)
	err = json.Unmarshal(bytes, msr)
	if err != nil {
		return nil, err
	}

	if msr.Code == "0" && msr.Data != nil {
		return msr.Data.List, nil
	} else {
		return nil, fmt.Errorf("%s(%s)", msr.Msg, msr.Code)
	}
}

const (
	MaskStoreJson = `[
      {
        "area": "崂山区",
        "address": "海口路33号-48号网点（皇家美孚旁）",
        "orgcode": "91370212730621124B",
        "name": "国风大药房永顺药店",
        "id": 105,
        "kznum": -1
      },
      {
        "area": "市北区",
        "address": "市北区延安三路13号东",
        "orgcode": "91370203706459083J",
        "name": "国风大药房中华药店",
        "id": 106,
        "kznum": -1
      },
      {
        "area": "李沧区",
        "address": "李沧区向阳路23号",
        "orgcode": "91370213863913606W",
        "name": "国风大药房崂山药店",
        "id": 107,
        "kznum": -1
      },
      {
        "area": "市北区",
        "address": "市北区同安路719号",
        "orgcode": "913702037334986276",
        "name": "国风大药房国风药店",
        "id": 108,
        "kznum": -1
      },
      {
        "area": "市南区",
        "address": "市南区高田路8号（大润发超市西侧）",
        "orgcode": "91370202770270007T      ",
        "name": "百姓阳光大药房一分店",
        "id": 109,
        "kznum": -1
      },
      {
        "area": "市北区",
        "address": "市北区标山路66-8",
        "orgcode": "91370203675262981C      ",
        "name": "百姓阳光大药房二十三分店",
        "id": 110,
        "kznum": -1
      },
      {
        "area": "胶州市",
        "address": "胶州市广州北路187号（金沙广场大润发北）",
        "orgcode": "9137028107328733XX",
        "name": "百姓阳光大药房五十八分店",
        "id": 111,
        "kznum": -1
      },
      {
        "area": "市北区",
        "address": "青岛市市北区洛阳路与商丘路交叉口",
        "orgcode": "913702135757790196        ",
        "name": "百姓阳光大药房三十七分店",
        "id": 112,
        "kznum": -1
      },
      {
        "area": "市北区",
        "address": "青岛市市北区敦化路49号甲",
        "orgcode": "91370203MA3D5W7L5E",
        "name": "青岛漱玉平民大药房第二十一分店",
        "id": 113,
        "kznum": -1
      },
      {
        "area": "市南区",
        "address": "青岛市市南区华严路4号",
        "orgcode": "91370202MA3D5T5F51",
        "name": "青岛漱玉平民大药房第七十分店",
        "id": 114,
        "kznum": -1
      },
      {
        "area": "市南区",
        "address": "青岛市市南区四川路29号丁",
        "orgcode": "91370202MA3D5T9M87",
        "name": "青岛漱玉平民大药房第四十八分店",
        "id": 115,
        "kznum": -1
      },
      {
        "area": "李沧区",
        "address": "青岛市李沧区书院路59号书院商城一楼西侧",
        "orgcode": "91370213MA3D6A792Y",
        "name": "青岛漱玉平民大药房第三十五分店",
        "id": 116,
        "kznum": -1
      },
      {
        "area": "平度市",
        "address": "平度市东阁街道办事处常州路17号",
        "orgcode": "91370281MA3BXRKL5G",
        "name": "青岛修和堂药业有限公司二十分店",
        "id": 117,
        "kznum": -1
      },
      {
        "area": "平度市",
        "address": "平度市南村镇双泉路52-4",
        "orgcode": "91370283334056004M",
        "name": "青岛修和堂药业有限公司十八分店",
        "id": 118,
        "kznum": -1
      },
      {
        "area": "黄岛区",
        "address": "青岛市黄岛区江山南路534号1号网点",
        "orgcode": "91370211MA5C45UR0P",
        "name": "青岛修和堂药业有限公司十七分店",
        "id": 119,
        "kznum": -1
      },
      {
        "area": "黄岛区",
        "address": "青岛市黄岛区紫金山路57号",
        "orgcode": "91370211MA3C4ACU45",
        "name": "青岛修和堂药业有限公司八十三分店",
        "id": 120,
        "kznum": -1
      },
      {
        "area": "城阳区",
        "address": "城阳区惜福镇驻地",
        "orgcode": "91370214766717451L",
        "name": "同方第二连锁店",
        "id": 121,
        "kznum": -1
      },
      {
        "area": "城阳区",
        "address": "城阳区春阳路大润发肯德基旁",
        "orgcode": "9137021466126461XC",
        "name": "同方第八十二连锁店",
        "id": 122,
        "kznum": -1
      },
      {
        "area": "黄岛区",
        "address": "黄岛区崇明岛路",
        "orgcode": "91370211MA3F3PPR03",
        "name": "同方280店",
        "id": 123,
        "kznum": -1
      },
      {
        "area": "即墨区",
        "address": "即墨区嵩山二路187-20号壹品华庭1层1层户",
        "orgcode": "91370282MA3D7XUH2P",
        "name": "同方208店",
        "id": 124,
        "kznum": -1
      },
      {
        "area": "胶州市",
        "address": "胶州市常州路291号",
        "orgcode": "9137028169376722XW",
        "name": "同方第九十七连锁店",
        "id": 125,
        "kznum": -1
      },
      {
        "area": "莱西市",
        "address": "莱西市济南路翡翠城南门12号",
        "orgcode": "91370285MA3EUY1N6M",
        "name": "同方382店",
        "id": 126,
        "kznum": -1
      },
      {
        "area": "城阳区",
        "address": "城阳区正阳路136号",
        "orgcode": "9137021470646291X2",
        "name": "青岛家佳源（城阳）购物中心有限公司",
        "id": 127,
        "kznum": -1
      },
      {
        "area": "即墨区",
        "address": "即墨区鹤山路999号",
        "orgcode": "9137028256470779X5",
        "name": "青岛家佳源（即墨）购物中心有限公司",
        "id": 128,
        "kznum": -1
      },
      {
        "area": "市南区",
        "address": "市南区费县路6号 ",
        "orgcode": "913702001635845060",
        "name": "青岛华联商厦股份有限公司",
        "id": 129,
        "kznum": -1
      },
      {
        "area": "市北区",
        "address": "市北区中山路149号",
        "orgcode": "913702007306189537",
        "name": "丽达中山路店",
        "id": 130,
        "kznum": -1
      },
      {
        "area": "崂山区",
        "address": "崂山区秦岭路18号",
        "orgcode": "913702006825570597",
        "name": "丽达崂山店",
        "id": 131,
        "kznum": -1
      },
      {
        "area": "李沧区",
        "address": "李沧区铜川路216号负一层超市",
        "orgcode": "91370213MA3DQENL2P",
        "name": "丽达绿城店",
        "id": 132,
        "kznum": -1
      },
      {
        "area": "市北区",
        "address": "市北区延安三路188号",
        "orgcode": "91370203321432955E",
        "name": "丽达延安三路店",
        "id": 133,
        "kznum": -1
      },
      {
        "area": "城阳区",
        "address": "城阳区崇阳路491号",
        "orgcode": "91370214733496939X",
        "name": "丽达城阳店",
        "id": 134,
        "kznum": -1
      },
      {
        "area": "胶州市",
        "address": "胶州市郑州西路5号",
        "orgcode": "9137028172404966XP",
        "name": "丽达胶州店",
        "id": 135,
        "kznum": -1
      },
      {
        "area": "市南区",
        "address": "市南区宁夏路162号",
        "orgcode": "91370200725550149T",
        "name": "大润发集团",
        "id": 136,
        "kznum": -1
      },
      {
        "area": "城阳区",
        "address": "城阳区春阳路167号",
        "orgcode": "91370214717867795B",
        "name": "大润发集团",
        "id": 137,
        "kznum": -1
      },
      {
        "area": "即墨区",
        "address": "即墨区振华街93号",
        "orgcode": "9137020071788240XD",
        "name": "大润发集团",
        "id": 138,
        "kznum": -1
      },
      {
        "area": "平度市",
        "address": "平度市苏州路4号",
        "orgcode": "913702006903383743",
        "name": "大润发集团",
        "id": 139,
        "kznum": -1
      },
      {
        "area": "胶州市",
        "address": "胶州市广州北路187号",
        "orgcode": "913702815577029589",
        "name": "大润发集团",
        "id": 140,
        "kznum": -1
      },
      {
        "area": "李沧区",
        "address": "青岛市李沧区向阳路64号",
        "orgcode": "91370200163921867G",
        "name": "青岛利客来集团股份有限公司购物中心A座超市",
        "id": 141,
        "kznum": -1
      },
      {
        "area": "城阳区",
        "address": "青岛市城阳区正阳路388号",
        "orgcode": "913702145577388383",
        "name": "青岛利客来集团城阳购物广场",
        "id": 142,
        "kznum": -1
      },
      {
        "area": "市南区",
        "address": "青岛市市南区中山路67号",
        "orgcode": "91370200693790083H",
        "name": "青岛悦喜客来中山购物有限公司",
        "id": 143,
        "kznum": -1
      },
      {
        "area": "莱西市",
        "address": "莱西市烟台路62号",
        "orgcode": "913702857372994812",
        "name": "青岛利客来集团莱西购物有限公司",
        "id": 144,
        "kznum": -1
      },
      {
        "area": "平度市",
        "address": "平度市福州路80号",
        "orgcode": "91370283794046368B",
        "name": "青岛利客来集团平度购物有限公司",
        "id": 145,
        "kznum": -1
      },
      {
        "area": "崂山区",
        "address": "崂山区海尔路83号",
        "orgcode": "91370212MA3R5EF40N",
        "name": "利群金鼎广场",
        "id": 146,
        "kznum": -1
      },
      {
        "area": "市北区",
        "address": "市北区台东三路77号",
        "orgcode": "91370200760291787J",
        "name": "利群商厦",
        "id": 147,
        "kznum": -1
      },
      {
        "area": "即墨区",
        "address": "即墨区鹤山路144号",
        "orgcode": "91370282718013054R",
        "name": "利群集团即墨商厦",
        "id": 148,
        "kznum": -1
      },
      {
        "area": "胶州市",
        "address": "胶州市澳门路298号",
        "orgcode": "91370281693781355F",
        "name": "利群胶州购物广场",
        "id": 149,
        "kznum": -1
      },
      {
        "area": "黄岛区",
        "address": "黄岛区凤凰山路169号",
        "orgcode": "913702116752587123",
        "name": "胶南购物中心",
        "id": 150,
        "kznum": -1
      },
      {
        "area": "莱西市",
        "address": "莱西市上海中路25号",
        "orgcode": "91370285086488907K",
        "name": "利群莱西购物广场",
        "id": 151,
        "kznum": -1
      },
      {
        "area": "平度市",
        "address": "平度市扬州路58号",
        "orgcode": "9137028379081901XG",
        "name": "北方国贸购物中心（平度店）",
        "id": 152,
        "kznum": -1
      },
      {
        "area": "市南区",
        "address": "市南区闽江路37号",
        "orgcode": "92370202MA3PXHDX49",
        "name": "齐脉健康管理中心",
        "id": 153,
        "kznum": -1
      },
      {
        "area": "市北区",
        "address": "市北区昌乐路1号1-217",
        "orgcode": "91370203MA3CMMX892",
        "name": "青岛品质生活百货商店",
        "id": 154,
        "kznum": -1
      }
    ]`
)
