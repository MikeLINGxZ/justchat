package utils

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/jpeg"
	"image/png"
	"os"

	"golang.org/x/image/draw"
)

// GenerateImagePreview 生成200x200的base64图片预览
func GenerateImagePreview(imagePath string) (string, error) {
	// 打开图片文件
	file, err := os.Open(imagePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 解码图片
	img, format, err := image.Decode(file)
	if err != nil {
		return "", err
	}

	// 创建200x200的目标图片
	dst := image.NewRGBA(image.Rect(0, 0, 200, 200))

	// 使用双线性插值缩放图片
	draw.BiLinear.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Over, nil)

	// 将图片编码为base64
	var buf bytes.Buffer
	var mimeType string

	switch format {
	case "jpeg", "jpg":
		mimeType = "image/jpeg"
		err = jpeg.Encode(&buf, dst, &jpeg.Options{Quality: 85})
	case "png":
		mimeType = "image/png"
		err = png.Encode(&buf, dst)
	default:
		// 对于其他格式，默认使用jpeg
		mimeType = "image/jpeg"
		err = jpeg.Encode(&buf, dst, &jpeg.Options{Quality: 85})
	}

	if err != nil {
		return "", err
	}

	// 转换为base64
	base64Str := base64.StdEncoding.EncodeToString(buf.Bytes())
	return "data:" + mimeType + ";base64," + base64Str, nil
}
