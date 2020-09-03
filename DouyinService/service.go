package DouyinService

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Info struct {
	UserInfo struct {
		UniqueId string `json:"unique_id"`
		TotalFavorited string `json:"total_favorited"`
		Nickname string `json:"nickname"`
		FollowerCount int `json:"follower_count"`
		AwemeCount int `json:"aweme_count"`
		Signature string `json:"signature"`
		AvatarMedium struct {
			Uri string `json:"uri"`
			UrlList []string `json:"url_list"`
		} `json:"avatar_medium"`
	} `json:"user_info"`
}

func GetDouyinInfo(secID string) (info Info, error error) {
	url := "https://www.iesdouyin.com/web/api/v2/user/info/?sec_uid=" + secID
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return info, err
	}

	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("accept-language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Add("Cookie", "_ga=GA1.2.685263550.1587277283; _gid=GA1.2.143250871.1587911549; tt_webid=6820028204934923790; _ba=BA0.2-20200301-5199e-c7q9NP0laGm7KfaPfGcH")
	req.Header.Add("status", "302")
	res, err := client.Do(req)
	if err == nil {
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return info, err
		}

		err = json.Unmarshal(b, &info)
		if err != nil {
			return info, err
		}
	}
	return
}