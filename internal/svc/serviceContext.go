package svc

import (
	"errors"
	"github.com/zeromicro/go-zero/core/bloom"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"shortURL/internal/config"
	"shortURL/model"
	"shortURL/pkg/sequence"
)

type ServiceContext struct {
	Config      config.Config
	ShortUrlDB  model.ShortUrlMapModel // mysql model
	SequenceDB  sequence.Sequence      // 发号器
	BlackWords  map[string]struct{}    // 违禁词map
	BloomFilter *bloom.Filter
}

func NewServiceContext(c config.Config) *ServiceContext {
	//mysql连接
	conn := sqlx.NewMysql(c.ShortUrlDB.DSN)

	//初始化 违禁词map
	var m = make(map[string]struct{})
	for _, b := range c.BlackList {
		m[b] = struct{}{}
	}

	//初始化 布隆过滤器
	bloomFilter := bloom.New(redis.MustNewRedis(redis.RedisConf{
		Host: c.CacheConf[0].Host,
		Type: redis.NodeType,
	}), "show service bloom", 500)
	loadShorturlToBloomFilter(conn, bloomFilter)

	return &ServiceContext{
		Config:      c,
		ShortUrlDB:  model.NewShortUrlMapModel(conn, c.CacheConf),
		SequenceDB:  sequence.NewMySQL(c.SequenceDB.DSN),
		BlackWords:  m,
		BloomFilter: bloomFilter,
	}
}

// 服务启动时 加载初始化 布隆过滤器
func loadShorturlToBloomFilter(conn sqlx.SqlConn, filter *bloom.Filter) error {
	if conn == nil || filter == nil {
		return errors.New("svc loadShorturlToBloomFilter failed, something is nil")
	}

	//查询shorturl总数total
	var total int
	err := conn.QueryRow(&total, `select COUNT(*) from short_url_mp where is_del=0`)
	if err != nil {
		logx.Errorw("svc  conn.QueryRow failed", logx.Field("err", err))
		return err
	}

	//依total分页将数据add进布隆过滤器
	pageSize := 50
	pageNum := total / pageSize
	if total%pageSize != 0 {
		pageNum++
	}

	for page := 1; page <= pageNum; page++ {
		var shorturls []string

		offset := (page - 1) * pageSize
		err = conn.QueryRows(&shorturls, `select surl from short_url_map where is_del=0 limit ?,?`, offset, pageSize)
		if err != nil {
			logx.Errorw("svc  conn.QueryRows failed", logx.Field("err", err))
			return err
		}

		for _, su := range shorturls {
			err = filter.Add([]byte(su))
			if err != nil {
				logx.Errorw("svc filter.Add failed", logx.Field("err", err))
				return err
			}
		}
	}

	return nil
}
