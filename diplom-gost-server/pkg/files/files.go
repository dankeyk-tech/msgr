package files

import (
	custom_errors "diplom-chat-gost-server/pkg/custom-errors"
	"encoding/base64"
	"github.com/nfnt/resize"
	"github.com/nickalie/go-webpbin"
	"github.com/valyala/fasthttp"
	"image"
	"image/jpeg"
	"image/png"
	"mime/multipart"
	"os"
	"strings"
)

func SaveAvatar(fileHeader *multipart.FileHeader, id int64, configPath string) (string, *custom_errors.ErrHttp) {
	fileNameArr := strings.Split(fileHeader.Filename, ".")

	var fileType int
	if fileNameArr[len(fileNameArr)-1] == "jpeg" || fileNameArr[len(fileNameArr)-1] == "jpg" {
		fileType = 1
	} else if fileNameArr[len(fileNameArr)-1] == "png" {
		fileType = 2
	} else if fileNameArr[len(fileNameArr)-1] == "webp" {
		fileType = 3
	}

	if fileType == 0 {
		return "", custom_errors.New(fasthttp.StatusUnprocessableEntity, "wrong file expansion")
	}

	fileOpened, err := fileHeader.Open()
	if err != nil {
		return "", custom_errors.New(fasthttp.StatusInternalServerError, "file open: "+err.Error())
	}

	path := "./temp.png"

	file, err := os.Create(path)
	if err != nil {
		return "", custom_errors.New(fasthttp.StatusInternalServerError, "file create: "+err.Error())
	}

	var input image.Image
	if fileType == 3 {
		input, err = webpbin.Decode(fileOpened)
		if err != nil {
			return "", custom_errors.New(fasthttp.StatusInternalServerError, "webp decode: "+err.Error())
		}
	}

	if fileType == 2 {
		input, err = png.Decode(fileOpened)
		if err != nil {
			return "", custom_errors.New(fasthttp.StatusInternalServerError, "png decode: "+err.Error())
		}
	}

	if fileType == 1 {
		input, err = jpeg.Decode(fileOpened)
		if err != nil {
			return "", custom_errors.New(fasthttp.StatusInternalServerError, "jpeg decode: "+err.Error())
		}
	}

	if input.Bounds().Dx() > 700 {
		coef := float32(700) / float32(input.Bounds().Dx())
		widthNew := float32(input.Bounds().Dx()) * coef
		heightNew := float32(input.Bounds().Dy()) * coef
		input = resize.Resize(uint(widthNew), uint(heightNew), input, resize.Lanczos3)
	}

	err = png.Encode(file, input)
	if err != nil {
		return "", custom_errors.New(fasthttp.StatusInternalServerError, "png encode: "+err.Error())
	}

	fileOpened.Close()
	file.Close()

	data, err := os.ReadFile(path)
	if err != nil {
		return "", custom_errors.New(fasthttp.StatusInternalServerError, "read file: "+err.Error())
	}
	fString := base64.StdEncoding.EncodeToString(data)

	if err = os.Remove(path); err != nil {
		return "", custom_errors.New(fasthttp.StatusInternalServerError, "remove: "+err.Error())
	}

	return fString, nil
}
