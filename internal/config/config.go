package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf

	//长链转短链 mysql
	ShortUrlDB struct {
		DSN string
	}

	//发号器 mysql
	SequenceDB struct {
		DSN string
	}

	//redis 缓存
	CacheConf cache.CacheConf

	//短链接黑名单/违禁词
	BlackList []string

	//短链接重定向服务域名
	ShortUrlDomain string
}
