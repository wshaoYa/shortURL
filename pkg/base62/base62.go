package base62

import "slices"

// 安全性考虑，打乱顺序
// const characters = `0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ`
const characters = "MNjklrdiHJKL3DFU45PW78ughQe126bnmqv0YftyREpozGaBVCXs9OIxcwTZAS"

func Uint64ToBase62(x uint64) string {
	if x == 0 {
		return string(characters[0])
	}

	var (
		bs   []byte
		base uint64 = 62
	)

	for x > 0 {
		bs = append(bs, characters[x%base])
		x /= base
	}

	slices.Reverse(bs)
	return string(bs)
}
