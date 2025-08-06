package utils

// GenerateNanoID 生成NanoID
// size: ID长度，默认21位
// alphabet: 字符集，默认使用URL友好的字符
import (
	gonanoid "github.com/matoous/go-nanoid/v2"
)

func GenerateNanoID(size ...int) (string, error) {
	if len(size) == 0 {
		return gonanoid.New() // 默认21位
	}
	return gonanoid.Generate("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", size[0])
}
