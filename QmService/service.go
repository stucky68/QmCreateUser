package QmService

import (
	"QmCreateUser/Utils"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"sort"
	"strconv"
	"strings"
	"time"
)

type QmService struct {
	bduss string
}

//https://quanmin.baidu.com/mvideo/api?api_name=usermessage&usermessage&refresh_type=4&curr_timestamp=1599143837&aid=1 查看消息
type config struct {
	Setportrait bool `json:"Setportrait"`
	SetAutograph bool `json:"SetAutograph"`
	SetBirthday bool `json:"SetBirthday"`
	SetNickName bool `json:"SetNickName"`
}

var g_config config

func init()  {
	rand.Seed(time.Now().Unix())
	data := Utils.ReadFileData("./config.json")
	err := json.Unmarshal([]byte(data), &g_config)
	if err != nil {
		panic(err)
	}

}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

func calculateSig(m map[string]string, signKey string) string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	sb := ""
	for k := range keys {
		sb += keys[k]
		sb += "="
		sb += m[keys[k]]
		sb += "&"
	}
	sb += "sign_key="
	sb += signKey

	md5Ctx := md5.New()
	md5Ctx.Write([]byte(sb))
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString( cipherStr)
}

func NewQmService(bduss string) *QmService {
	return &QmService{bduss:bduss}
}

func (service *QmService) GetShareLink() (string, error) {
	url := "https://quanmin.baidu.com/appui/user/mine?log=vhk&tn=1021212y&ctn=1021212y&imei=440000000123735&od=&cuid=133B4B0017CB08044A8983B54D34C3F5%7CV6E7HVO6Y&bdboxcuid=null&os=android&osbranch=a0&ua=810_1440_270&ut=MuMu_6.0.1_23_Android&uh=Netease,cancro,unknown,1&apiv=1.0.0.10&appv=210&version=2.3.2.10&life=1585137086&clife=1585137086&hid=4D364BE5F5663D4E17C49955B93D9D3A&network=1&network_state=20&sids=3022_1-3025_2-3043_1-3049_3-3057_1-3123_4-3159_1-3169_1-3000066_1&teenager=0&oaid=&activity_ext=&c3_aid=A00-CQALUGBCGZMA3PX64TID2VGXKHBDVXSD-B4LHAKIP&api_name=mine&sign=3be6d6f294008e8cd500cc9946728edb"
	method := "POST"

	payload := strings.NewReader("ext=mine")

	client := &http.Client {}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return "", nil
	}
	req.Header.Add("Cookie", "BDUSS=" + service.bduss)
	req.Header.Add("Content-Type", "text/plain")

	res, err := client.Do(req)
	if err != nil {
		return "", nil
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", nil
	}

	result := &struct {
		Mine struct {
			Status int `json:"status"`
			Msg string `json:"msg"`
			Data struct {
				ShareInfo struct {
					Link string `json:"link"`
				} `json:"shareInfo"`
			} `json:"data"`
		} `json:"mine"`
	}{}

	err = json.Unmarshal(body, result)
	if err != nil {
		return "", err
	}

	if result.Mine.Status != 0 {
		return "", errors.New(result.Mine.Msg)
	}
	return result.Mine.Data.ShareInfo.Link, nil
}

func (service *QmService) SetAutograph(autograph string) error {
	return service.userProfileSubmit("userprofilesubmit=autograph%3D%" + autograph)
}

func (service *QmService) SetBirthday(birthday string) error {
	return service.userProfileSubmit("userprofilesubmit=birthday%3D" + birthday)
}

func (service *QmService) SetSex(isWoMan int) error {
	return service.userProfileSubmit("userprofilesubmit=sex%3D" + strconv.Itoa(isWoMan))
}

func (service *QmService) SetNickName(nickName string) error {
	return service.userProfileSubmit("userprofilesubmit=nickname%3D" + nickName)
}

func (service *QmService) userProfileSubmit(str string) error {
	url := "https://quanmin.baidu.com/mvideo/api?log=vhk&tn=1021212y&ctn=1021212y&imei=440000000123735&od=&cuid=133B4B0017CB08044A8983B54D34C3F5%7CV6E7HVO6Y&bdboxcuid=null&os=android&osbranch=a0&ua=810_1440_270&ut=MuMu_6.0.1_23_Android&uh=Netease,cancro,unknown,1&apiv=1.0.0.10&appv=210&version=2.3.2.10&life=1585137086&clife=1585137086&hid=4D364BE5F5663D4E17C49955B93D9D3A&network=1&network_state=20&sids=3022_1-3025_2-3043_1-3049_3-3057_1-3123_4-3159_1-3169_1-3000066_1&teenager=0&oaid=&activity_ext=&c3_aid=A00-CQALUGBCGZMA3PX64TID2VGXKHBDVXSD-B4LHAKIP&api_name=userprofilesubmit&sign=729a525d2a4f62c437ee2d986a92ad79"
	method := "POST"

	payload := strings.NewReader(str)

	client := &http.Client {}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return err
	}
	req.Header.Add("Cookie", "BDUSS=" + service.bduss)
	req.Header.Add("Content-Type", "text/plain")

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	result := &struct {
		Timestamp int `json:"timestamp"`
		Logid string `json:"logid"`
		ServLogin bool `json:"servLogin"`
		Userprofilesubmit struct {
			Status int `json:"status"`
			Msg string `json:"msg"`
		} `json:"userprofilesubmit"`
	}{}

	err = json.Unmarshal(body, result)
	if err != nil {
		return err
	}


	if result.Userprofilesubmit.Status != 0 {
		return errors.New(result.Userprofilesubmit.Msg)
	}
	return nil
}

