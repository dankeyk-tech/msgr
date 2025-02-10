package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	api_db "diplom-chat-gost-server/internal/api-db"
	"diplom-chat-gost-server/internal/model"
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
	"strconv"
	"strings"
)

func GetChatKeyHandler(ctx *fasthttp.RequestCtx) (res *model.GetChatKeyRes, message string, code int) {
	defer func() {
		resFinal := model.Response{Data: res, Code: code, Message: message}
		if code != fasthttp.StatusOK {
			log.Error().Err(errors.New(message)).Msg("")
		}
		jsonRes, _ := json.Marshal(resFinal)
		ctx.Response.SetStatusCode(resFinal.Code)
		ctx.Response.SetBody(jsonRes)
	}()

	if !ctx.IsGet() {
		return nil, "handler: wrong method", fasthttp.StatusMethodNotAllowed
	}

	tokenGet := ctx.Request.Header.Peek("Authorization")
	if tokenGet == nil {
		return nil, "handler: no authorization header", fasthttp.StatusUnauthorized
	}

	claims := model.JWTCustomClaims{}

	if len(strings.Split(string(tokenGet), " ")) == 1 {
		return nil, "handler: token is invalid", fasthttp.StatusForbidden
	}

	if _, err := jwt.ParseWithClaims(strings.Split(string(tokenGet), " ")[1], &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWT.Key), nil
	}); err != nil {
		return nil, "handler: jwt parse with claims: " + err.Error(), fasthttp.StatusUnauthorized
	}

	id, err := strconv.Atoi(claims.Audience[0])
	if err != nil {
		return nil, "handler: atoi audience: " + err.Error(), fasthttp.StatusInternalServerError
	}

	chatID, err := strconv.Atoi(strings.Trim(string(ctx.QueryArgs().Peek("chat-id")), " "))
	if err != nil {
		return nil, "handler: atoi chat-id: " + err.Error(), fasthttp.StatusInternalServerError
	}

	key, errCustom := api_db.GetChatKeyDB(int64(chatID), int64(id), db)
	if errCustom != nil {
		return nil, "handler: get chat messages DB: " + errCustom.Error(), errCustom.Code
	}

	res = &model.GetChatKeyRes{Key: key}

	return res, "OK", fasthttp.StatusOK
}

