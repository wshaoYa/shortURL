package urltool

import (
	"net/url"
	"path"
	"shortURL/internal/Error"
)

// GetBasePath 获取末尾base路径
func GetBasePath(target string) (string, error) {
	parseURL, err := url.Parse(target)
	if err != nil {
		return "", err
	}
	//避免相对路径 （但其实根本就走不到这步应该，长链接有效校验就拦住了）
	if len(parseURL.Host) == 0 {
		return "", Error.ErrUrlNoHost
	}
	return path.Base(parseURL.Path), nil
}
