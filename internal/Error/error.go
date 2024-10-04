package Error

import "errors"

var (
	ErrInvalidLongURL  = errors.New("不可访问/非法的长链接")
	ErrInvalidShortURL = errors.New("不可访问/非法的短链接")
	ErrInternal        = errors.New("内部错误")
	ErrUrlNoHost       = errors.New("url缺少host信息")
	ErrLoopConvert     = errors.New("不能将短链接再次生成短链接")
)
