package model

type MessageShortItem struct {
	ID          int64   `json:"id"`
	MessageType int32   `json:"message_type"`
	Text        []int32 `json:"text"`
	Date        int64   `json:"date"`
	MyMessage   bool    `json:"my_message"`
}

type GetAllMessagesRes struct {
	Data    []*MessageShortItem `json:"data"`
	Code    int                 `json:"code"`
	Message string              `json:"message"`
}

type GetChatByUsersRes struct {
	Data    int64  `json:"data"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type SendMessageReq struct {
	ReceiverID int64   `json:"receiver_id"`
	Text       []int32 `json:"text"`
	Type       int     `json:"type"`
}

type SendMessageRes struct {
	Data    any    `json:"data"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}
