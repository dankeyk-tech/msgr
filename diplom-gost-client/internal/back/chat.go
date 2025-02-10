package back

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"diplom-chat-gost/internal/model"
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/restream/reindexer/v3"
	"io"
	"net/http"
	"strconv"
)

func GetChatByUsers(id int64, config model.Config) (int64, error) {
	req, err := http.NewRequest("GET", config.ServerDomain+"/get/chat-by-users?id="+strconv.Itoa(int(id)), nil)
	if err != nil {
		return -1, errors.New("new request: " + err.Error())
	}
	req.Header.Set("Authorization", "Bearer "+config.Token)

	client := &http.Client{}

	response, err := client.Do(req)
	if err != nil {
		return -1, errors.New("do request: " + err.Error())
	}

	if response.StatusCode != http.StatusOK {
		return -1, errors.New("bad request: " + strconv.Itoa(response.StatusCode))
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return -1, errors.New("read all: " + err.Error())
	}

	var res model.GetChatByUsersRes

	if err = json.Unmarshal(body, &res); err != nil {
		return -1, errors.New("unmarshal: " + err.Error())
	}

	if res.Code != 200 {
		return -1, errors.New("response message: " + res.Message)
	}

	return res.Data, nil
}

func GetAllChats(config model.Config) ([]*model.ChatShortItem, error) {
	req, err := http.NewRequest("GET", config.ServerDomain+"/get-all/chat", nil)
	if err != nil {
		return nil, errors.New("new request: " + err.Error())
	}
	req.Header.Set("Authorization", "Bearer "+config.Token)

	client := &http.Client{}

	response, err := client.Do(req)
	if err != nil {
		return nil, errors.New("do request: " + err.Error())
	}

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("bad request: " + strconv.Itoa(response.StatusCode))
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.New("read all: " + err.Error())
	}

	var res model.GetAllChatRes

	if err = json.Unmarshal(body, &res); err != nil {
		return nil, errors.New("unmarshal: " + err.Error())
	}

	if res.Code != 200 {
		return nil, errors.New("response message: " + res.Message)
	}

	claims := model.JWTCustomClaims{}

	if _, err = jwt.ParseWithClaims(config.Token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWT.Key), nil
	}); err != nil {
		return nil, errors.New("jwt parse with claims: " + err.Error())
	}

	uid, err := strconv.Atoi(claims.Audience[0])
	if err != nil {
		return nil, errors.New("atoi audience: " + err.Error())
	}

	rec, found := config.DB.Query("key").WhereInt64("uid", reindexer.EQ, int64(uid)).Get()
	if !found {
		return nil, errors.New("key with this uid doesn't exist")
	}

	privateKeyInt := rec.(*model.KeyItem).Key
	privateKeyByte := make([]byte, len(privateKeyInt))
	for idx, elem := range privateKeyInt {
		privateKeyByte[idx] = byte(elem)
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyByte)
	if err != nil {
		return nil, errors.New("parse private key: " + err.Error())
	}

	for idx, elem := range res.Data {
		if elem.LastMessageType == 3 {
			textByte := []byte("Файл")
			textInt := make([]int32, len(textByte))
			for idxText, elemText := range textByte {
				textInt[idxText] = int32(elemText)
			}
			res.Data[idx].LastMessageText = textInt
			continue
		} else if elem.LastMessageType == 2 {
			textByte := []byte("Изображение")
			textInt := make([]int32, len(textByte))
			for idxText, elemText := range textByte {
				textInt[idxText] = int32(elemText)
			}
			res.Data[idx].LastMessageText = textInt
			continue
		}

		encryptedKey, err := GetChatKey(config, elem.ID)
		if err != nil {
			return nil, errors.New("get chat key: " + err.Error())
		}

		encryptedKeyByte := make([]byte, len(encryptedKey))
		for idx, elem := range encryptedKey {
			encryptedKeyByte[idx] = byte(elem)
		}

		chatKey, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, encryptedKeyByte)
		if err != nil {
			return nil, errors.New("decrypt chat key: " + err.Error())
		}

		text := make([]byte, len(elem.LastMessageText))
		for idx1, elem1 := range elem.LastMessageText {
			text[idx1] = byte(elem1)
		}

		decryptedMessage := DecryptK(text, string(chatKey))

		decryptedMessageInt := make([]int32, len(decryptedMessage))
		for idx1, elem1 := range decryptedMessage {
			decryptedMessageInt[idx1] = int32(elem1)
		}

		res.Data[idx].LastMessageText = decryptedMessageInt
	}

	return res.Data, nil
}

