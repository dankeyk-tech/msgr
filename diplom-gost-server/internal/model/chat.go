package model

type ChatItem struct {
	ID              int64          `json:"id" reindex:"id,hash,pk"`
	UIDs            []int64        `json:"uids" reindex:"uids,hash"`
	LastMessageID   int64          `json:"last_message_id" reindex:"last_message_id,hash"`
	LastMessageDate int64          `json:"last_message_date" reindex:"last_message_date,hash"`
	User            []*UserItem    `json:"-" reindex:"user,,joined"`
	Message         []*MessageItem `json:"-" reindex:"message,,joined"`
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

type CheckChatRes struct {
	Key []int32 `json:"key"`
}

type CheckChatReq struct {
	ReceiverID int64 `json:"receiver_id"`
}