func httpGet(url string) ([]byte, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func randate() time.Time {
	min := time.Date(1980, 1, 0, 0, 0, 0,0, time.UTC).Unix()
	max := time.Date(2000, 1, 0, 0, 0, 0,0, time.UTC).Unix()
	delta := max-min
	sec := rand.Int63n(delta)+ min
	return time.Unix(sec,0)
}

func (service *QmService) Process(nickName, signature, imgUrl string) string {
	if g_config.Setportrait {
		imgData, err := httpGet(imgUrl)
		if err != nil {
			Utils.Log("获取头像失败:" + err.Error())
		} else {

			err = service.Setportrait(imgData)
			if err != nil {
				Utils.Log("上传头像失败:" + err.Error())
			} else {
				Utils.Log("上传头像成功")
			}
		}
	}

	if g_config.SetNickName {
		err := service.SetNickName(nickName)
		if err != nil {
			Utils.Log("设置昵称失败:" + err.Error())
		} else {
			Utils.Log("设置昵称成功")
		}
	}

	if g_config.SetAutograph {
		err := service.SetAutograph(signature)
		if err != nil {
			Utils.Log("设置签名失败:" + err.Error())
		} else {
			Utils.Log("设置签名成功")
		}
	}

	if g_config.SetBirthday {
		time := randate()
		err := service.SetBirthday(time.Format("20060102"))
		if err != nil {
			Utils.Log("设置生日失败:" + err.Error())
		} else {
			Utils.Log("设置生日成功")
		}
	}

	shareLink, err := service.GetShareLink()
	if err != nil {
		Utils.Log("获取共享链接失败:" + err.Error())
	}
	return shareLink
}

func (service *QmService) Setportrait(imageFile []byte) error {
	url := "https://passport.baidu.com/v2/sapi/center/setportrait"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	m := make(map[string]string)
	m["client"] = "android"
	m["cuid"] = "133B4B0017CB08044A8983B54D34C3F5"
	m["clientid"] = "133B4B0017CB08044A8983B54D34C3F5"
	m["zid"] = "VhWdd6aWoVz4pg6CYjp1yysdG8r69POjHcxLHy1pyQ7Cwbj8kVb_mEkhPqKH3ASpKY_6e9v6QxAEYmDzuhNA5FA"
	m["clientip"] = "10.0.2.15"
	m["appid"] = "1"
	m["tpl"] = "bdmv"
	m["app_version"] = "2.3.2.10"
	m["sdk_version"] = "8.9.3"
	m["sdkversion"] = "8.9.3"
	m["bduss"] = service.bduss
	m["portrait_type"] = "0"

	for k, v := range m {
		writer.WriteField(k, v)
	}
	sig := calculateSig(m, "0c7c877d1c7825fa4438c44dbb645d1b")
	writer.WriteField("sig", sig)

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`, escapeQuotes("file"), escapeQuotes("portrait.jpg")))
	h.Set("Content-Type", "image/jpeg")

	part1, _ := writer.CreatePart(h)
	part1.Write(imageFile)
	writer.Close()

	client := &http.Client {}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return err
	}

	req.Header.Add("XRAY-TRACEID", "ec8f1b76-3f3c-4e65-936b-a0a26ffaf0c8")
	req.Header.Add("XRAY-REQ-FUNC-ST-DNS", "httpsUrlConn;" + strconv.FormatInt(time.Now().UnixNano() / 10e5, 10) + ";5")
	req.Header.Add("Content-Type", "multipart/form-data;boundary=" + writer.Boundary())
	req.Header.Add("User-Agent", "tpl:bdmv;android_sapi_v8.9.3")
	req.Header.Add("Host", "passport.baidu.com")
	req.Header.Add("Connection", "Keep-Alive")
	req.Header.Add("Accept-Encoding", "gzip")
	req.Header.Add("Content-Length", strconv.Itoa(len(payload.String())))
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	result := &struct {
		Errno int `json:"errno"`
		Errmsg string `json:"errmsg"`
	}{}

	err = json.Unmarshal(body, result)
	if err != nil {
		return err
	}

	if result.Errno != 0 {
		return errors.New(result.Errmsg)
	}
	return nil
}

