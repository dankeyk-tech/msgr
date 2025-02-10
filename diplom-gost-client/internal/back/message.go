package back

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"diplom-chat-gost/internal/model"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fyne.io/fyne/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nfnt/resize"
	"github.com/nickalie/go-webpbin"
	"github.com/restream/reindexer/v3"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func SendMessage(id int64, text string, msgType int, config model.Config) error {
	key, err := CheckExistingChat(config, id)
	if err != nil {
		return errors.New("check existing chat: " + err.Error())
	}

	claims := model.JWTCustomClaims{}

	if _, err = jwt.ParseWithClaims(config.Token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWT.Key), nil
	}); err != nil {
		return errors.New("jwt parse with claims: " + err.Error())
	}

	idKey, err := strconv.Atoi(claims.Audience[0])
	if err != nil {
		return errors.New("atoi audience: " + err.Error())

	}

	rec, found := config.DB.Query("key").WhereInt64("uid", reindexer.EQ, int64(idKey)).Get()
	if !found {
		return errors.New("get key: key with this uid doesn't exist")
	}

	privateKeyInt := rec.(*model.KeyItem).Key
	privateKeyByte := make([]byte, len(privateKeyInt))
	for idx, elem := range privateKeyInt {
		privateKeyByte[idx] = byte(elem)
	}

	private, err := x509.ParsePKCS1PrivateKey(privateKeyByte)
	if err != nil {
		return errors.New("parse private key: " + err.Error())
	}

	chatKeyEncrypted := make([]byte, len(key.Key))
	for idx, elem := range key.Key {
		chatKeyEncrypted[idx] = byte(elem)
	}

	chatKey, err := rsa.DecryptPKCS1v15(rand.Reader, private, chatKeyEncrypted)
	if err != nil {
		return errors.New("decrypt chat key: " + err.Error())
	}

	encryptedMessage := EncryptK([]byte(text), string(chatKey))

	encryptedMessageInt := make([]int32, len(encryptedMessage))
	for idx, elem := range encryptedMessage {
		encryptedMessageInt[idx] = int32(elem)
	}

	req := model.SendMessageReq{
		ReceiverID: id,
		Text:       encryptedMessageInt,
		Type:       msgType,
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return errors.New("marshal" + err.Error())
	}

	r := bytes.NewReader(reqBody)

	request, err := http.NewRequest("POST", config.ServerDomain+"/send/message", r)
	if err != nil {
		return errors.New("new request: " + err.Error())
	}
	request.Header.Set("Authorization", "Bearer "+config.Token)

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return errors.New("do request: " + err.Error())
	}

	if response.StatusCode != http.StatusOK {
		return errors.New("bad request: " + strconv.Itoa(response.StatusCode))
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return errors.New("read all: " + err.Error())
	}

	var res model.SendMessageRes

	if err = json.Unmarshal(body, &res); err != nil {
		return errors.New("unmarshal: " + err.Error())
	}

	if res.Code != 200 {
		return errors.New("response message: " + res.Message)
	}

	return nil
}

func SendFiles(id int64, reader fyne.URIReadCloser, fType string, config model.Config) error {
	var text string
	var msgType int

	if fType == "photo" {
		msgType = 2

		fileNameArr := strings.Split(reader.URI().Name(), ".")

		var fileType int
		if fileNameArr[len(fileNameArr)-1] == "jpeg" || fileNameArr[len(fileNameArr)-1] == "jpg" {
			fileType = 1
		} else if fileNameArr[len(fileNameArr)-1] == "png" {
			fileType = 2
		} else if fileNameArr[len(fileNameArr)-1] == "webp" {
			fileType = 3
		}

		if fileType == 0 {
			return errors.New("check expansion: wrong file expansion")
		}

		fileOpened, err := os.Open(reader.URI().Path())
		if err != nil {
			return errors.New("file open: " + err.Error())
		}

		file, err := os.Create("./temp.png")
		if err != nil {
			return errors.New("file create: " + err.Error())
		}

		var input image.Image
		if fileType == 3 {
			input, err = webpbin.Decode(fileOpened)
			if err != nil {
				return errors.New("webp decode: " + err.Error())
			}
		}

		if fileType == 2 {
			input, err = png.Decode(fileOpened)
			if err != nil {
				return errors.New("png decode: " + err.Error())
			}
		}

		if fileType == 1 {
			input, err = jpeg.Decode(fileOpened)
			if err != nil {
				return errors.New("jpeg decode: " + err.Error())
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
			errors.New("png encode: " + err.Error())
		}

		fileOpened.Close()
		file.Close()

		data, err := os.ReadFile("./temp.png")
		if err != nil {
			return errors.New("read file: " + err.Error())
		}
		text = base64.StdEncoding.EncodeToString(data)

		if err = os.Remove("./temp.png"); err != nil {
			return errors.New("remove: " + err.Error())
		}
	} else {
		msgType = 3
		data, err := os.ReadFile(reader.URI().Path())
		if err != nil {
			return errors.New("read file: " + err.Error())
		}
		text = strings.ReplaceAll(reader.URI().Name(), " ", "_") + " " + base64.StdEncoding.EncodeToString(data)
	}

	if err := SendMessage(id, text, msgType, config); err != nil {
		return errors.New("send message: " + err.Error())
	}

	return nil
}
