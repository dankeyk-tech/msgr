package main

import (
	api_db "diplom-chat-gost-server/internal/api-db"
	"diplom-chat-gost-server/internal/model"
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
	"strconv"
	"strings"
	"time"
)

func SendMessageHandler(ctx *fasthttp.RequestCtx) (message string, code int) {
	defer func() {
		resFinal := model.Response{Data: nil, Code: code, Message: message}
		if code != fasthttp.StatusOK {
			log.Error().Err(errors.New(message)).Msg("")
		}
		jsonRes, _ := json.Marshal(resFinal)
		ctx.Response.SetStatusCode(resFinal.Code)
		ctx.Response.SetBody(jsonRes)
	}()

	if !ctx.IsPost() {
		return "handler: wrong method", fasthttp.StatusMethodNotAllowed
	}

	tokenGet := ctx.Request.Header.Peek("Authorization")
	if tokenGet == nil {
		return "handler: no authorization header", fasthttp.StatusUnauthorized
	}

	claims := model.JWTCustomClaims{}

	if len(strings.Split(string(tokenGet), " ")) == 1 {
		return "handler: token is invalid", fasthttp.StatusForbidden
	}

	if _, err := jwt.ParseWithClaims(strings.Split(string(tokenGet), " ")[1], &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWT.Key), nil
	}); err != nil {
		return "handler: jwt parse with claims: " + err.Error(), fasthttp.StatusUnauthorized
	}

	id, err := strconv.Atoi(claims.Audience[0])
	if err != nil {
		return "handler: atoi audience: " + err.Error(), fasthttp.StatusInternalServerError
	}

	var req *model.SendMessageReq

	if err = json.Unmarshal(ctx.PostBody(), &req); err != nil {
		return "handler: unmarshal request: " + err.Error() + ": wrong format of input data", fasthttp.StatusUnprocessableEntity
	}

	chatID := api_db.CheckExistingChatDB(int64(id), req.ReceiverID, db)

	if errCustom := api_db.CreateMessageDB(&model.MessageItem{
		ID:          0,
		ChatID:      chatID,
		UID:         int64(id),
		MessageType: req.Type,
		Text:        req.Text,
		Read:        0,
		Date:        time.Now().Unix(),
	}, db); errCustom != nil {
		return "handler: create message DB: " + errCustom.Error(), errCustom.Code
	}

	return "OK", fasthttp.StatusOK
}
