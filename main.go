package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Result struct {
	Errno int `json:"errno"`
	Errmsg string `json:"errmsg"`
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

func main() {
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
	m["bduss"] = "S1rZlJGZHItMS0xfnp0Vjk4MUFBSHdLdDBwaExXRHEtVzR1bi1hLU53bVJlWWRlRUFBQUFBJCQAAAAAAQAAAAEAAABJDhECZXE4NTQ1MzMzMAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAJHsX16R7F9eS"
	m["portrait_type"] = "0"

	for k, v := range m {
		writer.WriteField(k, v)
	}
	sig := calculateSig(m, "0c7c877d1c7825fa4438c44dbb645d1b")
	fmt.Println(sig)
	writer.WriteField("sig", sig)

	//_ = writer.WriteField("zid", "VhWdd6aWoVz4pg6CYjp1yysdG8r69POjHcxLHy1pyQ7Cwbj8kVb_mEkhPqKH3ASpKY_6e9v6QxAEYmDzuhNA5FA")
	//_ = writer.WriteField("app_version", "2.3.2.10")
	//_ = writer.WriteField("cuid", "133B4B0017CB08044A8983B54D34C3F5")
	//_ = writer.WriteField("sdkversion", "8.9.3")
	//_ = writer.WriteField("client", "android")
	//_ = writer.WriteField("sdk_version", "8.9.3")
	//_ = writer.WriteField("clientip", "10.0.2.15")
	//_ = writer.WriteField("portrait_type", "0")
	//_ = writer.WriteField("appid", "1")
	//_ = writer.WriteField("clientid", "133B4B0017CB08044A8983B54D34C3F5")
	//_ = writer.WriteField("bduss", "lPNGt4S3dsWTF2QXJCZWNkb2tJRldJenpmeHVyU3NIeXhZQnlMVEkyNlZMbmRmSUFBQUFBJCQAAAAAAAAAAAEAAADhyKIic3R1Y2t5ODg4OAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAJWhT1-VoU9fUz")
	//_ = writer.WriteField("sig", "a312f1b53cc1bc3d3feaf0f249463c53")
	//_ = writer.WriteField("tpl", "bdmv")
	file, _ := ioutil.ReadFile("C:/Users/stucky/Desktop/20200902214009.jpg")

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`, escapeQuotes("file"), escapeQuotes("portrait.jpg")))
	h.Set("Content-Type", "image/jpeg")

	part1, _ := writer.CreatePart(h)
	part1.Write(file)

	err := writer.Close()
	if err != nil {
		fmt.Println(err)
	}

	client := &http.Client {
	}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
	}

	req.Header.Add("XRAY-TRACEID", "42a53425-a2ed-4064-9e74-2eda077fc580")
	req.Header.Add("XRAY-REQ-FUNC-ST-DNS", "httpsUrlConn;" + strconv.FormatInt(time.Now().UnixNano() / 10e5, 10) + ";5")
	req.Header.Add("Content-Type", "multipart/form-data;boundary=" + writer.Boundary())
	req.Header.Add("User-Agent", "tpl:bdmv;android_sapi_v8.9.3")
	req.Header.Add("Host", "passport.baidu.com")
	req.Header.Add("Connection", "Keep-Alive")
	req.Header.Add("Accept-Encoding", "gzip")
	req.Header.Add("Content-Length", strconv.Itoa(len(payload.String())))
	res, err := client.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	result := &Result{}
	json.Unmarshal(body, result)

	fmt.Println(string(body), result.Errmsg)
}
