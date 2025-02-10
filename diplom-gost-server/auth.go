package main

import (
	api_db "diplom-chat-gost-server/internal/api-db"
	"diplom-chat-gost-server/internal/model"
	"diplom-chat-gost-server/pkg/jwtRegister"
	"diplom-chat-gost-server/pkg/password"
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
	"strconv"
	"time"
)

func SignUpHandler(ctx *fasthttp.RequestCtx) (res *model.SingUpRes, message string, code int) {
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

	var req *model.SingUpReq

	if err := json.Unmarshal(ctx.PostBody(), &req); err != nil {
		return nil, "handler: unmarshal request: " + err.Error() + ": wrong format of input data", fasthttp.StatusUnprocessableEntity
	}

	if errCustom := req.Validate(); errCustom != nil {
		return nil, "handler: validate request: " + errCustom.Error(), errCustom.Code
	}

	if errCustom := api_db.CheckExistingUserEmailDB(req.Email, db); errCustom != nil {
		return nil, "handler: check existing user email DB: " + errCustom.Error(), errCustom.Code
	}

	passChan := make(chan string)
	go password.GenPass(req.Password, passChan)

	uid, errCustom := api_db.CreateUserDB(&model.UserItem{
		ID:         0,
		Email:      req.Email,
		Password:   <-passChan,
		Photo:      "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",
		Surname:    req.Surname,
		Name:       req.Name,
		Status:     1,
		CreateDate: time.Now().Unix(),
		OpenKey:    []int32{},
	}, db)
	if errCustom != nil {
		return nil, "handler: create user DB: " + errCustom.Error(), errCustom.Code
	}

	token, errCustom := jwtRegister.GenerateToken(&model.JWTCustomClaims{
		Email: req.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "https://" + config.JWT.Domain,
			Audience:  jwt.ClaimStrings{strconv.Itoa(int(uid))},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})
	if errCustom != nil {
		return
	}

	return &model.SingUpRes{Token: token}, "OK", fasthttp.StatusOK
}

func SignInHandler(ctx *fasthttp.RequestCtx) (res *model.SingInRes, message string, code int) {
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

	var req *model.SingInReq

	if err := json.Unmarshal(ctx.PostBody(), &req); err != nil {
		return nil, "handler: unmarshal request: " + err.Error() + ": wrong format of input data", fasthttp.StatusUnprocessableEntity
	}

	if errCustom := req.Validate(); errCustom != nil {
		return nil, "handler: validate request: " + errCustom.Error(), errCustom.Code
	}

	user, errCustom := api_db.GetUserByEmailDB(req.Email, db)
	if errCustom != nil {
		return nil, "handler: get user by email DB: " + errCustom.Error(), errCustom.Code
	}

	passChan := make(chan bool)
	go password.CheckPass(user.Password, req.Password, passChan)

	if user.Status == 0 {
		return nil, "handler: check status failed: user is banned", fasthttp.StatusForbidden
	}

	if !<-passChan {
		return nil, "handler: check password failed: wrong password", fasthttp.StatusForbidden
	}

	token, errCustom := jwtRegister.GenerateToken(&model.JWTCustomClaims{
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "https://" + config.JWT.Domain,
			Audience:  jwt.ClaimStrings{strconv.Itoa(int(user.ID))},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})
	if errCustom != nil {
		return
	}

	return &model.SingInRes{
		Token: token,
	}, "OK", fasthttp.StatusOK
}