func GetChat(id int64, config model.Config) ([]*model.MessageShortItem, error) {
	req, err := http.NewRequest("GET", config.ServerDomain+"/get/chat?chat-id="+strconv.Itoa(int(id)), nil)
	if err != nil {
		return nil, errors.New("new request: " + err.Error())
	}
	req.Header.Set("Authorization", "Bearer "+config.Token)

	client := &http.Client{}

	response, err := client.Do(req)
	if err != nil {
		return nil, errors.New("do request: " + err.Error())
	}

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("bad request: " + strconv.Itoa(response.StatusCode))
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.New("read all: " + err.Error())
	}

	var res model.GetAllMessagesRes

	if err = json.Unmarshal(body, &res); err != nil {
		return nil, errors.New("unmarshal: " + err.Error())
	}

	if res.Code != 200 {
		return nil, errors.New("response message: " + res.Message)
	}

	encryptedKey, err := GetChatKey(config, id)
	if err != nil {
		return nil, errors.New("get chat key: " + err.Error())
	}

	claims := model.JWTCustomClaims{}

	if _, err = jwt.ParseWithClaims(config.Token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWT.Key), nil
	}); err != nil {
		return nil, errors.New("jwt parse with claims: " + err.Error())
	}

	uid, err := strconv.Atoi(claims.Audience[0])
	if err != nil {
		return nil, errors.New("atoi audience: " + err.Error())
	}

	rec, found := config.DB.Query("key").WhereInt64("uid", reindexer.EQ, int64(uid)).Get()
	if !found {
		return nil, errors.New("key with this uid doesn't exist")
	}

	privateKeyInt := rec.(*model.KeyItem).Key
	privateKeyByte := make([]byte, len(privateKeyInt))
	for idx, elem := range privateKeyInt {
		privateKeyByte[idx] = byte(elem)
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyByte)
	if err != nil {
		return nil, errors.New("parse private key: " + err.Error())
	}

	encryptedKeyByte := make([]byte, len(encryptedKey))
	for idx, elem := range encryptedKey {
		encryptedKeyByte[idx] = byte(elem)
	}

	chatKey, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, encryptedKeyByte)
	if err != nil {
		return nil, errors.New("decrypt chat key: " + err.Error())
	}

	for idx, elem := range res.Data {
		text := make([]byte, len(elem.Text))
		for idx1, elem1 := range elem.Text {
			text[idx1] = byte(elem1)
		}

		decryptedMessage := DecryptK(text, string(chatKey))

		decryptedMessageInt := make([]int32, len(decryptedMessage))
		for idx1, elem1 := range decryptedMessage {
			decryptedMessageInt[idx1] = int32(elem1)
		}

		res.Data[idx].Text = decryptedMessageInt
	}

	return res.Data, nil
}

func CheckExistingChat(config model.Config, receiverID int64) (*model.CheckChatItem, error) {
	req := model.CheckChatReq{
		ReceiverID: receiverID,
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, errors.New("marshal" + err.Error())
	}

	r := bytes.NewReader(reqBody)

	request, err := http.NewRequest("POST", config.ServerDomain+"/check/chat", r)
	if err != nil {
		return nil, errors.New("new request: " + err.Error())
	}
	request.Header.Set("Authorization", "Bearer "+config.Token)

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return nil, errors.New("do request: " + err.Error())
	}

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("bad request: " + strconv.Itoa(response.StatusCode))
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.New("read all: " + err.Error())
	}

	var res model.CheckChatRes

	if err = json.Unmarshal(body, &res); err != nil {
		return nil, errors.New("unmarshal: " + err.Error())
	}

	if res.Code != 200 {
		return nil, errors.New("response message: " + res.Message)
	}

	return res.Data, nil
}
