package model

import (
	custom_errors "diplom-chat-gost-server/pkg/custom-errors"
	"github.com/valyala/fasthttp"
	"strings"
)

type UserItem struct {
	ID         int64    `json:"id" reindex:"id,hash,pk"`
	Email      string   `json:"email" reindex:"email,hash"`
	Password   string   `json:"password" reindex:"password,hash"`
	Photo      string   `json:"photo" reindex:"photo,hash"`
	Surname    string   `json:"surname" reindex:"surname,hash"`
	Name       string   `json:"name" reindex:"name,hash"`
	Status     int32    `json:"status" reindex:"status,hash"`
	CreateDate int64    `json:"create_date" reindex:"create_date,hash"`
	OpenKey    []int32  `json:"open_key" reindex:"open_key,hash"`
	_          struct{} `reindex:"email+surname+name=search,text,composite"`
}

type GetAccountRes struct {
	ID         int64  `json:"id"`
	Email      string `json:"email"`
	Photo      string `json:"photo"`
	Surname    string `json:"surname"`
	Name       string `json:"name"`
	CreateDate int64  `json:"create_date"`
}

type UserSnippetItem struct {
	ID      int64  `json:"id"`
	Surname string `json:"surname"`
	Name    string `json:"name"`
	Photo   string `json:"photo"`
	Email   string `json:"email"`
	ChatID  int64  `json:"chat_id"`
}

type UpdateDataReq struct {
	Surname string `json:"surname"`
	Name    string `json:"name"`
}

func (req *UpdateDataReq) Validate() *custom_errors.ErrHttp {
	req.Surname = strings.Trim(req.Surname, " ")
	req.Name = strings.Trim(req.Name, " ")
	if req.Surname == "" {
		return custom_errors.New(fasthttp.StatusUnprocessableEntity, "field surname can't be empty")
	}
	if req.Name == "" {
		return custom_errors.New(fasthttp.StatusUnprocessableEntity, "field name can't be empty")
	}

	return nil
}