func CheckChatHandler(ctx *fasthttp.RequestCtx) (res *model.CheckChatRes, message string, code int) {
	defer func() {
		resFinal := model.Response{Data: res, Code: code, Message: message}
		if code != fasthttp.StatusOK {
			log.Error().Err(errors.New(message)).Msg("")
		}
		jsonRes, _ := json.Marshal(resFinal)
		ctx.Response.SetStatusCode(resFinal.Code)
		ctx.Response.SetBody(jsonRes)
	}()

	if !ctx.IsPost() {
		return nil, "handler: wrong method", fasthttp.StatusMethodNotAllowed
	}

	tokenGet := ctx.Request.Header.Peek("Authorization")
	if tokenGet == nil {
		return nil, "handler: no authorization header", fasthttp.StatusUnauthorized
	}

	claims := model.JWTCustomClaims{}

	if len(strings.Split(string(tokenGet), " ")) == 1 {
		return nil, "handler: token is invalid", fasthttp.StatusForbidden
	}

	if _, err := jwt.ParseWithClaims(strings.Split(string(tokenGet), " ")[1], &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWT.Key), nil
	}); err != nil {
		return nil, "handler: jwt parse with claims: " + err.Error(), fasthttp.StatusUnauthorized
	}

	id, err := strconv.Atoi(claims.Audience[0])
	if err != nil {
		return nil, "handler: atoi audience: " + err.Error(), fasthttp.StatusInternalServerError
	}

	var req *model.CheckChatReq

	if err = json.Unmarshal(ctx.PostBody(), &req); err != nil {
		return nil, "handler: unmarshal request: " + err.Error() + ": wrong format of input data", fasthttp.StatusUnprocessableEntity
	}

	chatID := api_db.CheckExistingChatDB(int64(id), req.ReceiverID, db)
	if chatID == -1 {
		if errCustom := api_db.CheckExistingUserIDDB(req.ReceiverID, db); errCustom != nil {
			return nil, "handler: check existing user ID DB: " + errCustom.Error(), fasthttp.StatusUnprocessableEntity
		}
		if errCustom := api_db.CheckExistingUserIDDB(int64(id), db); errCustom != nil {
			return nil, "handler: check existing user ID DB: " + errCustom.Error(), fasthttp.StatusUnprocessableEntity
		}

		if errCustom := api_db.CreateChatDB(int64(id), req.ReceiverID, db); errCustom != nil {
			return nil, "handler: create chat DB: " + errCustom.Error(), errCustom.Code
		}

		chatID = api_db.CheckExistingChatDB(int64(id), req.ReceiverID, db)

		publicFirstKey, errCustom := api_db.GetUserKeyByIDDB(int64(id), db)
		if errCustom != nil {
			return nil, "handler: get first user key by id DB: " + errCustom.Error(), errCustom.Code
		}

		publicSecondKey, errCustom := api_db.GetUserKeyByIDDB(req.ReceiverID, db)
		if errCustom != nil {
			return nil, "handler: get second user key by id DB: " + errCustom.Error(), errCustom.Code
		}

		publicFirstKeyByte := make([]byte, len(publicFirstKey))
		for idx, elem := range publicFirstKey {
			publicFirstKeyByte[idx] = byte(elem)
		}

		publicSecondKeyByte := make([]byte, len(publicSecondKey))
		for idx, elem := range publicSecondKey {
			publicSecondKeyByte[idx] = byte(elem)
		}

		publicFirst, err := x509.ParsePKCS1PublicKey(publicFirstKeyByte)
		if err != nil {
			return nil, "handler: parse first public key: " + err.Error(), fasthttp.StatusInternalServerError
		}

		publicSecond, err := x509.ParsePKCS1PublicKey(publicSecondKeyByte)
		if err != nil {
			return nil, "handler: parse second public key: " + err.Error(), fasthttp.StatusInternalServerError
		}

		chatKey := make([]byte, 32)
		if _, err = rand.Read(chatKey); err != nil {
			return nil, "handler: read: " + err.Error(), fasthttp.StatusInternalServerError
		}

		firstChatEncrypt, err := rsa.EncryptPKCS1v15(rand.Reader, publicFirst, chatKey)
		if err != nil {
			return nil, "handler: encrypt chat key: " + err.Error(), fasthttp.StatusInternalServerError
		}
		firstChatEncryptInt := make([]int32, len(firstChatEncrypt))
		for idx, elem := range firstChatEncrypt {
			firstChatEncryptInt[idx] = int32(elem)
		}

		secondChatEncrypt, err := rsa.EncryptPKCS1v15(rand.Reader, publicSecond, chatKey)
		if err != nil {
			return nil, "handler: encrypt chat key: " + err.Error(), fasthttp.StatusInternalServerError
		}
		secondChatEncryptInt := make([]int32, len(secondChatEncrypt))
		for idx, elem := range secondChatEncrypt {
			secondChatEncryptInt[idx] = int32(elem)
		}

		if errCustom = api_db.CreateChatKeyDB(&model.ChatKeyItem{
			ChatID:    chatID,
			FirstUID:  int64(id),
			FirstKey:  firstChatEncryptInt,
			SecondUID: req.ReceiverID,
			SecondKey: secondChatEncryptInt,
		}, db); errCustom != nil {
			return nil, "handler: create chat key DB: " + errCustom.Error(), errCustom.Code
		}

		res = &model.CheckChatRes{Key: firstChatEncryptInt}
	} else {
		key, errCustom := api_db.GetChatKeyDB(chatID, int64(id), db)
		if errCustom != nil {
			return nil, "handler: get chat key DB: " + errCustom.Error(), errCustom.Code
		}
		res = &model.CheckChatRes{Key: key}
	}

	return res, "OK", fasthttp.StatusOK
}

func GetChatByUsersHandler(ctx *fasthttp.RequestCtx) (res int64, message string, code int) {
	defer func() {
		resFinal := model.Response{Data: res, Code: code, Message: message}
		if code != fasthttp.StatusOK {
			log.Error().Err(errors.New(message)).Msg("")
		}
		jsonRes, _ := json.Marshal(resFinal)
		ctx.Response.SetStatusCode(resFinal.Code)
		ctx.Response.SetBody(jsonRes)
	}()

	if !ctx.IsGet() {
		return -1, "handler: wrong method", fasthttp.StatusMethodNotAllowed
	}

	tokenGet := ctx.Request.Header.Peek("Authorization")
	if tokenGet == nil {
		return -1, "handler: no authorization header", fasthttp.StatusUnauthorized
	}

	claims := model.JWTCustomClaims{}

	if len(strings.Split(string(tokenGet), " ")) == 1 {
		return -1, "handler: token is invalid", fasthttp.StatusForbidden
	}

	if _, err := jwt.ParseWithClaims(strings.Split(string(tokenGet), " ")[1], &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWT.Key), nil
	}); err != nil {
		return -1, "handler: jwt parse with claims: " + err.Error(), fasthttp.StatusUnauthorized
	}

	id, err := strconv.Atoi(claims.Audience[0])
	if err != nil {
		return -1, "handler: atoi audience: " + err.Error(), fasthttp.StatusInternalServerError
	}

	uid, err := strconv.Atoi(strings.Trim(string(ctx.QueryArgs().Peek("id")), " "))
	if err != nil {
		return -1, "handler: atoi id: " + err.Error(), fasthttp.StatusInternalServerError
	}

	res, errCustom := api_db.GetChatByUsersDB(int64(id), int64(uid), db)
	if errCustom != nil {
		return -1, "handler: get chat by users DB: " + errCustom.Error(), errCustom.Code
	}

	return res, "OK", fasthttp.StatusOK
}

