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
)

func SearchUserHandler(ctx *fasthttp.RequestCtx) (res []*model.UserSnippetItem, message string, code int) {
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

	search := strings.Trim(string(ctx.QueryArgs().Peek("search")), " ")

	res, errCustom := api_db.SearchUserDB(int64(id), search, db)
	if errCustom != nil {
		return nil, "handler: search user DB: " + errCustom.Error(), errCustom.Code
	}

	for idx := range res {
		res[idx].ChatID = api_db.SearchUserChat([]int64{res[idx].ID, int64(id)}, db)
	}

	return res, "OK", fasthttp.StatusOK
}
