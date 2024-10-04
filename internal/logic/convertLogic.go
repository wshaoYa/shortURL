package logic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"shortURL/internal/Error"
	"shortURL/internal/svc"
	"shortURL/internal/types"
	"shortURL/model"
	"shortURL/pkg/base62"
	"shortURL/pkg/connect"
	"shortURL/pkg/md5"
	"shortURL/pkg/urltool"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConvertLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConvertLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConvertLogic {
	return &ConvertLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 拼装resp中的shortUrl
func (l *ConvertLogic) assemblyShortUrl(shortUrl string) string {
	return fmt.Sprintf("%v/%v", l.svcCtx.Config.ShortUrlDomain, shortUrl)
}

func (l *ConvertLogic) Convert(req *types.ConvertReq) (resp *types.ConvertResp, err error) {
	//长链可访问校验
	if !connect.Ping(req.LongURL) {
		return nil, Error.ErrInvalidLongURL
	}

	//转md5
	md5Val := md5.GetMd5([]byte(req.LongURL))

	//已存在校验
	u, err := l.svcCtx.ShortUrlDB.FindOneByMd5(l.ctx, sql.NullString{String: md5Val, Valid: true})
	if !errors.Is(err, sqlx.ErrNotFound) {
		if err == nil { //已存在
			return &types.ConvertResp{ShortURL: l.assemblyShortUrl(u.Surl.String)}, nil
		}
		//报未知错误
		logx.Errorw("ShortUrlDB.FindOneByMd5 failed", logx.Field("err", err))
		return nil, Error.ErrInternal
	}

	//循环转链校验--取出末尾path，校验是否为短链
	shortUrl, err := urltool.GetBasePath(req.LongURL)
	if err != nil {
		logx.Errorw("urltool.GetBasePath failed", logx.Field("long url", req.LongURL), logx.Field("err", err))
		return nil, Error.ErrInternal
	}
	_, err = l.svcCtx.ShortUrlDB.FindOneBySurl(l.ctx, sql.NullString{String: shortUrl, Valid: true})
	if !errors.Is(err, sqlx.ErrNotFound) {
		if err == nil {
			return nil, Error.ErrLoopConvert
		}
		//报未知错误
		logx.Errorw("ShortUrlDB.FindOneBySurl failed", logx.Field("err", err))
		return nil, Error.ErrInternal
	}

	for {
		//发号器 取号--注意安全考虑和黑名单单词
		ID, err := l.svcCtx.SequenceDB.Next()
		if err != nil {
			logx.Errorw("logic l.svcCtx.SequenceDB.Next() failed", logx.Field("err", err))
			return nil, Error.ErrInternal
		}

		//转62进制
		shortUrl = base62.Uint64ToBase62(ID)

		//短链违法词校验
		if _, ok := l.svcCtx.BlackWords[shortUrl]; !ok {
			break // 非违禁词
		}
	}

	//入库
	_, err = l.svcCtx.ShortUrlDB.Insert(l.ctx, &model.ShortUrlMap{
		Lurl: sql.NullString{String: req.LongURL, Valid: true},
		Md5:  sql.NullString{String: md5Val, Valid: true},
		Surl: sql.NullString{String: shortUrl, Valid: true},
	})
	if err != nil {
		logx.Errorw("logic l.svcCtx.ShortUrlDB.Insert failed", logx.Field("err", err))
		return nil, Error.ErrInternal
	}

	//布隆过滤器 add
	err = l.svcCtx.BloomFilter.AddCtx(l.ctx, []byte(shortUrl))
	if err != nil {
		logx.Errorw("logic l.svcCtx.BloomFilter.AddCtx failed", logx.Field("err", err))
		return nil, Error.ErrInternal
	}

	//返回响应--顺便拼接短链接服务的host
	return &types.ConvertResp{ShortURL: l.assemblyShortUrl(shortUrl)}, nil
}