func GetAllChatHandler(ctx *fasthttp.RequestCtx) (res []*model.ChatShortItem, message string, code int) {
	defer func() {
		resFinal := model.Response{Data: res, Code: code, Message: message}
		if code != fasthttp.StatusOK {
			log.Error().Err(errors.New(message)).Msg("")
		}
		jsonRes, _ := json.Marshal(resFinal)
		ctx.Response.SetStatusCode(resFinal.Code)
		ctx.Response.SetBody(jsonRes)
	}()

	if !ctx.IsGet() {
		return nil, "handler: wrong method", fasthttp.StatusMethodNotAllowed
	}

	tokenGet := ctx.Request.Header.Peek("Authorization")
	if tokenGet == nil {
		return nil, "handler: no authorization header", fasthttp.StatusUnauthorized
	}

	claims := model.JWTCustomClaims{}

	if len(strings.Split(string(tokenGet), " ")) == 1 {
		return nil, "handler: token is invalid", fasthttp.StatusForbidden
	}

	if _, err := jwt.ParseWithClaims(strings.Split(string(tokenGet), " ")[1], &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWT.Key), nil
	}); err != nil {
		return nil, "handler: jwt parse with claims: " + err.Error(), fasthttp.StatusUnauthorized
	}

	id, err := strconv.Atoi(claims.Audience[0])
	if err != nil {
		return nil, "handler: atoi audience: " + err.Error(), fasthttp.StatusInternalServerError
	}

	res, errCustom := api_db.GetUserChatsDB(int64(id), db)
	if errCustom != nil {
		return nil, "handler: get user chats DB: " + errCustom.Error(), errCustom.Code
	}

	return res, "OK", fasthttp.StatusOK
}

func GetChatHandler(ctx *fasthttp.RequestCtx) (res []*model.MessageShortItem, message string, code int) {
	defer func() {
		resFinal := model.Response{Data: res, Code: code, Message: message}
		if code != fasthttp.StatusOK {
			log.Error().Err(errors.New(message)).Msg("")
		}
		jsonRes, _ := json.Marshal(resFinal)
		ctx.Response.SetStatusCode(resFinal.Code)
		ctx.Response.SetBody(jsonRes)
	}()

	if !ctx.IsGet() {
		return nil, "handler: wrong method", fasthttp.StatusMethodNotAllowed
	}

	tokenGet := ctx.Request.Header.Peek("Authorization")
	if tokenGet == nil {
		return nil, "handler: no authorization header", fasthttp.StatusUnauthorized
	}

	claims := model.JWTCustomClaims{}

	if len(strings.Split(string(tokenGet), " ")) == 1 {
		return nil, "handler: token is invalid", fasthttp.StatusForbidden
	}

	if _, err := jwt.ParseWithClaims(strings.Split(string(tokenGet), " ")[1], &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWT.Key), nil
	}); err != nil {
		return nil, "handler: jwt parse with claims: " + err.Error(), fasthttp.StatusUnauthorized
	}

	id, err := strconv.Atoi(claims.Audience[0])
	if err != nil {
		return nil, "handler: atoi audience: " + err.Error(), fasthttp.StatusInternalServerError
	}

	chatID, err := strconv.Atoi(strings.Trim(string(ctx.QueryArgs().Peek("chat-id")), " "))
	if err != nil {
		return nil, "handler: atoi chat-id: " + err.Error(), fasthttp.StatusInternalServerError
	}

	res, errCustom := api_db.GetChatMessagesDB(int64(id), int64(chatID), db)
	if errCustom != nil {
		return nil, "handler: get chat messages DB: " + errCustom.Error(), errCustom.Code
	}

	return res, "OK", fasthttp.StatusOK
}
