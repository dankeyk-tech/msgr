package model

import (
	custom_errors "diplom-chat-gost-server/pkg/custom-errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/valyala/fasthttp"
	"strings"
)

type Response struct {
	Data    any    `json:"data"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type SingInReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SingInRes struct {
	Token string `json:"token"`
}

type SingUpReq struct {
	Surname  string `json:"surname"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SingUpRes struct {
	Token string `json:"token"`
}

type ChangePassReq struct {
	Password string `json:"password"`
}

type UpdateOpenKeyReq struct {
	Key []int32 `json:"key"`
}

type ChangePassRes struct {
	Email string `json:"email"`
	Hash  string `json:"hash"`
}

type ChangePassConfirmReq struct {
	Hash     string `json:"hash"`
	AuthCode int    `json:"auth_code"`
}

type ChangePassConfirmRes struct {
	Token string `json:"token"`
}

type JWTCustomClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func (req *SingInReq) Validate() *custom_errors.ErrHttp {
	req.Email = strings.Trim(req.Email, " ")
	req.Password = strings.Trim(req.Password, " ")
	if req.Email == "" {
		return custom_errors.New(fasthttp.StatusUnprocessableEntity, "field email can't be empty")
	}
	if req.Password == "" {
		return custom_errors.New(fasthttp.StatusUnprocessableEntity, "field password can't be empty")
	}

	return nil
}

func (req *SingUpReq) Validate() *custom_errors.ErrHttp {
	req.Email = strings.Trim(req.Email, " ")
	req.Password = strings.Trim(req.Password, " ")
	req.Surname = strings.Trim(req.Surname, " ")
	req.Name = strings.Trim(req.Name, " ")
	if req.Email == "" {
		return custom_errors.New(fasthttp.StatusUnprocessableEntity, "field email can't be empty")
	}
	if req.Password == "" {
		return custom_errors.New(fasthttp.StatusUnprocessableEntity, "field password can't be empty")
	}
	if req.Surname == "" {
		return custom_errors.New(fasthttp.StatusUnprocessableEntity, "field surname can't be empty")
	}
	if req.Name == "" {
		return custom_errors.New(fasthttp.StatusUnprocessableEntity, "field name can't be empty")
	}

	return nil
}

func (req *ChangePassConfirmReq) Validate() *custom_errors.ErrHttp {
	req.Hash = strings.Trim(req.Hash, " ")
	if req.Hash == "" {
		return custom_errors.New(fasthttp.StatusUnprocessableEntity, "field hash can't be empty")
	}
	if req.AuthCode == 0 {
		return custom_errors.New(fasthttp.StatusUnprocessableEntity, "field auth_code can't be empty")
	}

	return nil
}
