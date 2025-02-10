package model

type GetAllChatRes struct {
	Data    []*ChatShortItem `json:"data"`
	Code    int              `json:"code"`
	Message string           `json:"message"`
}

type ChatGuiItem struct {
	Surname string
	Name    string
	Text    string
	Time    string
}

type ChatShortItem struct {
	ID              int64   `json:"id"`
	ReceiverID      int64   `json:"receiver_id"`
	ReceiverName    string  `json:"receiver_name"`
	ReceiverSurname string  `json:"receiver_surname"`
	ReceiverPhoto   string  `json:"receiver_photo"`
	LastMessageText []int32 `json:"last_message_text"`
	LastMessageDate int64   `json:"last_message_date"`
	LastMessageType int32   `json:"last_message_type"`
	Read            int32   `json:"read"`
	MyMessage       int32   `json:"my_message"`
}

type CheckChatReq struct {
	ReceiverID int64 `json:"receiver_id"`
}

type CheckChatRes struct {
	Data    *CheckChatItem `json:"data"`
	Code    int            `json:"code"`
	Message string         `json:"message"`
}

type CheckChatItem struct {
	Key []int32 `json:"key"`
}

type GetChatKeyRes struct {
	Data    *GetChatKeyItem `json:"data"`
	Code    int             `json:"code"`
	Message string          `json:"message"`
}

type GetChatKeyItem struct {
	Key []int32 `json:"key"`
}
