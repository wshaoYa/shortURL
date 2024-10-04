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

	"github.com/zeromicro/go-zero/core/logx"
)

type ShowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShowLogic {
	return &ShowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShowLogic) Show(req *types.ShowReq) (resp *types.ShowResp, err error) {
	//布隆过滤器 拦截过滤
	has, err := l.svcCtx.BloomFilter.ExistsCtx(l.ctx, []byte(req.ShortURL))
	if err != nil {
		logx.Errorw("logic l.svcCtx.BloomFilter.ExistsCtx failed", logx.Field("err", err))
		return nil, Error.ErrInternal
	}
	if !has {
		return nil, Error.ErrInvalidShortURL
	}

	//校验短链是否存在
	// go-zero的缓存原生支持singleflight 防止缓存击穿
	fmt.Println("通过bloom filter，开始查询缓存/DB")
	u, err := l.svcCtx.ShortUrlDB.FindOneBySurl(l.ctx, sql.NullString{String: req.ShortURL, Valid: true})
	if err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			return nil, Error.ErrInvalidShortURL
		}
		logx.Errorw("logic l.svcCtx.ShortUrlDB.FindOneBySurl failed", logx.Field("err", err))
		return nil, Error.ErrInternal
	}

	//返回响应
	return &types.ShowResp{LongURL: u.Lurl.String}, nil
}
