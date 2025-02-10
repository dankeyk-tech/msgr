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

func UpdateOpenKeyHandler(ctx *fasthttp.RequestCtx) (message string, code int) {
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

	var req *model.UpdateOpenKeyReq

	if err := json.Unmarshal(ctx.PostBody(), &req); err != nil {
		return "handler: unmarshal request: " + err.Error() + ": wrong format of input data", fasthttp.StatusUnprocessableEntity
	}

	id, err := strconv.Atoi(claims.Audience[0])
	if err != nil {
		return "handler: atoi audience: " + err.Error(), fasthttp.StatusInternalServerError
	}

	if errCustom := api_db.UpdateOpenKeyDB(int64(id), req.Key, db); errCustom != nil {
		return "handler: update open key DB: " + errCustom.Error(), errCustom.Code
	}

	return "OK", fasthttp.StatusOK
}

/*
func ChangePasswordHandler(ctx *fasthttp.RequestCtx) (res *model.ChangePassRes, message string, code int) {
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

	var req *model.ChangePassReq

	if err := json.Unmarshal(ctx.PostBody(), &req); err != nil {
		return nil, "handler: unmarshal request: " + err.Error() + ": wrong format of input data", fasthttp.StatusUnprocessableEntity
	}

	salt := make([]byte, 256)
	if _, err := crand.Read(salt); err != nil {
		return nil, "handler: crypto read: " + err.Error(), fasthttp.StatusInternalServerError
	}

	hashVer := sha3.Sum512([]byte(fmt.Sprintf("%s%s", claims.Email, hex.EncodeToString(salt))))

	authCode := rand.Intn(999999-100000+1) + 100000

	go email.SendMail(claims.Email, "Смена пароля", "Для смены пароля используйте код: "+strconv.Itoa(authCode), config.Email)

	if errCustom := api_db.CreateAuthDB(&model.AuthItem{
		Key:      string_builder.BuildStrings([]string{claims.Email, "_-_", fmt.Sprintf("%x", hashVer)}),
		Salt:     hex.EncodeToString(salt),
		Password: req.Password,
		AuthCode: authCode,
		Date:     time.Now().Unix(),
	}, db); errCustom != nil {
		return nil, "handler: create auth DB: " + errCustom.Error(), errCustom.Code
	}

	return &model.ChangePassRes{
		Email: claims.Email,
		Hash:  fmt.Sprintf("%x", hashVer),
	}, "OK", fasthttp.StatusOK
}

func ChangePasswordConfirmHandler(ctx *fasthttp.RequestCtx) (message string, code int) {
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

	var req *model.ChangePassConfirmReq

	if err := json.Unmarshal(ctx.PostBody(), &req); err != nil {
		return "handler: unmarshal request: " + err.Error() + ": wrong format of input data", fasthttp.StatusUnprocessableEntity
	}

	if errCustom := req.Validate(); errCustom != nil {
		return "handler: validate request: " + errCustom.Error(), errCustom.Code
	}

	auth, errCustom := api_db.GetAuthDB(string_builder.BuildStrings([]string{claims.Email, "_-_", req.Hash}), db)
	if errCustom != nil {
		return "handler: get auth DB: " + errCustom.Error(), errCustom.Code
	}

	passChan := make(chan string)
	go password.GenPass(auth.Password, passChan)

	if auth.AuthCode != req.AuthCode {
		return "handler: check auth code: wrong code", fasthttp.StatusForbidden
	}

	if errCustom = api_db.UpdatePasswordDB(<-passChan, claims.Email, db); errCustom != nil {
		return "handler: update password DB: " + errCustom.Error(), errCustom.Code
	}

	return "OK", fasthttp.StatusOK
}
*/
/*
func UpdatePhotoHandler(ctx *fasthttp.RequestCtx) (message string, code int) {
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

	file, err := ctx.FormFile("file")
	if err != nil {
		return "handler: form file: " + err.Error(), fasthttp.StatusUnprocessableEntity
	}

	path, errCustom := files.SaveAvatar(file, int64(id), config.FilePath)
	if errCustom != nil {
		return "handler: save avatar: " + errCustom.Error(), errCustom.Code
	}

	if errCustom = api_db.UpdatePhotoDB(path, int64(id), db); errCustom != nil {
		return "handler: update photo: " + errCustom.Error(), errCustom.Code
	}

	return "OK", fasthttp.StatusOK
}
*/

func UpdateDataHandler(ctx *fasthttp.RequestCtx) (message string, code int) {
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

	var req *model.UpdateDataReq

	if err = json.Unmarshal(ctx.PostBody(), &req); err != nil {
		return "handler: unmarshal request: " + err.Error() + ": wrong format of input data", fasthttp.StatusUnprocessableEntity
	}

	if errCustom := req.Validate(); errCustom != nil {
		return "handler: validate request: " + errCustom.Error(), errCustom.Code
	}

	if errCustom := api_db.UpdateDataDB(req, int64(id), db); errCustom != nil {
		return "handler: update data DB: " + errCustom.Error(), errCustom.Code
	}

	return "OK", fasthttp.StatusOK
}

func GetAccountHandler(ctx *fasthttp.RequestCtx) (res *model.GetAccountRes, message string, code int) {
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

	res, errCustom := api_db.GetAccountDB(int64(id), db)
	if errCustom != nil {
		return nil, "handler: get account DB: " + errCustom.Error(), errCustom.Code
	}

	return res, "OK", fasthttp.StatusOK
}
