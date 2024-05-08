package utils

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/png"
)

func RGBAToBase64(rgba *image.RGBA) (string, error) {
	// 创建一个与RGBA图像大小相同的字节缓冲区
	var buf bytes.Buffer

	// 将RGBA图像编码为PNG
	err := png.Encode(&buf, rgba)
	if err != nil {
		return "", err
	}

	// 将PNG字节编码为Base64字符串
	base64Str := base64.StdEncoding.EncodeToString(buf.Bytes())

	return base64Str, nil
}
