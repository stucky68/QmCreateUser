package main

import (
	"QmCreateUser/DouyinService"
	"QmCreateUser/QmService"
	"QmCreateUser/Utils"
	"fmt"
	"math/rand"
	"net/http"
	url2 "net/url"
	"strings"
)

func getSecID(url string) string {
	method := "GET"
	client := &http.Client {
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return ""
	}
	res, err := client.Do(req)
	if err != nil {
		return ""
	}

	if len(res.Header["Location"]) > 0 {
		u, err := url2.Parse(res.Header["Location"][0])
		if err != nil {
			 return ""
		}
		m := u.Query()
		if len(m["sec_uid"]) > 0 {
			return m["sec_uid"][0]
		}
	}
	return ""
}

func main() {
	nickNameData := Utils.ReadFileData("./nickName.txt")
	if nickNameData == "" {
		panic("昵称表不能为空")
	}
	nickNames := strings.Split(nickNameData, "\r\n")

	fileData := Utils.ReadFileData("./data.txt")
	if fileData != "" {
		items := strings.Split(fileData, "\r\n")
		for index, value := range items {
			item := strings.Split(value, "|")
			if len(item) == 2 {
				Utils.Log(fmt.Sprintf("开始处理第%d行数据", index+1))
				qmService := QmService.NewQmService(item[0])
				secID := getSecID(item[1])
				if secID == "" {
					Utils.Log("获取SecID失败")
					continue
				}

				info, err := DouyinService.GetDouyinInfo(secID)
				if err != nil {
					Utils.Log(err)
					continue
				}

				if len(info.UserInfo.AvatarMedium.UrlList) > 0 {
					nickName := Utils.FilterNickName(info.UserInfo.Nickname)
					nickName = strings.ReplaceAll(nickName, "小姐", "小姐姐")

					nickNameRand :=rand.Intn(len(nickNames) - 1)

					begin := rand.Intn(100) % 2
					if begin == 0 {
						nickName = nickNames[nickNameRand] + nickName
					} else {
						nickName = nickName + nickNames[nickNameRand]
					}
					shareLink := qmService.Process(nickName, info.UserInfo.Signature,info.UserInfo.AvatarMedium.UrlList[0])
					Utils.Log("全民分享链接:" + shareLink)
				} else {
					Utils.Log("获取抖音头像失败 SecID:" + item[1])
				}
			}
		}
	}
}
