package connect

import (
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"time"
)

// 客户端
var cli = http.Client{
	Transport: &http.Transport{
		DisableKeepAlives: true,
	},
	Timeout: 15 * time.Second,
}

// Ping 测试target url是否可访问通
func Ping(target string) bool {
	resp, err := cli.Get(target)
	if err != nil {
		logx.Errorw("connect Ping failed", logx.Field("err", err))
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK // 不兼容跳转的长链接
}
