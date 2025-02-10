package model

import (
	"errors"
)

type SingInReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SingInRes struct {
	Data    SingInResData `json:"data"`
	Code    int           `json:"code"`
	Message string        `json:"message"`
}

type SingInResData struct {
	Token string `json:"token"`
}

func (req *SingInReq) Validate() error {
	if req.Email == "" {
		return errors.New("Поле \"Почта\" не может быть пустым!")
	}

	if req.Password == "" {
		return errors.New("Поле \"Пароль\" не может быть пустым!")
	}

	return nil
}
