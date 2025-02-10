package model

type UserItem struct {
	ID         int64    `json:"id" reindex:"id,hash,pk"`
	Email      string   `json:"email" reindex:"email,hash"`
	Password   string   `json:"password" reindex:"password,hash"`
	Photo      string   `json:"photo" reindex:"photo,hash"`
	Surname    string   `json:"surname" reindex:"surname,hash"`
	Name       string   `json:"name" reindex:"name,hash"`
	Status     int32    `json:"status" reindex:"status,hash"`
	CreateDate int64    `json:"create_date" reindex:"create_date,hash"`
	OpenKey    string   `json:"open_key" reindex:"open_key,hash"`
	_          struct{} `reindex:"email+surname+name=search,text,composite"`
}

type SearchUserRes struct {
	Data    []*UserSnippetItem `json:"data"`
	Code    int                `json:"code"`
	Message string             `json:"message"`
}

type UserSnippetItem struct {
	ID      int64  `json:"id"`
	ChatID  int64  `json:"chat_id"`
	Photo   string `json:"photo"`
	Surname string `json:"surname"`
	Name    string `json:"name"`
	Email   string `json:"email"`
}

type GetAccountRes struct {
	Data    *AccountItem `json:"data"`
	Code    int          `json:"code"`
	Message string       `json:"message"`
}

type UpdateOpenKeyReq struct {
	Key []int32 `json:"key"`
}

type UpdateOpenKeyRes struct {
	Data    any    `json:"data"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type AccountItem struct {
	ID         int64  `json:"id"`
	Email      string `json:"email"`
	Photo      string `json:"photo"`
	Surname    string `json:"surname"`
	Name       string `json:"name"`
	CreateDate int64  `json:"create_date"`
}
