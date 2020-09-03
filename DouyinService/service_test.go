package DouyinService

import (
	"fmt"
	"testing"
)

func TestGetDouyinInfo(t *testing.T) {
	info, _ := GetDouyinInfo("MS4wLjABAAAAH4tB-PlRkVElPgI5ldAe0Z_WLRjWje5vAjFMWQFykmM")
	fmt.Println(info.UserInfo.Nickname, info.UserInfo.Signature, info.UserInfo.AvatarMedium.UrlList[0])
}
